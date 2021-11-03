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
	autonomousContainerPatchDataSourceRepresentation = map[string]interface{}{
		"autonomous_container_database_id": Representation{RepType: Required, Create: `${oci_database_autonomous_container_database.test_autonomous_container_database.id}`},
		"compartment_id":                   Representation{RepType: Required, Create: `${var.compartment_id}`},
	}

	AutonomousContainerPatchResourceConfig = GenerateResourceFromRepresentationMap("oci_database_autonomous_container_database", "test_autonomous_container_database", Required, Create, autonomousContainerDatabaseRepresentation) +
		AutonomousExadataInfrastructureResourceConfig
)

// issue-routing-tag: database/default
func TestDatabaseAutonomousContainerPatchResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestDatabaseAutonomousContainerPatchResource_basic")
	defer httpreplay.SaveScenario()

	config := ProviderTestConfig()

	compartmentId := GetEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	datasourceName := "data.oci_database_autonomous_container_patches.test_autonomous_container_patches"

	SaveConfigContent("", "", "", t)

	ResourceTest(t, nil, []resource.TestStep{
		// verify datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_database_autonomous_container_patches", "test_autonomous_container_patches", Required, Create, autonomousContainerPatchDataSourceRepresentation) +
				compartmentIdVariableStr + AutonomousContainerPatchResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(datasourceName, "autonomous_container_database_id"),
				resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),

				resource.TestCheckResourceAttrSet(datasourceName, "autonomous_patches.#"),
				resource.TestCheckResourceAttrSet(datasourceName, "autonomous_patches.0.description"),
				resource.TestCheckResourceAttrSet(datasourceName, "autonomous_patches.0.id"),
				resource.TestCheckResourceAttrSet(datasourceName, "autonomous_patches.0.patch_model"),
				resource.TestCheckResourceAttrSet(datasourceName, "autonomous_patches.0.quarter"),
				resource.TestCheckResourceAttrSet(datasourceName, "autonomous_patches.0.type"),
				resource.TestCheckResourceAttrSet(datasourceName, "autonomous_patches.0.version"),
				resource.TestCheckResourceAttrSet(datasourceName, "autonomous_patches.0.year"),
			),
		},
	})
}
