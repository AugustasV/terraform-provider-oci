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
	imageShapeSingularDataSourceRepresentation = map[string]interface{}{
		"image_id":   Representation{RepType: Required, Create: `${var.FlexInstanceImageOCID[var.region]}`},
		"shape_name": Representation{RepType: Required, Create: `VM.Standard.E3.Flex`},
	}

	imageShapeDataSourceRepresentation = map[string]interface{}{
		"image_id": Representation{RepType: Required, Create: `${var.FlexInstanceImageOCID[var.region]}`},
	}

	ImageShapeResourceConfig = FlexVmImageIdsVariable
)

// issue-routing-tag: core/computeImaging
func TestCoreImageShapeResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestCoreImageShapeResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	datasourceName := "data.oci_core_image_shapes.test_image_shapes"
	singularDatasourceName := "data.oci_core_image_shape.test_image_shape"

	SaveConfigContent("", "", "", t)

	ResourceTest(t, nil, []resource.TestStep{
		// verify datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_core_image_shapes", "test_image_shapes", Required, Create, imageShapeDataSourceRepresentation) +
				compartmentIdVariableStr + ImageShapeResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(datasourceName, "image_id"),

				resource.TestCheckResourceAttrSet(datasourceName, "image_shape_compatibilities.#"),
				resource.TestCheckResourceAttrSet(datasourceName, "image_shape_compatibilities.0.image_id"),
				resource.TestCheckResourceAttrSet(datasourceName, "image_shape_compatibilities.0.shape"),
			),
		},
		// verify singular datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_core_image_shape", "test_image_shape", Required, Create, imageShapeSingularDataSourceRepresentation) +
				compartmentIdVariableStr + ImageShapeResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(singularDatasourceName, "image_id"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "shape_name"),

				resource.TestCheckResourceAttrSet(singularDatasourceName, "shape"),
			),
		},
	})
}
