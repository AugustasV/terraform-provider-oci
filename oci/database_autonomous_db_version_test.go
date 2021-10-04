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
	autonomousDbVersionDataSourceRepresentation = map[string]interface{}{
		"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id}`},
		"db_workload":    Representation{RepType: Optional, Create: `OLTP`},
	}

	AutonomousDbVersionResourceConfig = ""
)

// issue-routing-tag: database/default
func TestDatabaseAutonomousDbVersionResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestDatabaseAutonomousDbVersionResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	datasourceName := "data.oci_database_autonomous_db_versions.test_autonomous_db_versions"

	SaveConfigContent("", "", "", t)

	ResourceTest(t, nil, []resource.TestStep{
		// verify datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_database_autonomous_db_versions", "test_autonomous_db_versions", Required, Create, autonomousDbVersionDataSourceRepresentation) +
				compartmentIdVariableStr + AutonomousDbVersionResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),

				resource.TestCheckResourceAttrSet(datasourceName, "autonomous_db_versions.#"),
				resource.TestCheckResourceAttrSet(datasourceName, "autonomous_db_versions.0.db_workload"),
				resource.TestCheckResourceAttrSet(datasourceName, "autonomous_db_versions.0.details"),
				resource.TestCheckResourceAttrSet(datasourceName, "autonomous_db_versions.0.is_dedicated"),
				resource.TestCheckResourceAttrSet(datasourceName, "autonomous_db_versions.0.version"),
			),
		},

		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_database_autonomous_db_versions", "test_autonomous_db_versions", Optional, Create, autonomousDbVersionDataSourceRepresentation) +
				compartmentIdVariableStr + AutonomousDbVersionResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(datasourceName, "db_workload", "OLTP"),

				resource.TestCheckResourceAttrSet(datasourceName, "autonomous_db_versions.#"),
				resource.TestCheckResourceAttrSet(datasourceName, "autonomous_db_versions.0.db_workload"),
				resource.TestCheckResourceAttrSet(datasourceName, "autonomous_db_versions.0.details"),
				resource.TestCheckResourceAttrSet(datasourceName, "autonomous_db_versions.0.is_dedicated"),
				resource.TestCheckResourceAttrSet(datasourceName, "autonomous_db_versions.0.is_default_for_free"),
				resource.TestCheckResourceAttrSet(datasourceName, "autonomous_db_versions.0.is_default_for_paid"),
				resource.TestCheckResourceAttrSet(datasourceName, "autonomous_db_versions.0.is_free_tier_enabled"),
				resource.TestCheckResourceAttrSet(datasourceName, "autonomous_db_versions.0.is_paid_enabled"),
				resource.TestCheckResourceAttrSet(datasourceName, "autonomous_db_versions.0.version"),
			),
		},
	})
}
