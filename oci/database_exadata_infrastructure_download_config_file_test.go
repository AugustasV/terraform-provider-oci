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
	exadataInfrastructureDownloadConfigFileSingularDataSourceRepresentation = map[string]interface{}{
		"exadata_infrastructure_id": Representation{RepType: Required, Create: `${oci_database_exadata_infrastructure.test_exadata_infrastructure.id}`},
		"base64_encode_content":     Representation{RepType: Optional, Create: `true`},
	}

	ExadataInfrastructureDownloadConfigFileResourceConfig = GenerateResourceFromRepresentationMap("oci_database_exadata_infrastructure", "test_exadata_infrastructure", Required, Create, exadataInfrastructureRepresentation)
)

// issue-routing-tag: database/ExaCC
func TestDatabaseExadataInfrastructureDownloadConfigFileResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestDatabaseExadataInfrastructureDownloadConfigFileResource_basic")
	defer httpreplay.SaveScenario()

	config := ProviderTestConfig()

	compartmentId := GetEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	singularDatasourceName := "data.oci_database_exadata_infrastructure_download_config_file.test_exadata_infrastructure_download_config_file"

	SaveConfigContent("", "", "", t)

	ResourceTest(t, nil, []resource.TestStep{
		// verify singular datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_database_exadata_infrastructure_download_config_file", "test_exadata_infrastructure_download_config_file", Required, Create, exadataInfrastructureDownloadConfigFileSingularDataSourceRepresentation) +
				compartmentIdVariableStr + ExadataInfrastructureDownloadConfigFileResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(singularDatasourceName, "exadata_infrastructure_id"),
				resource.TestCheckResourceAttr(singularDatasourceName, "base64_encode_content", "false"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "content"),
			),
		},

		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_database_exadata_infrastructure_download_config_file", "test_exadata_infrastructure_download_config_file", Optional, Create, exadataInfrastructureDownloadConfigFileSingularDataSourceRepresentation) +
				compartmentIdVariableStr + ExadataInfrastructureDownloadConfigFileResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(singularDatasourceName, "exadata_infrastructure_id"),
				resource.TestCheckResourceAttr(singularDatasourceName, "base64_encode_content", "true"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "content"),
			),
		},
	})
}
