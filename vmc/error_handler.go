/* Copyright 2020 VMware, Inc.
   SPDX-License-Identifier: MPL-2.0 */

package vmc

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std"
	e "github.com/vmware/vsphere-automation-sdk-go/lib/vapi/std/errors"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/bindings"
	"github.com/vmware/vsphere-automation-sdk-go/runtime/data"
	"github.com/vmware/vsphere-automation-sdk-go/services/vmc/model"
	"log"
)

func printAPIError(apiError model.ErrorResponse) []string {
	var detailedErrorResponse []string
	if len(apiError.ErrorMessages) > 0 && apiError.ErrorCode != "0" {
		var details string
		for i := 0; i < len(apiError.ErrorMessages); i++ {
			details += fmt.Sprintf("%s (code)%v ", apiError.ErrorMessages[i], apiError.ErrorCode[i])
		}

		detailedErrorResponse = append(detailedErrorResponse, details)
		return detailedErrorResponse
	}

	if len(apiError.ErrorMessages) > 0 {
		return apiError.ErrorMessages
	}

	if apiError.ErrorCode != "" {
		detailedErrorResponse = append(detailedErrorResponse, fmt.Sprintf("Error code :  %s ", apiError.ErrorCode))
		return detailedErrorResponse
	}

	return detailedErrorResponse
}

func logVapiErrorData(message string, vAPIMessages []std.LocalizableMessage, vapiType *e.ErrorType, apiErrorDataValue *data.StructValue) error {

	if apiErrorDataValue == nil {
		if len(vAPIMessages) > 0 {
			return fmt.Errorf("%s (%s)", message, vAPIMessages[0].DefaultMessage)
		}
		if vapiType != nil {
			return fmt.Errorf("%s (%s)", message, *vapiType)
		}

		return fmt.Errorf("%s (no additional details provided)", message)
	}

	var typeConverter = bindings.NewTypeConverter()
	typeConverter.SetMode(bindings.REST)
	data, err := typeConverter.ConvertToGolang(apiErrorDataValue, model.ErrorResponseBindingType())

	if err != nil {
		log.Printf("[ERROR]: Failed to extract error details: %s", err)
		if len(vAPIMessages) > 0 {
			return fmt.Errorf("%s (%s)", message, vAPIMessages[0].DefaultMessage)
		}
		// error type is the only piece of info we have here
		if vapiType != nil {
			return fmt.Errorf("%s (%s)", message, *vapiType)
		}

		return fmt.Errorf("%s (no additional details provided)", message)
	}

	apiError := data.(model.ErrorResponse)

	details := fmt.Sprintf(" %s: %s", message, printAPIError(apiError))
	log.Printf("[ERROR]: %s", details)
	return fmt.Errorf(details)
}

func logAPIError(message string, err error) error {
	if vapiError, ok := err.(e.InvalidRequest); ok {
		return logVapiErrorData(message, vapiError.Messages, vapiError.ErrorType, vapiError.Data)
	}
	if vapiError, ok := err.(e.NotFound); ok {
		return logVapiErrorData(message, vapiError.Messages, vapiError.ErrorType, vapiError.Data)
	}
	if vapiError, ok := err.(e.Unauthorized); ok {
		return logVapiErrorData(message, vapiError.Messages, vapiError.ErrorType, vapiError.Data)
	}
	if vapiError, ok := err.(e.Unauthenticated); ok {
		return logVapiErrorData(message, vapiError.Messages, vapiError.ErrorType, vapiError.Data)
	}
	if vapiError, ok := err.(e.InternalServerError); ok {
		return logVapiErrorData(message, vapiError.Messages, vapiError.ErrorType, vapiError.Data)
	}
	if vapiError, ok := err.(e.ServiceUnavailable); ok {
		return logVapiErrorData(message, vapiError.Messages, vapiError.ErrorType, vapiError.Data)
	}
	if vapiError, ok := err.(e.AlreadyExists); ok {
		return logVapiErrorData(message, vapiError.Messages, vapiError.ErrorType, vapiError.Data)
	}
	if vapiError, ok := err.(e.AlreadyInDesiredState); ok {
		return logVapiErrorData(message, vapiError.Messages, vapiError.ErrorType, vapiError.Data)
	}
	return err
}

func isNotFoundError(err error) bool {
	if _, ok := err.(e.NotFound); ok {
		return true
	}

	return false
}

func HandleCreateError(resourceType string, err error) error {
	msg := fmt.Sprintf("Failed to create %s", resourceType)
	return logAPIError(msg, err)
}

func HandleUpdateError(resourceType string, err error) error {
	msg := fmt.Sprintf("Failed to update %s", resourceType)
	return logAPIError(msg, err)
}

func HandleListError(resourceType string, err error) error {
	msg := fmt.Sprintf("Failed to read %s", resourceType)
	return logAPIError(msg, err)
}

func HandleReadError(d *schema.ResourceData, resourceType string, resourceID string, err error) error {
	msg := fmt.Sprintf("Failed to read %s %s", resourceType, resourceID)
	if isNotFoundError(err) {
		d.SetId("")
		log.Printf(msg)
		return nil
	}
	return logAPIError(msg, err)
}

func HandleDataSourceReadError(d *schema.ResourceData, resourceType string, err error) error {
	msg := fmt.Sprintf("Failed to read %s ", resourceType)
	return logAPIError(msg, err)
}

func HandleDeleteError(resourceType string, resourceID string, err error) error {
	if isNotFoundError(err) {
		log.Printf("[WARNING] %s %s not found on backend", resourceType, resourceID)
		// We don't want to fail apply on this
		return nil
	}
	msg := fmt.Sprintf("Failed to delete %s %s", resourceType, resourceID)
	return logAPIError(msg, err)
}
