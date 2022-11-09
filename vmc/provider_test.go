/* Copyright 2019-2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"github.com/vmware/terraform-provider-vmc/vmc/constants"
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
	if v := os.Getenv(constants.APIToken); v == "" {
		t.Fatal(constants.APIToken + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.OrgID); v == "" {
		t.Fatal(constants.OrgID + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.OrgDisplayName); v == "" {
		t.Fatal(constants.OrgDisplayName + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.TestSddcID); v == "" {
		t.Fatal(constants.TestSddcID + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.TestSddcName); v == "" {
		t.Fatal(constants.TestSddcName + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.AwsAccountNumber); v == "" {
		t.Fatal(constants.AwsAccountNumber + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.NsxtReverseProxyURL); v == "" {
		t.Fatal(constants.NsxtReverseProxyURL + " must be set for acceptance tests")
	}
}

// testAccPreCheckZerocloud this function validates a smaller set ot
// environment variables needed for lightweight E2E testing using
// the Zerocloud SDDC cloud provider option
func testAccPreCheckZerocloud(t *testing.T) {
	if v := os.Getenv(constants.VmcURL); v == "" {
		t.Fatal(constants.VmcURL + " must be set for Zerocloud acceptance tests")
	}
	if v := os.Getenv(constants.CspURL); v == "" {
		t.Fatal(constants.CspURL + " must be set for Zerocloud acceptance tests")
	}
	if v := os.Getenv(constants.APIToken); v == "" {
		t.Fatal(constants.APIToken + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.OrgID); v == "" {
		t.Fatal(constants.OrgID + " must be set for acceptance tests")
	}
	if v := os.Getenv(constants.AwsAccountNumber); v == "" {
		t.Fatal(constants.AwsAccountNumber + " must be set for acceptance tests")
	}
}
