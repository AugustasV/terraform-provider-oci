// Copyright (c) 2017, 2021, Oracle and/or its affiliates. All rights reserved.
// Licensed under the Mozilla Public License v2.0

package oci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"

	"github.com/oracle/oci-go-sdk/v51/core"
)

var (
	defaultRouteTable = `
resource "oci_core_default_route_table" "default" {
	manage_default_resource_id = "${oci_core_virtual_network.t.default_route_table_id}"
	route_rules {
		cidr_block = "0.0.0.0/0"
		network_entity_id = "${oci_core_internet_gateway.internet-gateway1.id}"
	}
}
`

	RouteTableScenarioTestDependencies = VcnResourceConfig + VcnResourceDependencies + ObjectStorageCoreService +
		GenerateResourceFromRepresentationMap("oci_core_local_peering_gateway", "test_local_peering_gateway", Required, Create, localPeeringGatewayRepresentation) +
		`
	resource "oci_core_internet_gateway" "test_internet_gateway" {
		compartment_id = "${var.compartment_id}"
		vcn_id = "${oci_core_vcn.test_vcn.id}"
		display_name = "-tf-internet-gateway"
	}

	resource "oci_core_service_gateway" "test_service_gateway" {
		#Required
		compartment_id = "${var.compartment_id}"
		services {
			service_id = "${lookup(data.oci_core_services.test_services.services[0], "id")}"
		}
		vcn_id = "${oci_core_vcn.test_vcn.id}"
	}`

	ObjectStorageCoreService = `data "oci_core_services" "test_services" {
  		filter {
    		name   = "name"
    		values = ["OCI .* Object Storage"]
			regex  = true
  		}
	}
	`

	routeTableRouteRulesRepresentationWithCidrBlock = map[string]interface{}{
		"network_entity_id": Representation{RepType: Required, Create: `${oci_core_internet_gateway.test_internet_gateway.id}`},
		"cidr_block":        Representation{RepType: Required, Create: `0.0.0.0/0`, Update: `10.0.0.0/8`},
	}
	routeTableRouteRulesRepresentationWithServiceCidr = map[string]interface{}{
		"network_entity_id": Representation{RepType: Required, Create: `${oci_core_service_gateway.test_service_gateway.id}`},
		"destination":       Representation{RepType: Required, Create: `${lookup(data.oci_core_services.test_services.services[0], "cidr_block")}`},
		"destination_type":  Representation{RepType: Required, Create: `SERVICE_CIDR_BLOCK`},
	}
	routeTableRouteRulesRepresentationWithServiceCidrAddingCidrBlock = map[string]interface{}{
		"network_entity_id": Representation{RepType: Required, Create: `${oci_core_service_gateway.test_service_gateway.id}`},
		"cidr_block":        Representation{RepType: Required, Create: `${lookup(data.oci_core_services.test_services.services[0], "cidr_block")}`},
		"destination":       Representation{RepType: Required, Create: `${lookup(data.oci_core_services.test_services.services[0], "cidr_block")}`},
		"destination_type":  Representation{RepType: Required, Create: `SERVICE_CIDR_BLOCK`},
	}
	routeTableRepresentationWithServiceCidr = GetUpdatedRepresentationCopy("route_rules", []RepresentationGroup{
		{Required, routeTableRouteRulesRepresentationWithServiceCidr},
		{Required, routeTableRouteRulesRepresentationWithCidrBlock}},
		routeTableRepresentationWithRouteRulesReqired,
	)
	routeTableRepresentationWithServiceCidrAddingCidrBlock = GetUpdatedRepresentationCopy("route_rules", []RepresentationGroup{
		{Required, routeTableRouteRulesRepresentationWithServiceCidrAddingCidrBlock},
		{Required, routeTableRouteRulesRepresentationWithCidrBlock}},
		routeTableRepresentationWithRouteRulesReqired,
	)
	routeTableRepresentationWithRouteRulesReqired = RepresentationCopyWithNewProperties(routeTableRepresentation, map[string]interface{}{
		"route_rules": RepresentationGroup{Required, routeTableRouteRulesRepresentationWithCidrBlock},
	})
)

// We needed to add a lot of special code to handle this case because of the terraform deficiency on differentiating values from statefile and from the config
// We test all the edge cases for that code here.
// issue-routing-tag: core/virtualNetwork
func TestResourceCoreRouteTable_deprecatedCidrBlock(t *testing.T) {
	httpreplay.SetScenario("TestResourceCoreRouteTable_deprecatedCidrBlock")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	resourceName := "oci_core_route_table.test_route_table"

	var resId, resId2 string

	ResourceTest(t, testAccCheckCoreRouteTableDestroy, []resource.TestStep{
		// verify Create
		{
			Config: config + compartmentIdVariableStr + RouteTableScenarioTestDependencies +
				GenerateResourceFromRepresentationMap("oci_core_route_table", "test_route_table", Required, Create, routeTableRepresentationWithRouteRulesReqired),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(resourceName, "route_rules.#", "1"),
				CheckResourceSetContainsElementWithProperties(resourceName, "route_rules", map[string]string{
					"cidr_block": "0.0.0.0/0",
				},
					[]string{
						"network_entity_id",
					}),
				resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),

				func(s *terraform.State) (err error) {
					resId, err = FromInstanceState(s, resourceName, "id")
					return err
				},
			),
		},
		// verify Update to deprecated cidr_block
		{
			Config: config + compartmentIdVariableStr + RouteTableScenarioTestDependencies +
				GenerateResourceFromRepresentationMap("oci_core_route_table", "test_route_table", Required, Update, routeTableRepresentationWithRouteRulesReqired),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttr(resourceName, "route_rules.#", "1"),
				CheckResourceSetContainsElementWithProperties(resourceName, "route_rules", map[string]string{"cidr_block": "10.0.0.0/8"}, []string{"network_entity_id"}),
				resource.TestCheckResourceAttrSet(resourceName, "state"),
				resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),

				func(s *terraform.State) (err error) {
					resId2, err = FromInstanceState(s, resourceName, "id")
					if resId != resId2 {
						return fmt.Errorf("Resource recreated when it was supposed to be updated.")
					}
					return err
				},
			),
		},
		// verify Update to network_id
		{
			Config: config + compartmentIdVariableStr + RouteTableScenarioTestDependencies +
				GenerateResourceFromRepresentationMap("oci_core_route_table", "test_route_table", Required, Update,
					GetUpdatedRepresentationCopy("route_rules.network_entity_id", Representation{RepType: Required, Create: `${oci_core_local_peering_gateway.test_local_peering_gateway.id}`},
						routeTableRepresentationWithRouteRulesReqired,
					)),

			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttr(resourceName, "route_rules.#", "1"),
				CheckResourceSetContainsElementWithProperties(resourceName, "route_rules", map[string]string{"cidr_block": "10.0.0.0/8"}, []string{"network_entity_id"}),
				resource.TestCheckResourceAttrSet(resourceName, "state"),
				resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),

				func(s *terraform.State) (err error) {
					resId2, err = FromInstanceState(s, resourceName, "id")
					if resId != resId2 {
						return fmt.Errorf("Resource recreated when it was supposed to be updated.")
					}
					return err
				},
			),
		},
		// verify Create with destination_type
		{
			Config: config + compartmentIdVariableStr + RouteTableScenarioTestDependencies +
				GenerateResourceFromRepresentationMap("oci_core_route_table", "test_route_table", Required, Create, routeTableRepresentationWithServiceCidr),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttr(resourceName, "route_rules.#", "2"),
				CheckResourceSetContainsElementWithProperties(resourceName, "route_rules", map[string]string{"destination_type": "SERVICE_CIDR_BLOCK"}, []string{"network_entity_id", "destination"}),
				CheckResourceSetContainsElementWithProperties(resourceName, "route_rules", map[string]string{"destination_type": "CIDR_BLOCK", "destination": "0.0.0.0/0"}, []string{"network_entity_id"}),
				resource.TestCheckResourceAttrSet(resourceName, "state"),
				resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),
			),
		},
		// verify Update after having a destination_type rule
		{
			Config: config + compartmentIdVariableStr + RouteTableScenarioTestDependencies +
				GenerateResourceFromRepresentationMap("oci_core_route_table", "test_route_table", Required, Update, routeTableRepresentationWithServiceCidr),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttr(resourceName, "route_rules.#", "2"),
				CheckResourceSetContainsElementWithProperties(resourceName, "route_rules", map[string]string{"destination_type": "SERVICE_CIDR_BLOCK"}, []string{"network_entity_id", "destination"}),
				CheckResourceSetContainsElementWithProperties(resourceName, "route_rules", map[string]string{"destination_type": "CIDR_BLOCK", "destination": "10.0.0.0/8"}, []string{"network_entity_id"}),
				resource.TestCheckResourceAttrSet(resourceName, "state"),
				resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),
			),
		},
		// verify adding cidr_block to a rule that has destination already
		{
			Config: config + compartmentIdVariableStr + RouteTableScenarioTestDependencies +
				GenerateResourceFromRepresentationMap("oci_core_route_table", "test_route_table", Required, Update, routeTableRepresentationWithServiceCidrAddingCidrBlock),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttr(resourceName, "route_rules.#", "2"),
				CheckResourceSetContainsElementWithProperties(resourceName, "route_rules", map[string]string{"destination_type": "SERVICE_CIDR_BLOCK"}, []string{"network_entity_id", "destination"}),
				CheckResourceSetContainsElementWithProperties(resourceName, "route_rules", map[string]string{"destination_type": "CIDR_BLOCK", "destination": "10.0.0.0/8"}, []string{"network_entity_id"}),
				resource.TestCheckResourceAttrSet(resourceName, "state"),
				resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),
			),
		},
		// We need to test that updating network entity also works when specifying destination instead of cidr_block
		// delete before next Create
		{
			Config: config + compartmentIdVariableStr + RouteTableScenarioTestDependencies,
		},
		//Create with optionals and destination
		{
			Config: config + compartmentIdVariableStr + RouteTableScenarioTestDependencies +
				GenerateResourceFromRepresentationMap("oci_core_route_table", "test_route_table", Optional, Update, routeTableRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "display_name", "displayName2"),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttr(resourceName, "route_rules.#", "1"),
				CheckResourceSetContainsElementWithProperties(resourceName, "route_rules", map[string]string{
					"destination":      "10.0.0.0/8",
					"destination_type": "CIDR_BLOCK",
				},
					[]string{
						"network_entity_id",
					}),
				resource.TestCheckResourceAttrSet(resourceName, "state"),
				resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),

				func(s *terraform.State) (err error) {
					resId, err = FromInstanceState(s, resourceName, "id")
					return err
				},
			),
		},
		// verify updates to network entity when using destination
		{
			Config: config + compartmentIdVariableStr + RouteTableScenarioTestDependencies +
				GenerateResourceFromRepresentationMap("oci_core_route_table", "test_route_table", Optional, Update,
					GetUpdatedRepresentationCopy("route_rules.network_entity_id", Representation{RepType: Required, Create: `${oci_core_local_peering_gateway.test_local_peering_gateway.id}`},
						routeTableRepresentation,
					)),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "display_name", "displayName2"),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttr(resourceName, "route_rules.#", "1"),
				CheckResourceSetContainsElementWithProperties(resourceName, "route_rules", map[string]string{"destination": "10.0.0.0/8", "destination_type": "CIDR_BLOCK"}, []string{"network_entity_id"}),
				resource.TestCheckResourceAttrSet(resourceName, "state"),
				resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),

				func(s *terraform.State) (err error) {
					resId2, err = FromInstanceState(s, resourceName, "id")
					if resId != resId2 {
						return fmt.Errorf("Resource recreated when it was supposed to be updated.")
					}
					return err
				},
			),
		},
	})
}

// issue-routing-tag: core/virtualNetwork
func TestResourceCoreRouteTable_defaultResource(t *testing.T) {
	httpreplay.SetScenario("TestResourceCoreRouteTable_defaultResource")
	defer httpreplay.SaveScenario()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	provider := testAccProvider
	config := testProviderConfig() + compartmentIdVariableStr + `
		resource "oci_core_virtual_network" "t" {
			compartment_id = "${var.compartment_id}"
			cidr_block = "10.0.0.0/16"
			display_name = "-tf-vcn"
		}

		resource "oci_core_internet_gateway" "internet-gateway1" {
			compartment_id = "${var.compartment_id}"
			vcn_id = "${oci_core_virtual_network.t.id}"
			display_name = "-tf-internet-gateway"
		}`

	compartmentIdU := getEnvSettingWithDefault("compartment_id_for_update", compartmentId)
	compartmentIdUVariableStr := fmt.Sprintf("variable \"compartment_id_for_update\" { default = \"%s\" }\n", compartmentIdU)
	resourceName := "oci_core_route_table.t"
	defaultResourceName := "oci_core_default_route_table.default"

	resource.Test(t, resource.TestCase{
		Providers: map[string]terraform.ResourceProvider{
			"oci": provider,
		},
		Steps: []resource.TestStep{
			// verify Create without rules
			{
				Config: config + `
					resource "oci_core_route_table" "t" {
						compartment_id = "${var.compartment_id}"
						vcn_id = "${oci_core_virtual_network.t.id}"
					}

					resource "oci_core_default_route_table" "default" {
						manage_default_resource_id = "${oci_core_virtual_network.t.default_route_table_id}"
					}`,
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "display_name"),
					resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),
					resource.TestCheckResourceAttrSet(resourceName, "compartment_id"),
					resource.TestCheckResourceAttr(resourceName, "state", string(core.RouteTableLifecycleStateAvailable)),
					resource.TestCheckResourceAttr(resourceName, "route_rules.#", "0"),
					resource.TestCheckResourceAttrSet(defaultResourceName, "manage_default_resource_id"),
					resource.TestCheckResourceAttrSet(defaultResourceName, "display_name"),
					resource.TestCheckResourceAttr(defaultResourceName, "route_rules.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "state", string(core.RouteTableLifecycleStateAvailable)),
				),
			},
			// verify add rule
			{
				Config: config + `
					resource "oci_core_route_table" "t" {
						compartment_id = "${var.compartment_id}"
						vcn_id = "${oci_core_virtual_network.t.id}"
						route_rules {
							cidr_block = "0.0.0.0/0"
							network_entity_id = "${oci_core_internet_gateway.internet-gateway1.id}"
						}
					}` + defaultRouteTable,
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttrSet(resourceName, "display_name"),
					resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),
					resource.TestCheckResourceAttrSet(resourceName, "compartment_id"),
					resource.TestCheckResourceAttr(resourceName, "state", string(core.RouteTableLifecycleStateAvailable)),
					resource.TestCheckResourceAttr(resourceName, "route_rules.#", "1"),
					CheckResourceSetContainsElementWithProperties(resourceName, "route_rules", map[string]string{
						"cidr_block": "0.0.0.0/0",
					},
						[]string{
							"network_entity_id",
						}),
					resource.TestCheckResourceAttrSet(defaultResourceName, "manage_default_resource_id"),
					resource.TestCheckResourceAttrSet(defaultResourceName, "compartment_id"),
					resource.TestCheckResourceAttr(defaultResourceName, "state", string(core.RouteTableLifecycleStateAvailable)),
					resource.TestCheckResourceAttrSet(defaultResourceName, "display_name"),
					resource.TestCheckResourceAttr(defaultResourceName, "route_rules.#", "1"),
					CheckResourceSetContainsElementWithProperties(defaultResourceName, "route_rules", map[string]string{
						"cidr_block": "0.0.0.0/0",
					},
						[]string{
							"network_entity_id",
						}),
				),
			},
			// verify Update
			{
				Config: compartmentIdUVariableStr + config + `
					resource "oci_core_route_table" "t" {
						compartment_id = "${var.compartment_id}"
						vcn_id = "${oci_core_virtual_network.t.id}"
						display_name = "-tf-route-table"
						route_rules {
							cidr_block = "0.0.0.0/0"
							network_entity_id = "${oci_core_internet_gateway.internet-gateway1.id}"
						}
						route_rules {
							cidr_block = "10.0.0.0/8"
							network_entity_id = "${oci_core_internet_gateway.internet-gateway1.id}"
						}
					}
					resource "oci_core_default_route_table" "default" {
						manage_default_resource_id = "${oci_core_virtual_network.t.default_route_table_id}"
						display_name = "default-tf-route-table"
						compartment_id = "${var.compartment_id_for_update}"
						route_rules {
							cidr_block = "0.0.0.0/0"
							network_entity_id = "${oci_core_internet_gateway.internet-gateway1.id}"
						}
						route_rules {
							cidr_block = "10.0.0.0/8"
							network_entity_id = "${oci_core_internet_gateway.internet-gateway1.id}"
						}
					}`,
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "-tf-route-table"),
					resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),
					resource.TestCheckResourceAttrSet(resourceName, "compartment_id"),
					resource.TestCheckResourceAttr(resourceName, "route_rules.#", "2"),
					CheckResourceSetContainsElementWithProperties(resourceName, "route_rules", map[string]string{
						"cidr_block": "0.0.0.0/0",
					},
						[]string{
							"network_entity_id",
						}),
					CheckResourceSetContainsElementWithProperties(resourceName, "route_rules", map[string]string{
						"cidr_block": "10.0.0.0/8",
					},
						[]string{
							"network_entity_id",
						}),
					resource.TestCheckResourceAttr(resourceName, "state", string(core.RouteTableLifecycleStateAvailable)),
					resource.TestCheckResourceAttrSet(defaultResourceName, "manage_default_resource_id"),
					resource.TestCheckResourceAttr(defaultResourceName, "compartment_id", compartmentIdU),
					resource.TestCheckResourceAttr(defaultResourceName, "display_name", "default-tf-route-table"),
					resource.TestCheckResourceAttr(defaultResourceName, "route_rules.#", "2"),
					CheckResourceSetContainsElementWithProperties(defaultResourceName, "route_rules", map[string]string{
						"cidr_block": "0.0.0.0/0",
					},
						[]string{
							"network_entity_id",
						}),
					CheckResourceSetContainsElementWithProperties(defaultResourceName, "route_rules", map[string]string{
						"cidr_block": "10.0.0.0/8",
					},
						[]string{
							"network_entity_id",
						}),
					resource.TestCheckResourceAttr(defaultResourceName, "state", string(core.RouteTableLifecycleStateAvailable)),
				),
			},
			// verify default resource delete
			{
				Config: config,
				Check:  nil,
			},
			// verify adding the default resource back to the config
			{
				Config: config + defaultRouteTable,
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttrSet(defaultResourceName, "manage_default_resource_id"),
					resource.TestCheckResourceAttrSet(defaultResourceName, "display_name"),
					resource.TestCheckResourceAttrSet(defaultResourceName, "compartment_id"),
					resource.TestCheckResourceAttr(defaultResourceName, "route_rules.#", "1"),
					CheckResourceSetContainsElementWithProperties(defaultResourceName, "route_rules", map[string]string{
						"cidr_block": "0.0.0.0/0",
					},
						[]string{
							"network_entity_id",
						}),
					resource.TestCheckResourceAttr(defaultResourceName, "state", string(core.RouteTableLifecycleStateAvailable)),
				),
			},
		},
	})
}
