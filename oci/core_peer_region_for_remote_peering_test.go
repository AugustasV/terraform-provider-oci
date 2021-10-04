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
	peerRegionForRemotePeeringDataSourceRepresentation = map[string]interface{}{}

	PeerRegionForRemotePeeringResourceConfig = ""
)

// issue-routing-tag: core/default
func TestCorePeerRegionForRemotePeeringResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestCorePeerRegionForRemotePeeringResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	datasourceName := "data.oci_core_peer_region_for_remote_peerings.test_peer_region_for_remote_peerings"

	SaveConfigContent("", "", "", t)

	ResourceTest(t, nil, []resource.TestStep{
		// verify datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_core_peer_region_for_remote_peerings", "test_peer_region_for_remote_peerings", Required, Create, peerRegionForRemotePeeringDataSourceRepresentation) +
				compartmentIdVariableStr + PeerRegionForRemotePeeringResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(

				resource.TestCheckResourceAttrSet(datasourceName, "peer_region_for_remote_peerings.#"),
				resource.TestCheckResourceAttrSet(datasourceName, "peer_region_for_remote_peerings.0.name"),
			),
		},
	})
}
