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
	edgeSubnetDataSourceRepresentation = map[string]interface{}{}

	EdgeSubnetResourceConfig = ""
)

// issue-routing-tag: waas/default
func TestWaasEdgeSubnetResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestWaasEdgeSubnetResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	datasourceName := "data.oci_waas_edge_subnets.test_edge_subnets"

	SaveConfigContent("", "", "", t)

	ResourceTest(t, nil, []resource.TestStep{
		// verify datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_waas_edge_subnets", "test_edge_subnets", Required, Create, edgeSubnetDataSourceRepresentation) +
				compartmentIdVariableStr + EdgeSubnetResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(

				resource.TestCheckResourceAttrSet(datasourceName, "edge_subnets.#"),
				resource.TestCheckResourceAttrSet(datasourceName, "edge_subnets.0.cidr"),
				resource.TestCheckResourceAttrSet(datasourceName, "edge_subnets.0.region"),
				resource.TestCheckResourceAttrSet(datasourceName, "edge_subnets.0.time_modified"),
			),
		},
	})
}
