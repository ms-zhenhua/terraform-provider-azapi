package services

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Azure/terraform-provider-azapi/internal/clients"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"time"

	"log"
	"strings"

	"github.com/Azure/terraform-provider-azapi/internal/azure"
	"github.com/Azure/terraform-provider-azapi/internal/azure/types"
	"github.com/Azure/terraform-provider-azapi/utils"
)

func preflightValidation(meta interface{}, body map[string]interface{}) error {
	client := meta.(*clients.Client).ResourceClient
	ctx, cancel := context.WithTimeout(meta.(*clients.Client).StopContext, *schema.DefaultTimeout(3 * time.Minute))
	defer cancel()

	_, err := client.PreflightValidate(ctx, body)
	return err
}

func schemaValidation(azureResourceType, apiVersion string, resourceDef *types.ResourceType, body interface{}) error {
	log.Printf("[INFO] prepare validation for resource type: %s, api-version: %s", azureResourceType, apiVersion)
	versions := azure.GetApiVersions(azureResourceType)
	if len(versions) == 0 {
		return schemaValidationError(fmt.Sprintf("the `type` is invalid.\n resource type %s can't be found.\n", azureResourceType))
	}
	isVersionValid := false
	for _, version := range versions {
		if version == apiVersion {
			isVersionValid = true
			break
		}
	}
	if !isVersionValid {
		return schemaValidationError(fmt.Sprintf("the `type`'s api-version is invalid.\n The supported versions are [%s].\n", strings.Join(versions, ", ")))
	}

	if resourceDef != nil {
		errors := (*resourceDef).Validate(utils.NormalizeObject(body), "")
		if len(errors) != 0 {
			errorMsg := "the `body` is invalid:\n"
			for _, err := range errors {
				errorMsg += fmt.Sprintf("%s\n", err.Error())
			}
			return schemaValidationError(errorMsg)
		}
	}
	return nil
}

func schemaValidationError(detail string) error {
	return fmt.Errorf("embedded schema validation failed: %s You can try to update `azapi` provider to "+
		"the latest version or disable the validation using the feature flag `schema_validation_enabled = false` "+
		"within the resource block", detail)
}

func isResourceHasProperty(resourceDef *types.ResourceType, property string) bool {
	if resourceDef == nil || resourceDef.Body == nil || resourceDef.Body.Type == nil {
		return false
	}
	objectType, ok := (*resourceDef.Body.Type).(*types.ObjectType)
	if !ok {
		return false
	}
	if prop, ok := objectType.Properties[property]; ok {
		if !prop.IsReadOnly() {
			return true
		}
	}
	return false
}

func flattenOutput(responseBody interface{}, paths []interface{}) string {
	for _, path := range paths {
		if path == "*" {
			if v, ok := responseBody.(string); ok {
				return v
			}
			outputJson, _ := json.Marshal(responseBody)
			return string(outputJson)
		}
	}

	var output interface{}
	output = make(map[string]interface{})
	for _, path := range paths {
		part := utils.ExtractObject(responseBody, path.(string))
		if part == nil {
			continue
		}
		output = utils.GetMergedJson(output, part)
	}
	outputJson, _ := json.Marshal(output)
	return string(outputJson)
}
