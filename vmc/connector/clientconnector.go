/* Copyright 2019-2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

// Package connector provides helper methods that provides client.Connector, required to call VMC APIs.
package connector

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vmware/terraform-provider-vmc/vmc/constants"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/security"
)

type Authenticator interface {
	Authenticate() error
}

type ConnectorWrapper struct {
	client.Connector
	RefreshToken string
	OrgID        string
	VmcURL       string
	CspURL       string
}

func (c *ConnectorWrapper) Authenticate() error {
	var err error
	httpClient := http.Client{}
	c.Connector, err = NewClientConnectorByRefreshToken(c.RefreshToken, c.VmcURL, c.CspURL, httpClient)
	if err != nil {
		return err
	}
	return nil
}

// NewClientConnectorByRefreshToken returns client connector to any VMC service by using OAuth authentication using Refresh Token.
func NewClientConnectorByRefreshToken(refreshToken, serviceUrl, cspURL string,
	httpClient http.Client) (client.Connector, error) {

	if len(serviceUrl) <= 0 {
		serviceUrl = constants.DefaultVmcUrl
	}

	if len(cspURL) <= 0 {
		cspURL = constants.DefaultCspUrl +
			constants.CspRefreshUrlSuffix
	} else {
		cspURL = cspURL +
			constants.CspRefreshUrlSuffix
	}

	securityCtx, err := SecurityContextByRefreshToken(refreshToken, cspURL)
	if err != nil {
		return nil, err
	}

	connector := client.NewRestConnector(serviceUrl, httpClient)
	connector.SetSecurityContext(securityCtx)

	return connector, nil
}

// SecurityContextByRefreshToken returns Security Context with access token that is received from Cloud Service Provider using Refresh Token by OAuth authentication scheme.
func SecurityContextByRefreshToken(refreshToken string, cspURL string) (core.SecurityContext, error) {
	payload := strings.NewReader("refresh_token=" + refreshToken)

	req, _ := http.NewRequest("POST", cspURL, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		b, _ := io.ReadAll(res.Body)
		return nil, fmt.Errorf("response from Cloud Service Provider contains status code %d : %s", res.StatusCode, string(b))
	}

	defer res.Body.Close()

	var jsondata map[string]interface{}
	err = json.NewDecoder(res.Body).Decode(&jsondata)
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
		return nil, errors.New("Cloud Service Provider authentication response does not contain access token")
	}

	securityCtx := security.NewOauthSecurityContext(accessToken)
	return securityCtx, nil
}