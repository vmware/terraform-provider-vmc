/* Copyright 2019 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/model"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/orgs/sddcs"
	"log"
	"strings"
	"time"
)

func resourcePublicIP() *schema.Resource {
	return &schema.Resource{
		Create: resourcePublicIPCreate,
		Read:   resourcePublicIPRead,
		Update: resourcePublicIPUpdate,
		Delete: resourcePublicIPDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(10 * time.Minute),
			Update: schema.DefaultTimeout(10 * time.Minute),
			Delete: schema.DefaultTimeout(10 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"org_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Organization identifier",
			},
			"sddc_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Sddc Identifier",
			},
			"allocation_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "IP Allocation ID",
			},
			"public_ip": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Allocated Public IP",
			},
			"private_ip": {
				Type:        schema.TypeString,
				Elem:        &schema.Schema{Type: schema.TypeString},
				Optional:    true,
				Description: "Workload VM private IP",
			},
			"name": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Workload VM name",
			},
			"dnat_rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "DNAT rule ID",
			},
			"snat_rule_id": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "SNAT rule ID",
			},
		},
	}
}

func resourcePublicIPCreate(d *schema.ResourceData, m interface{}) error {
	connector := (m.(*ConnectorWrapper)).Connector

	orgID := d.Get("org_id").(string)
	sddcID := d.Get("sddc_id").(string)

	privateIP := d.Get("private_ip").(string)
	workloadName := d.Get("name").(string)
	publicIPsClient := sddcs.NewDefaultPublicipsClient(connector)

	var sddcAllocatePublicIpSpec = &model.SddcAllocatePublicIpSpec{
		Count:      1,
		PrivateIps: []string{privateIP},
		Names:      []string{workloadName},
	}

	// Create Public IP
	task, err := publicIPsClient.Create(orgID, sddcID, *sddcAllocatePublicIpSpec)
	if err != nil {
		return fmt.Errorf("error while creating public IP : %v", err)
	}

	tasksClient := orgs.NewDefaultTasksClient(connector)

	return resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		task, err := tasksClient.Get(orgID, task.Id)

		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("error describing instance: %s", err))
		}
		if *task.Status != "FINISHED" {
			log.Print("Task not finished yet")
			return resource.RetryableError(fmt.Errorf("expected instance to be created but was in state %s", *task.Status))
		} else {
			publicIPClient := sddcs.NewDefaultPublicipsClient(connector)
			publicIPs, err := publicIPClient.List(orgID, sddcID)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("error while getting list of public IPs for SDDC %s: %v", d.Get("sddc_id").(string), err))
			}
			for i := 0; i < len(publicIPs); i++ {
				singleVal := publicIPs[i]
				if d.Get("private_ip").(string) == *(singleVal.AssociatedPrivateIp) {
					d.SetId(*(singleVal.AllocationId))
					break
				}
			}
			if d.Id() == "" {
				return resource.NonRetryableError(fmt.Errorf("error while getting the allocationID %v", err))
			}
			return resource.NonRetryableError(resourcePublicIPRead(d, m))
		}
	})
}

func resourcePublicIPRead(d *schema.ResourceData, m interface{}) error {

	connector := (m.(*ConnectorWrapper)).Connector
	publicIPClient := sddcs.NewDefaultPublicipsClient(connector)

	orgID := d.Get("org_id").(string)
	sddcID := d.Get("sddc_id").(string)
	allocationID := d.Id()
	//check if the SDDC exists
	sddc, err := getSDDC(connector, orgID, sddcID)
	if err != nil {
		if err.Error() == errors.NewNotFound().Error() {
			log.Printf("Can't get Public IP: The associated SDDC with ID %s not found", sddcID)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Can't get Public IP: Error while getting the associated SDDC with ID %s,%v", sddcID, err)
	}

	if *sddc.SddcState == "DELETED" {
		log.Printf("Can't get public IP: the associated SDDC with ID %s is already deleted", sddc.Id)
		d.SetId("")
		return nil
	}

	publicIP, err := publicIPClient.Get(orgID, sddcID, allocationID)
	if err != nil {
		if err.Error() == errors.NewNotFound().Error() {
			log.Printf("Public IP with allocation ID %s not found", allocationID)
			d.SetId("")
			return nil
		}
		return fmt.Errorf("Error while getting public IP details for %s: %v", allocationID, err)
	}

	d.SetId(*publicIP.AllocationId)
	d.Set("public_ip", publicIP.PublicIp)
	d.Set("name", publicIP.Name)
	d.Set("private_ip", publicIP.AssociatedPrivateIp)
	d.Set("dnat_rule_id", publicIP.DnatRuleId)
	d.Set("snat_rule_id", publicIP.SnatRuleId)
	return nil

}

func resourcePublicIPDelete(d *schema.ResourceData, m interface{}) error {

	connector := (m.(*ConnectorWrapper)).Connector
	publicIPClient := sddcs.NewDefaultPublicipsClient(connector)

	allocationID := d.Id()
	orgID := d.Get("org_id").(string)
	sddcID := d.Get("sddc_id").(string)
	publicIP := d.Get("public_ip").(string)

	task, err := publicIPClient.Delete(orgID, sddcID, allocationID)
	if err != nil {
		if err.Error() == errors.NewInvalidRequest().Error() {
			log.Printf("Can't Delete : Public IP  %s not found or already deleted %v", publicIP, err)
			return nil
		}
		return fmt.Errorf("Error while deleting public IP %s: %v", publicIP, err)
	}
	tasksClient := orgs.NewDefaultTasksClient(connector)
	return resource.Retry(d.Timeout(schema.TimeoutDelete), func() *resource.RetryError {
		task, err := tasksClient.Get(orgID, task.Id)
		if err != nil {
			return resource.NonRetryableError(fmt.Errorf("Error while deleting public IP %s: %v", publicIP, err))
		}
		if task.ErrorMessage != nil && strings.Contains(*task.ErrorMessage, "Entity is not found for Id "+allocationID) {
			log.Printf("Can't Delete : Public IP  %s not found or already deleted %v", publicIP, *task.ErrorMessage)
			return resource.NonRetryableError(nil)
		}
		if *task.Status != "FINISHED" {
			return resource.RetryableError(fmt.Errorf("Expected instance to be deleted but was in state %s", *task.Status))
		}
		return resource.NonRetryableError(nil)
	})
}

func resourcePublicIPUpdate(d *schema.ResourceData, m interface{}) error {
	connector := (m.(*ConnectorWrapper)).Connector
	publicIPClient := sddcs.NewDefaultPublicipsClient(connector)
	allocationID := d.Id()
	orgID := d.Get("org_id").(string)
	sddcID := d.Get("sddc_id").(string)
	publicIPName := d.Get("name").(string)
	associatedPrivateIP := d.Get("private_ip").(string)
	publicIP := d.Get("public_ip").(string)

	if d.HasChange("private_ip") {

		if d.Get("private_ip") == "" {
			//detach privateIP case
			newSDDCPublicIP := model.SddcPublicIp{
				PublicIp: publicIP,
				Name:     &publicIPName,
			}
			task, err := publicIPClient.Update(orgID, sddcID, allocationID, "detach", newSDDCPublicIP)
			if err != nil {
				return fmt.Errorf("error while detaching the public ip: %v", err)
			}
			tasksClient := orgs.NewDefaultTasksClient(connector)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				task, err := tasksClient.Get(orgID, task.Id)
				if err != nil {
					return resource.NonRetryableError(fmt.Errorf("Error while waiting for task sddc %s: %v", task.Id, err))
				}
				if *task.Status != "FINISHED" {
					return resource.RetryableError(fmt.Errorf("Expected IP to be detached but was in state %s", *task.Status))
				}
				return resource.NonRetryableError(resourcePublicIPRead(d, m))
			})
			if err != nil {
				return err
			}

		} else {
			//reattach privateIP case
			newSDDCPublicIP := model.SddcPublicIp{
				PublicIp:            publicIP,
				AssociatedPrivateIp: &associatedPrivateIP,
				Name:                &publicIPName,
			}
			task, err := publicIPClient.Update(orgID, sddcID, allocationID, "reattach", newSDDCPublicIP)
			if err != nil {
				return fmt.Errorf("error while reattaching the public IP : %v", err)
			}
			tasksClient := orgs.NewDefaultTasksClient(connector)
			err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
				task, err := tasksClient.Get(orgID, task.Id)
				if err != nil {
					return resource.NonRetryableError(fmt.Errorf("error while waiting for task sddc %s: %v", task.Id, err))
				}
				if *task.Status != "FINISHED" {
					return resource.RetryableError(fmt.Errorf("expected IP to be reattached but was in state %s", *task.Status))
				}
				return resource.NonRetryableError(resourcePublicIPRead(d, m))
			})
			if err != nil {
				return err
			}
		}

	} else if d.HasChange("name") {
		//rename case
		newSDDCPublicIP := model.SddcPublicIp{
			Name:                &publicIPName,
			AssociatedPrivateIp: &associatedPrivateIP,
		}
		task, err := publicIPClient.Update(orgID, sddcID, allocationID, "rename", newSDDCPublicIP)

		if err != nil {
			return fmt.Errorf("error while updating public IP for rename action type  : %v", err)
		}

		tasksClient := orgs.NewDefaultTasksClient(connector)
		err = resource.Retry(d.Timeout(schema.TimeoutUpdate), func() *resource.RetryError {
			task, err := tasksClient.Get(orgID, task.Id)
			if err != nil {
				return resource.NonRetryableError(fmt.Errorf("error while waiting for task sddc %s: %v", task.Id, err))
			}
			if *task.Status != "FINISHED" {
				return resource.RetryableError(fmt.Errorf("expected IP to be renamed but was in state %s", *task.Status))
			}
			return resource.NonRetryableError(resourcePublicIPRead(d, m))
		})
		if err != nil {
			return err
		}
	}
	return resourcePublicIPRead(d, m)
}
