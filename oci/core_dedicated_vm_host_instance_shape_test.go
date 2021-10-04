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
	dedicatedVmHostInstanceShapeDataSourceRepresentation = map[string]interface{}{
		"compartment_id":          Representation{RepType: Required, Create: `${var.compartment_id}`},
		"availability_domain":     Representation{RepType: Optional, Create: `${data.oci_identity_availability_domains.test_availability_domains.availability_domains.0.name}`},
		"dedicated_vm_host_shape": Representation{RepType: Optional, Create: `DVH.Standard2.52`},
	}

	DedicatedVmHostInstanceShapeResourceConfig = AvailabilityDomainConfig
)

// issue-routing-tag: core/default
func TestCoreDedicatedVmHostInstanceShapeResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestCoreDedicatedVmHostInstanceShapeResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	datasourceName := "data.oci_core_dedicated_vm_host_instance_shapes.test_dedicated_vm_host_instance_shapes"

	SaveConfigContent("", "", "", t)

	ResourceTest(t, nil, []resource.TestStep{
		// verify datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_core_dedicated_vm_host_instance_shapes", "test_dedicated_vm_host_instance_shapes", Required, Create, dedicatedVmHostInstanceShapeDataSourceRepresentation) +
				compartmentIdVariableStr + DedicatedVmHostInstanceShapeResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),

				resource.TestCheckResourceAttrSet(datasourceName, "dedicated_vm_host_instance_shapes.#"),
				resource.TestCheckResourceAttrSet(datasourceName, "dedicated_vm_host_instance_shapes.0.availability_domain"),
				resource.TestCheckResourceAttrSet(datasourceName, "dedicated_vm_host_instance_shapes.0.instance_shape_name"),
			),
		},
	})
}
