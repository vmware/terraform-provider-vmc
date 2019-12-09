/* Copyright 2019 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var testAccProviders map[string]terraform.ResourceProvider

var testAccProvider *schema.Provider

func init() {
	testAccProvider = Provider().(*schema.Provider)
	testAccProviders = map[string]terraform.ResourceProvider{
		"vmc": testAccProvider,
	}
}

func TestProvider(t *testing.T) {
	if err := Provider().(*schema.Provider).InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}

}

func TestProvider_impl(t *testing.T) {
	var _ terraform.ResourceProvider = Provider()
}

func testAccPreCheck(t *testing.T) {
	if v := os.Getenv("REFRESH_TOKEN"); v == "" {
		t.Fatal("REFRESH_TOKEN must be set for acceptance tests")
	}
	if v := os.Getenv("ORG_ID"); v == "" {
		t.Fatal("ORG_ID must be set for acceptance tests")
	}
	if v := os.Getenv("TEST_SDDC_ID"); v == "" {
		t.Fatal("TEST_SDDC_ID must be set for acceptance tests")
	}
}
