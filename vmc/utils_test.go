/* Copyright 2022 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/model"
	"testing"
)

func TestToHostInstanceType(t *testing.T) {
	type result struct {
		converted string
		err       error
	}
	type test struct {
		input string
		want  result
	}

	tests := []test{
		{input: HostInstancetypeI3, want: result{converted: model.SddcConfig_HOST_INSTANCE_TYPE_I3_METAL, err: nil}},
		{input: HostInstancetypeI3EN, want: result{converted: model.SddcConfig_HOST_INSTANCE_TYPE_I3EN_METAL, err: nil}},
		{input: HostInstancetypeI4I, want: result{converted: model.SddcConfig_HOST_INSTANCE_TYPE_I4I_METAL, err: nil}},
		{input: HostInstancetypeR5, want: result{converted: model.SddcConfig_HOST_INSTANCE_TYPE_R5_METAL, err: nil}},
		{input: "RandomString", want: result{converted: "", err: fmt.Errorf("unknown host instance type: RandomString")}},
	}

	for _, testCase := range tests {
		got, err := toHostInstanceType(testCase.input)
		assert.Equal(t, got, testCase.want.converted)
		assert.Equal(t, err, testCase.want.err)
	}
}
