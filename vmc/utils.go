/* Copyright 2020 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	uuid "github.com/satori/go.uuid"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs"
	"net/url"
)

var storageCapacityMap = map[string]int64{
	"15TB": 15003,
	"20TB": 20004,
	"25TB": 25005,
	"30TB": 30006,
	"35TB": 35007,
}

func GetSDDC(connector client.Connector, orgID string, sddcID string) (model.Sddc, error) {
	sddcClient := orgs.NewDefaultSddcsClient(connector)
	sddc, err := sddcClient.Get(orgID, sddcID)
	return sddc, err
}

func ConvertStorageCapacitytoInt(s string) int64 {
	storageCapacity := storageCapacityMap[s]
	return storageCapacity
}

// Mapping for deployment_type field
// During refresh/import state, return value of VMC API should be converted to uppercamel case in terraform
// to maintain consistency
func ConvertDeployType(s string) string {
	if s == "SINGLE_AZ" {
		return "SingleAZ"
	} else if s == "MULTI_AZ" {
		return "MultiAZ"
	} else {
		return ""
	}
}

func IsValidUUID(u string) error {
	_, err := uuid.FromString(u)
	if err != nil {
		return err
	}
	return nil
}

func IsValidURL(s string) error {
	_, err := url.ParseRequestURI(s)
	if err != nil {
		return err
	}
	return nil
}
