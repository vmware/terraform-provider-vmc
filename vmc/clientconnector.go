/* Copyright 2019 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

// Package vmc provides helper methods that provides client.Connector, required to call VMC APIs.
package vmc

import (
	"context"
	"crypto/tls"
	"encoding/base64"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/security"
	"gitlab.eng.vmware.com/vapi-sdk/csp-go-openapi-bindings"
	"log"
	"net/http"
)

// NewClientConnectorByRefreshToken returns client connector to any VMC service by using OAuth authentication using Refresh Token.
func NewClientConnectorByRefreshToken(refreshToken, serviceUrl, cspURL string,
	httpClient http.Client) (client.Connector, error) {

	/*if len(serviceUrl) <= 0 {
		serviceUrl = DefaultVMCServer
	}

	if len(cspURL) <= 0 {
		cspURL = DefaultCSPUrl
	}*/

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
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	config := openapiclient.NewConfiguration()
	config.HTTPClient = &http.Client{Transport: tr}
	config.AddDefaultHeader("authorization", "Basic"+refreshToken)

	auth := context.WithValue(context.Background(), openapiclient.ContextBasicAuth, openapiclient.APIKey{
		Key:    refreshToken,
		Prefix: "Basic",
	})

	APIClient := openapiclient.NewAPIClient(config)

	accessToken, resp, err := APIClient.AuthenticationApi.GetAccessTokenByApiRefreshTokenUsingPOST(auth, refreshToken)
	log.Printf("Response status from the call  : %s",resp.Status)
	log.Printf("Response code from the call  : %d",resp.StatusCode)
	log.Println(resp.Request)
	log.Printf("Accesstoken  : %v",accessToken)
	if resp.StatusCode != 200 {
		return nil, HandleCreateError("accesstoken using refresh_token",err)
	}
	if err != nil {
		return nil, HandleCreateError("accesstoken using refresh_token",err)
	}


	securityCtx := security.NewOauthSecurityContext(accessToken.AccessToken)
	return securityCtx, nil
}

func NewClientConnectorByClientID(clientID, clientSecret, serviceUrl, cspURL string,
	httpClient http.Client) (client.Connector, error) {


	securityCtx, err := SecurityContextByClientID(clientID, clientSecret, cspURL)
	if err != nil {
		return nil, err
	}

	connector := client.NewRestConnector(serviceUrl, httpClient)
	connector.SetSecurityContext(securityCtx)

	return connector, nil
}

func SecurityContextByClientID(clientID string, clientSecret string, cspURL string) (core.SecurityContext, error) {

	clientCredentials := clientID + ":" + clientSecret

	encodedClientCredentials := base64.StdEncoding.EncodeToString([]byte(clientCredentials))

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	config := openapiclient.NewConfiguration()
	log.Printf("config basepath %s",config.BasePath)
	config.HTTPClient = &http.Client{Transport: tr}
	config.AddDefaultHeader("authorization", "Basic"+encodedClientCredentials)

	auth := context.WithValue(context.Background(), openapiclient.ContextBasicAuth, openapiclient.APIKey{
		Key:    encodedClientCredentials,
		Prefix: "Basic",
	})

	APIClient := openapiclient.NewAPIClient(config)

	accessToken,response, err := APIClient.AuthenticationApi.GetTokenForAuthGrantTypeUsingPOST1(auth, "client_credentials", nil)
	log.Printf("Request %v",response.Request)
	log.Printf("Response code %d",response.StatusCode)
	log.Printf("Request Body %v",response.Body)
	if err != nil {
		return nil, HandleCreateError("accesstoken using client ID and client secret",err)
	}

	securityCtx := security.NewOauthSecurityContext(accessToken.AccessToken)
	return securityCtx, nil
}
