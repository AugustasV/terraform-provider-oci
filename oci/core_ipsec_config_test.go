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
	ipSecConnectionDeviceConfigSingularDataSourceRepresentation = map[string]interface{}{
		"ipsec_id": Representation{RepType: Required, Create: `${oci_core_ipsec.test_ip_sec_connection.id}`},
	}

	IpSecConnectionDeviceConfigResourceConfig = IpSecConnectionRequiredOnlyResource
)

// issue-routing-tag: core/default
func TestCoreIpSecConnectionDeviceConfigResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestCoreIpSecConnectionDeviceConfigResource_basic")
	defer httpreplay.SaveScenario()

	config := ProviderTestConfig()

	compartmentId := GetEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	singularDatasourceName := "data.oci_core_ipsec_config.test_ip_sec_connection_device_config"

	SaveConfigContent("", "", "", t)

	ResourceTest(t, nil, []resource.TestStep{
		// verify singular datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_core_ipsec_config", "test_ip_sec_connection_device_config", Required, Create, ipSecConnectionDeviceConfigSingularDataSourceRepresentation) +
				compartmentIdVariableStr + IpSecConnectionDeviceConfigResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(singularDatasourceName, "ipsec_id"),

				resource.TestCheckResourceAttrSet(singularDatasourceName, "compartment_id"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "id"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "time_created"),
				resource.TestCheckResourceAttr(singularDatasourceName, "tunnels.#", "2"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "tunnels.0.ip_address"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "tunnels.0.shared_secret"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "tunnels.0.time_created"),
			),
		},
	})
}
