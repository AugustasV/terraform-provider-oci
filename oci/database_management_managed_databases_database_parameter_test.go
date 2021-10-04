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
	managedDatabasesDatabaseParameterSingularDataSourceRepresentation = map[string]interface{}{
		"managed_database_id":        Representation{RepType: Required, Create: "ocid.database.testId"},
		"is_allowed_values_included": Representation{RepType: Optional, Create: `false`},
		"name":                       Representation{RepType: Required, Create: "open_cursors"},
		"source":                     Representation{RepType: Optional, Create: `CURRENT`},
	}

	managedDatabasesDatabaseParameterDataSourceRepresentation = map[string]interface{}{
		"managed_database_id":        Representation{RepType: Required, Create: "ocid.database.testId"},
		"is_allowed_values_included": Representation{RepType: Optional, Create: `false`},
		"name":                       Representation{RepType: Required, Create: `open_cursors`},
		"source":                     Representation{RepType: Optional, Create: `CURRENT`},
	}
)

// issue-routing-tag: database_management/default
func TestDatabaseManagementManagedDatabasesDatabaseParameterResource_basic(t *testing.T) {
	t.Skip("Skip this test till Database Management service provides a better way of testing this. It requires a live managed database instance")
	httpreplay.SetScenario("TestDatabaseManagementManagedDatabasesDatabaseParameterResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	datasourceName := "data.oci_database_management_managed_databases_database_parameters.test_managed_databases_database_parameters"
	singularDatasourceName := "data.oci_database_management_managed_databases_database_parameter.test_managed_databases_database_parameter"

	SaveConfigContent("", "", "", t)

	ResourceTest(t, nil, []resource.TestStep{
		// verify datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_database_management_managed_databases_database_parameters", "test_managed_databases_database_parameters", Required, Create, managedDatabasesDatabaseParameterDataSourceRepresentation) +
				compartmentIdVariableStr,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(datasourceName, "database_parameters_collection.#"),
				resource.TestCheckResourceAttrSet(datasourceName, "database_parameters_collection.0.database_name"),
				resource.TestCheckResourceAttrSet(datasourceName, "database_parameters_collection.0.database_sub_type"),
				resource.TestCheckResourceAttrSet(datasourceName, "database_parameters_collection.0.database_type"),
				resource.TestCheckResourceAttrSet(datasourceName, "database_parameters_collection.0.database_version"),
				resource.TestCheckResourceAttr(datasourceName, "database_parameters_collection.0.items.#", "1"),
			),
		},
		// verify singular datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_database_management_managed_databases_database_parameter", "test_managed_databases_database_parameter", Required, Create, managedDatabasesDatabaseParameterSingularDataSourceRepresentation) +
				compartmentIdVariableStr,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(singularDatasourceName, "database_name"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "database_sub_type"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "database_type"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "database_version"),
				resource.TestCheckResourceAttr(singularDatasourceName, "items.#", "1"),
			),
		},
	})
}
