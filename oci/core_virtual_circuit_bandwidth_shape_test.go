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
	virtualCircuitBandwidthShapeDataSourceRepresentation = map[string]interface{}{
		"provider_service_id": Representation{RepType: Required, Create: `${data.oci_core_fast_connect_provider_services.test_fast_connect_provider_services.fast_connect_provider_services.0.id}`},
	}

	VirtualCircuitBandwidthShapeResourceConfig = GenerateDataSourceFromRepresentationMap("oci_core_fast_connect_provider_services", "test_fast_connect_provider_services", Required, Create, fastConnectProviderServiceDataSourceRepresentation)
)

// issue-routing-tag: core/default
func TestCoreVirtualCircuitBandwidthShapeResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestCoreVirtualCircuitBandwidthShapeResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	datasourceName := "data.oci_core_virtual_circuit_bandwidth_shapes.test_virtual_circuit_bandwidth_shapes"

	SaveConfigContent("", "", "", t)

	ResourceTest(t, nil, []resource.TestStep{
		// verify datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_core_virtual_circuit_bandwidth_shapes", "test_virtual_circuit_bandwidth_shapes", Required, Create, virtualCircuitBandwidthShapeDataSourceRepresentation) +
				compartmentIdVariableStr + VirtualCircuitBandwidthShapeResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(

				resource.TestCheckResourceAttrSet(datasourceName, "virtual_circuit_bandwidth_shapes.#"),
				resource.TestCheckResourceAttrSet(datasourceName, "virtual_circuit_bandwidth_shapes.0.bandwidth_in_mbps"),
				resource.TestCheckResourceAttrSet(datasourceName, "virtual_circuit_bandwidth_shapes.0.name"),
			),
		},
	})
}
