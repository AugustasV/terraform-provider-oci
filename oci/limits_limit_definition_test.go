// Copyright (c) 2017, 2021, Oracle and/or its affiliates. All rights reserved.
// Licensed under the Mozilla Public License v2.0

package oci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	limitDefinitionDataSourceRepresentation = map[string]interface{}{
		"compartment_id": Representation{RepType: Required, Create: `${var.tenancy_ocid}`},
		"name":           Representation{RepType: Optional, Create: `custom-image-count`},
		"service_name":   Representation{RepType: Optional, Create: `${data.oci_limits_services.test_services.services.0.name}`},
	}

	LimitDefinitionResourceConfig = GenerateDataSourceFromRepresentationMap("oci_limits_services", "test_services", Required, Create, limitsServiceDataSourceRepresentation)
)

// issue-routing-tag: limits/default
func TestLimitsLimitDefinitionResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestLimitsLimitDefinitionResource_basic")
	defer httpreplay.SaveScenario()

	config := ProviderTestConfig()

	compartmentId := GetEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)
	tenancyId := GetEnvSettingWithBlankDefault("tenancy_ocid")

	datasourceName := "data.oci_limits_limit_definitions.test_limit_definitions"

	SaveConfigContent("", "", "", t)

	ResourceTest(t, nil, []resource.TestStep{
		// verify datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_limits_limit_definitions", "test_limit_definitions", Required, Create, limitDefinitionDataSourceRepresentation) +
				compartmentIdVariableStr + LimitDefinitionResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(datasourceName, "compartment_id", tenancyId),
				resource.TestCheckResourceAttrSet(datasourceName, "limit_definitions.#"),
				resource.TestCheckResourceAttrSet(datasourceName, "limit_definitions.0.are_quotas_supported"),
				resource.TestCheckResourceAttrSet(datasourceName, "limit_definitions.0.description"),
				resource.TestCheckResourceAttrSet(datasourceName, "limit_definitions.0.is_deprecated"),
				resource.TestCheckResourceAttrSet(datasourceName, "limit_definitions.0.is_dynamic"),
				resource.TestCheckResourceAttrSet(datasourceName, "limit_definitions.0.is_eligible_for_limit_increase"),
				resource.TestCheckResourceAttrSet(datasourceName, "limit_definitions.0.is_resource_availability_supported"),
				resource.TestCheckResourceAttrSet(datasourceName, "limit_definitions.0.name"),
				resource.TestCheckResourceAttrSet(datasourceName, "limit_definitions.0.scope_type"),
				resource.TestCheckResourceAttrSet(datasourceName, "limit_definitions.0.service_name"),
			),
		},
	})
}
