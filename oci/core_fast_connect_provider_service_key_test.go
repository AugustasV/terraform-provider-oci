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
	fastConnectProviderServiceKeySingularDataSourceRepresentation = map[string]interface{}{
		"provider_service_id":       Representation{RepType: Required, Create: `${data.oci_core_fast_connect_provider_services.test_fast_connect_provider_services.fast_connect_provider_services.0.id}`},
		"provider_service_key_name": Representation{RepType: Required, Create: `d8f7a443-28c2-4dcf-996c-286351908c58`},
	}

	FastConnectProviderServiceKeyResourceConfig = GenerateDataSourceFromRepresentationMap("oci_core_fast_connect_provider_services", "test_fast_connect_provider_services", Required, Create, fastConnectProviderServiceDataSourceRepresentation)
)

// issue-routing-tag: core/default
func TestCoreFastConnectProviderServiceKeyResource_basic(t *testing.T) {
	if httpreplay.ShouldRetryImmediately() {
		t.Skip("FastConnect failed in dxterraformtest tenancy IAD region with a known issue: TER-1232")
	}

	httpreplay.SetScenario("TestCoreFastConnectProviderServiceKeyResource_basic")
	defer httpreplay.SaveScenario()

	config := ProviderTestConfig()

	compartmentId := GetEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	singularDatasourceName := "data.oci_core_fast_connect_provider_service_key.test_fast_connect_provider_service_key"

	SaveConfigContent("", "", "", t)

	ResourceTest(t, nil, []resource.TestStep{
		// verify singular datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_core_fast_connect_provider_service_key", "test_fast_connect_provider_service_key", Required, Create, fastConnectProviderServiceKeySingularDataSourceRepresentation) +
				compartmentIdVariableStr + FastConnectProviderServiceKeyResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(singularDatasourceName, "provider_service_id"),
				resource.TestCheckResourceAttr(singularDatasourceName, "provider_service_key_name", "d8f7a443-28c2-4dcf-996c-286351908c58"),

				resource.TestCheckResourceAttrSet(singularDatasourceName, "bandwidth_shape_name"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "name"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "peering_location"),
			),
		},
	})
}
