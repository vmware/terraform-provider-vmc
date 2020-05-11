package vmc

import (
	"github.com/vmware/vsphere-automation-sdk-go/runtime/protocol/client"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs"
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
