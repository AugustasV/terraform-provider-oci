// Copyright (c) 2017, 2021, Oracle and/or its affiliates. All rights reserved.
// Licensed under the Mozilla Public License v2.0

package oci

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	bootVolumeAttachmentDataSourceRepresentation = map[string]interface{}{
		"availability_domain": Representation{RepType: Required, Create: `${data.oci_identity_availability_domains.test_availability_domains.availability_domains.0.name}`},
		"compartment_id":      Representation{RepType: Required, Create: `${var.compartment_id}`},
		"boot_volume_id":      Representation{RepType: Optional, Create: `${oci_core_instance.test_instance.boot_volume_id}`},
		"instance_id":         Representation{RepType: Optional, Create: `${oci_core_instance.test_instance.id}`},
	}

	BootVolumeAttachmentResourceConfig = GenerateResourceFromRepresentationMap("oci_core_boot_volume", "test_boot_volume", Required, Create, bootVolumeRepresentation) +
		GenerateResourceFromRepresentationMap("oci_core_subnet", "test_subnet", Required, Create, subnetRepresentation) +
		GenerateResourceFromRepresentationMap("oci_core_vcn", "test_vcn", Required, Create, vcnRepresentation) +
		OciImageIdsVariable +
		GenerateResourceFromRepresentationMap("oci_core_instance", "test_instance", Required, Create, instanceRepresentation) +
		AvailabilityDomainConfig
)

// issue-routing-tag: core/computeSharedOwnershipVmAndBm
func TestCoreBootVolumeAttachmentResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestCoreBootVolumeAttachmentResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	datasourceName := "data.oci_core_boot_volume_attachments.test_boot_volume_attachments"

	SaveConfigContent("", "", "", t)

	ResourceTest(t, nil, []resource.TestStep{
		// verify datasource can retrieve a specific attachment using server-side filtering
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_core_boot_volume_attachments", "test_boot_volume_attachments", Optional, Create, bootVolumeAttachmentDataSourceRepresentation) +
				compartmentIdVariableStr + BootVolumeAttachmentResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(datasourceName, "boot_volume_id"),
				resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttrSet(datasourceName, "instance_id"),

				resource.TestCheckResourceAttr(datasourceName, "boot_volume_attachments.#", "1"),
				resource.TestCheckResourceAttrSet(datasourceName, "boot_volume_attachments.0.availability_domain"),
				TestCheckResourceAttributesEqual(datasourceName, "boot_volume_attachments.0.boot_volume_id", "oci_core_instance.test_instance", "boot_volume_id"),
				resource.TestCheckResourceAttrSet(datasourceName, "boot_volume_attachments.0.compartment_id"),
				resource.TestCheckResourceAttrSet(datasourceName, "boot_volume_attachments.0.display_name"),
				resource.TestCheckResourceAttr(datasourceName, "boot_volume_attachments.0.encryption_in_transit_type", "NONE"),
				resource.TestCheckResourceAttrSet(datasourceName, "boot_volume_attachments.0.id"),
				TestCheckResourceAttributesEqual(datasourceName, "boot_volume_attachments.0.instance_id", "oci_core_instance.test_instance", "id"),
				resource.TestCheckResourceAttrSet(datasourceName, "boot_volume_attachments.0.state"),
				resource.TestCheckResourceAttrSet(datasourceName, "boot_volume_attachments.0.time_created"),
			),
		},
		// verify datasource can retrieve all boot volume attachments in a compartment by specifying no filtering options
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_core_boot_volume_attachments", "test_boot_volume_attachments", Required, Update, bootVolumeAttachmentDataSourceRepresentation) +
				compartmentIdVariableStr + BootVolumeAttachmentResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),
				resource.TestMatchResourceAttr(datasourceName, "boot_volume_attachments.#", regexp.MustCompile("[1-9][0-9]*")),
			),
		},
	})
}
