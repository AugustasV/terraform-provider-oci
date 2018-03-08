// Copyright (c) 2017, Oracle and/or its affiliates. All rights reserved.

package provider

import (
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/oracle/oci-go-sdk/core"
	"github.com/stretchr/testify/suite"
)

type DatasourceCoreRouteTableTestSuite struct {
	suite.Suite
	Config       string
	Providers    map[string]terraform.ResourceProvider
	ResourceName string
}

func (s *DatasourceCoreRouteTableTestSuite) SetupTest() {
	s.Providers = testAccProviders
	s.Config = legacyTestProviderConfig() + `
	resource "oci_core_virtual_network" "t" {
		compartment_id = "${var.compartment_id}"
		display_name = "-tf-vcn"
		cidr_block = "10.0.0.0/16"
	}`

	s.ResourceName = "data.oci_core_route_tables.t"
}

func (s *DatasourceCoreRouteTableTestSuite) TestAccDatasourceRouteTable_basic() {
	resource.Test(s.T(), resource.TestCase{
		PreventPostDestroyRefresh: true,
		Providers:                 s.Providers,
		Steps: []resource.TestStep{
			{
				ImportState:       true,
				ImportStateVerify: true,
				Config: s.Config + `
				data "oci_core_route_tables" "t" {
					compartment_id = "${var.compartment_id}"
					vcn_id = "${oci_core_virtual_network.t.id}"
					filter {
						name = "display_name"
						values = ["Default Route Table.*"]
						regex = true
					}
				}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet(s.ResourceName, "vcn_id"),
					resource.TestCheckResourceAttr(s.ResourceName, "route_tables.#", "1"),
					resource.TestCheckResourceAttr(s.ResourceName, "route_tables.0.display_name", "Default Route Table for -tf-vcn"),
					resource.TestCheckResourceAttr(s.ResourceName, "route_tables.0.state", string(core.RouteTableLifecycleStateAvailable)),
					resource.TestCheckResourceAttrSet(s.ResourceName, "route_tables.0.id"),
					resource.TestCheckResourceAttrSet(s.ResourceName, "route_tables.0.compartment_id"),
					resource.TestCheckResourceAttrSet(s.ResourceName, "route_tables.0.vcn_id"),
					resource.TestCheckResourceAttrSet(s.ResourceName, "route_tables.0.time_created"),
					resource.TestCheckResourceAttr(s.ResourceName, "route_tables.0.route_rules.#", "0"),
				),
			},
		},
	},
	)
}

func TestDatasourceCoreRouteTableTestSuite(t *testing.T) {
	suite.Run(t, new(DatasourceCoreRouteTableTestSuite))
}
