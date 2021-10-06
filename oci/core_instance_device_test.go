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
	instanceDeviceDataSourceRepresentation = map[string]interface{}{
		"instance_id":  Representation{RepType: Required, Create: `${oci_core_instance.test_instance.id}`},
		"is_available": Representation{RepType: Optional, Create: `true`},
		"name":         Representation{RepType: Optional, Create: `/dev/oracleoci/oraclevdb`},
	}

	InstanceDeviceResourceConfig = GenerateResourceFromRepresentationMap("oci_core_subnet", "test_subnet", Required, Create, subnetRepresentation) +
		GenerateResourceFromRepresentationMap("oci_core_vcn", "test_vcn", Required, Create, vcnRepresentation) +
		OciImageIdsVariable +
		GenerateResourceFromRepresentationMap("oci_core_instance", "test_instance", Required, Create, instanceRepresentation) +
		AvailabilityDomainConfig
)

// issue-routing-tag: core/computeSharedOwnershipVmAndBm
func TestCoreInstanceDeviceResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestCoreInstanceDeviceResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	datasourceName := "data.oci_core_instance_devices.test_instance_devices"

	SaveConfigContent("", "", "", t)

	ResourceTest(t, nil, []resource.TestStep{
		// verify datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_core_instance_devices", "test_instance_devices", Optional, Create, instanceDeviceDataSourceRepresentation) +
				compartmentIdVariableStr + InstanceDeviceResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(datasourceName, "instance_id"),
				resource.TestCheckResourceAttr(datasourceName, "is_available", "true"),
				resource.TestCheckResourceAttr(datasourceName, "name", "/dev/oracleoci/oraclevdb"),

				resource.TestCheckResourceAttrSet(datasourceName, "devices.#"),
				resource.TestCheckResourceAttrSet(datasourceName, "devices.0.is_available"),
				resource.TestCheckResourceAttrSet(datasourceName, "devices.0.name"),
			),
		},
	})
}
