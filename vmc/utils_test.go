/* Copyright 2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/model"
	"testing"
)

func TestToHostInstanceType(t *testing.T) {
	var convertedHostInstanceType, err = toHostInstanceType(HostInstancetypeI3)
	if err != nil {
		t.Errorf("No errors expected!")
	}
	if convertedHostInstanceType != model.SddcConfig_HOST_INSTANCE_TYPE_I3_METAL {
		t.Errorf("Expected %s, but got %s",
			model.SddcConfig_HOST_INSTANCE_TYPE_I3_METAL, convertedHostInstanceType)
	}

	convertedHostInstanceType, err = toHostInstanceType(HostInstancetypeI3EN)
	if err != nil {
		t.Errorf("No errors expected!")
	}
	if convertedHostInstanceType != model.SddcConfig_HOST_INSTANCE_TYPE_I3EN_METAL {
		t.Errorf("Expected %s, but got %s",
			model.SddcConfig_HOST_INSTANCE_TYPE_I3EN_METAL, convertedHostInstanceType)
	}

	convertedHostInstanceType, err = toHostInstanceType(HostInstancetypeI4I)
	if err != nil {
		t.Errorf("No errors expected!")
	}
	if convertedHostInstanceType != model.SddcConfig_HOST_INSTANCE_TYPE_I4I_METAL {
		t.Errorf("Expected %s, but got %s",
			model.SddcConfig_HOST_INSTANCE_TYPE_I4I_METAL, convertedHostInstanceType)
	}

	convertedHostInstanceType, err = toHostInstanceType(HostInstancetypeR5)
	if err != nil {
		t.Errorf("No errors expected!")
	}
	if convertedHostInstanceType != model.SddcConfig_HOST_INSTANCE_TYPE_R5_METAL {
		t.Errorf("Expected %s, but got %s",
			model.SddcConfig_HOST_INSTANCE_TYPE_R5_METAL, convertedHostInstanceType)
	}

	_, err = toHostInstanceType("RandomString")
	if err == nil {
		t.Errorf("Error expected!")
	}
}
