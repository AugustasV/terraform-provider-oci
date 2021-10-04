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
	publicVantagePointSingularDataSourceRepresentation = map[string]interface{}{
		"apm_domain_id": Representation{RepType: Required, Create: `${oci_apm_apm_domain.test_apm_domain.id}`},
		"display_name":  Representation{RepType: Optional, Create: `US East (Ashburn)`},
		"name":          Representation{RepType: Optional, Create: `OraclePublic-us-ashburn-1`},
	}

	publicVantagePointDataSourceRepresentation = map[string]interface{}{
		"apm_domain_id": Representation{RepType: Required, Create: `${oci_apm_apm_domain.test_apm_domain.id}`},
		"display_name":  Representation{RepType: Optional, Create: `US East (Ashburn)`},
		"name":          Representation{RepType: Optional, Create: `OraclePublic-us-ashburn-1`},
	}

	PublicVantagePointResourceConfig = ""
)

// issue-routing-tag: apm_synthetics/default
func TestApmSyntheticsPublicVantagePointResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestApmSyntheticsPublicVantagePointResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	datasourceName := "data.oci_apm_synthetics_public_vantage_points.test_public_vantage_points"
	singularDatasourceName := "data.oci_apm_synthetics_public_vantage_point.test_public_vantage_point"

	SaveConfigContent("", "", "", t)

	ResourceTest(t, nil, []resource.TestStep{
		// verify datasource
		{
			Config: config + GenerateResourceFromRepresentationMap("oci_apm_apm_domain", "test_apm_domain", Required, Create, apmDomainRepresentation) +
				GenerateDataSourceFromRepresentationMap("oci_apm_synthetics_public_vantage_points", "test_public_vantage_points", Optional, Create, publicVantagePointDataSourceRepresentation) +
				compartmentIdVariableStr + PublicVantagePointResourceConfig,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(datasourceName, "apm_domain_id"),
				resource.TestCheckResourceAttr(datasourceName, "name", "OraclePublic-us-ashburn-1"),
				resource.TestCheckResourceAttr(datasourceName, "display_name", "US East (Ashburn)"),

				resource.TestCheckResourceAttrSet(datasourceName, "public_vantage_point_collection.#"),
			),
		},
		// verify singular datasource
		{
			Config: config + GenerateResourceFromRepresentationMap("oci_apm_apm_domain", "test_apm_domain", Required, Create, apmDomainRepresentation) +
				GenerateDataSourceFromRepresentationMap("oci_apm_synthetics_public_vantage_point", "test_public_vantage_point", Optional, Create, publicVantagePointSingularDataSourceRepresentation) +
				compartmentIdVariableStr + PublicVantagePointResourceConfig,
			Check: resource.ComposeAggregateTestCheckFunc(
				resource.TestCheckResourceAttrSet(singularDatasourceName, "apm_domain_id"),
				resource.TestCheckResourceAttr(singularDatasourceName, "name", "OraclePublic-us-ashburn-1"),
				resource.TestCheckResourceAttr(singularDatasourceName, "display_name", "US East (Ashburn)"),
			),
		},
	})
}
