/* Copyright 2019 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

// Package vmc provides helper methods that provides client.Connector, required to call VMC APIs.
package vmc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/vmware/vsphere-automation-sdk-go/runtime/core"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/security"
)

// NewClientConnectorByRefreshToken returns client connector by using OAuth authentication using Refresh Token.
func NewClientConnectorByRefreshToken(refreshToken, vmcURL, cspURL string,
	httpClient http.Client) (client.Connector, error) {

	if len(vmcURL) <= 0 {
		vmcURL = DefaultVMCServer
	}

	if len(cspURL) <= 0 {
		cspURL = DefaultCSPUrl +
			CSPRefreshUrlSuffix
	} else {
		cspURL = cspURL +
			CSPRefreshUrlSuffix
	}

	securityCtx, err := SecurityContextByRefreshToken(refreshToken, cspURL)
	if err != nil {
		return nil, err
	}

	connector := client.NewRestConnector(vmcURL, httpClient)
	connector.SetSecurityContext(securityCtx)

	return connector, nil
}

// SecurityContextByRefreshToken returns Security Context with access token that is received from CSP using Refresh Token by OAuth authentication scheme.
func SecurityContextByRefreshToken(refreshToken string, cspURL string) (core.SecurityContext, error) {
	payload := strings.NewReader("refresh_token=" + refreshToken)

	req, _ := http.NewRequest("POST", cspURL, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	if res.StatusCode != 200 {
		b, _ := ioutil.ReadAll(res.Body)
		return nil, fmt.Errorf("status code %d trying to get csp auth token. %s", res.StatusCode, string(b))
	}

	defer res.Body.Close()

	var jsondata map[string]interface{}
	json.NewDecoder(res.Body).Decode(&jsondata)

	var accessToken string
	if token, ok := jsondata["access_token"]; ok {
		if accessTokenStr, ok := token.(string); ok {
			accessToken = accessTokenStr
		} else {
			errMsg := fmt.Sprintf("Invalid type for access_token, expected string, actual %s", reflect.TypeOf(token).String())
			return nil, errors.New(errMsg)
		}
	} else {
		return nil, errors.New("CSP auth response doesn't contain access token")
	}

	securityCtx := security.NewOauthSecurityContext(accessToken)
	return securityCtx, nil
}
