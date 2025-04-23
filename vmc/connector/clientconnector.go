// © Broadcom. All Rights Reserved.
// The term "Broadcom" refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: MPL-2.0

// Package connector provides helper methods that provides client.Connector, required to call VMC APIs.
package connector

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/security"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/vmware/terraform-provider-vmc/vmc/constants"
)

type Authenticator interface {
	Authenticate() error
}

type Wrapper struct {
	client.Connector
	RefreshToken string
	ClientID     string
	ClientSecret string
	OrgID        string
	VmcURL       string
	CspURL       string
}

func CopyWrapper(original Wrapper) *Wrapper {
	return &original
}

func (c *Wrapper) Authenticate() error {
	var err error
	httpClient := http.Client{}
	if len(c.RefreshToken) > 0 {
		c.Connector, err = newClientConnectorByRefreshToken(c.RefreshToken, c.VmcURL, c.CspURL, &httpClient)
		if err != nil {
			return err
		}
		return nil
	}
	if len(c.ClientID) > 0 && len(c.ClientSecret) > 0 {
		c.Connector, err = newClientConnectorByClientID(c.ClientID, c.ClientSecret, c.VmcURL, c.CspURL, &httpClient)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("no refreshToken or ClientID/ClientSecret provided")
}

// newClientConnectorByRefreshToken returns client connector to any VMC service by using OAuth authentication using Refresh Token.
func newClientConnectorByRefreshToken(refreshToken, serviceURL, cspURL string,
	httpClient *http.Client) (client.Connector, error) {

	if len(serviceURL) <= 0 {
		serviceURL = constants.DefaultVmcURL
	}

	if len(cspURL) <= 0 {
		cspURL = constants.DefaultCspURL +
			constants.CspRefreshURLSuffix
	} else {
		cspURL = cspURL +
			constants.CspRefreshURLSuffix
	}

	securityCtx, err := securityContextByRefreshToken(refreshToken, cspURL)
	if err != nil {
		return nil, err
	}

	connector := client.NewConnector(serviceURL, client.UsingRest(nil),
		client.WithHttpClient(httpClient), client.WithSecurityContext(securityCtx))

	return connector, nil
}

// SecurityContextByRefreshToken returns Security Context with access token that is received from Cloud Service Provider using Refresh Token by OAuth authentication scheme.
func securityContextByRefreshToken(refreshToken string, cspURL string) (core.SecurityContext, error) {
	payload := strings.NewReader("refresh_token=" + refreshToken)

	req, _ := http.NewRequest("POST", cspURL, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	securityCtx, err := parseAuthnResponse(res)
	if err != nil {
		return nil, err
	}
	return securityCtx, nil
}

// newClientConnectorByClientID returns client connector to any VMC service by using OAuth authentication using clientId and secret.
func newClientConnectorByClientID(clientID, clientSecret, serviceURL, cspURL string,
	httpClient *http.Client) (client.Connector, error) {

	if len(serviceURL) <= 0 {
		serviceURL = constants.DefaultVmcURL
	}

	if len(cspURL) <= 0 {
		cspURL = constants.DefaultCspURL +
			constants.CspTokenURLSuffix
	} else {
		cspURL = cspURL +
			constants.CspTokenURLSuffix
	}

	securityCtx, err := securityContextByClientID(clientID, clientSecret, cspURL)
	if err != nil {
		return nil, err
	}

	connector := client.NewConnector(serviceURL, client.UsingRest(nil),
		client.WithHttpClient(httpClient), client.WithSecurityContext(securityCtx))

	return connector, nil
}

func securityContextByClientID(clientID string, clientSecret string, cspTokenEndpointURL string) (core.SecurityContext, error) {
	oauth2Config := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     cspTokenEndpointURL,
	}
	token, err := oauth2Config.Token(context.TODO())
	if err != nil {
		return nil, err
	}
	return security.NewOauthSecurityContext(token.AccessToken), nil
}

func parseAuthnResponse(response *http.Response) (*security.OauthSecurityContext, error) {
	if response.StatusCode != 200 {
		b, _ := io.ReadAll(response.Body)
		return nil, fmt.Errorf("response from Cloud Service Provider contains status code %d : %s", response.StatusCode, string(b))
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}(response.Body)

	var jsondata map[string]interface{}
	err := json.NewDecoder(response.Body).Decode(&jsondata)
	if err != nil {
		return nil, fmt.Errorf("error decoding response : %v", err)
	}

	var accessToken string
	if token, ok := jsondata["access_token"]; ok {
		if accessTokenStr, ok := token.(string); ok {
			accessToken = accessTokenStr
		} else {
			errMsg := fmt.Sprintf("Invalid type for access_token, expected string, actual %s", reflect.TypeOf(token).String())
			return nil, errors.New(errMsg)
		}
	} else {
		return nil, errors.New("cloud Service Provider authentication response does not contain access token")
	}

	securityCtx := security.NewOauthSecurityContext(accessToken)
	return securityCtx, nil
}
