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
	mysqlVersionDataSourceRepresentation = map[string]interface{}{
		"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id}`},
	}

	MysqlVersionResourceConfig = GenerateDataSourceFromRepresentationMap("oci_mysql_mysql_versions", "test_mysql_versions", Required, Create, mysqlVersionDataSourceRepresentation)
)

// issue-routing-tag: mysql/default
func TestMysqlMysqlVersionResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestMysqlMysqlVersionResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	datasourceName := "data.oci_mysql_mysql_versions.test_mysql_versions"

	SaveConfigContent("", "", "", t)

	ResourceTest(t, nil, []resource.TestStep{
		// verify datasource
		{
			Config: config + compartmentIdVariableStr + MysqlVersionResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),

				resource.TestCheckResourceAttrSet(datasourceName, "versions.#"),
				resource.TestCheckResourceAttrSet(datasourceName, "versions.0.version_family"),
			),
		},
	})
}
