/* Copyright 2019 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var testAccProviders map[string]*schema.Provider

var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider()
	testAccProviders = map[string]*schema.Provider{
		"vmc": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}

}

func TestProvider_impl(t *testing.T) {
	var _ *schema.Provider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv(APIToken); v == "" {
		t.Fatal(APIToken + " must be set for acceptance tests")
	}
	if v := os.Getenv(OrgID); v == "" {
		t.Fatal(OrgID + " must be set for acceptance tests")
	}
	if v := os.Getenv(OrgDisplayName); v == "" {
		t.Fatal(OrgDisplayName + " must be set for acceptance tests")
	}
	if v := os.Getenv(TestSDDCId); v == "" {
		t.Fatal(TestSDDCId + " must be set for acceptance tests")
	}
	if v := os.Getenv(TestSDDCName); v == "" {
		t.Fatal(TestSDDCName + " must be set for acceptance tests")
	}
	if v := os.Getenv(AWSAccountNumber); v == "" {
		t.Fatal(AWSAccountNumber + " must be set for acceptance tests")
	}
	if v := os.Getenv(NSXTReverseProxyUrl); v == "" {
		t.Fatal(NSXTReverseProxyUrl + " must be set for acceptance tests")
	}
}
