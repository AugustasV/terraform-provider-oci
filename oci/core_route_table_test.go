// Copyright (c) 2017, 2021, Oracle and/or its affiliates. All rights reserved.
// Licensed under the Mozilla Public License v2.0

package oci

import (
	"context"
	"fmt"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/oracle/oci-go-sdk/v52/common"
	oci_core "github.com/oracle/oci-go-sdk/v52/core"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	RouteTableRequiredOnlyResource = RouteTableResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_core_route_table", "test_route_table", Required, Create, routeTableRepresentation)

	RouteTableResource = RouteTableResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_core_route_table", "test_route_table", Optional, Create, routeTableRepresentation)

	routeTableDataSourceRepresentation = map[string]interface{}{
		"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id}`},
		"display_name":   Representation{RepType: Optional, Create: `MyRouteTable`, Update: `displayName2`},
		"state":          Representation{RepType: Optional, Create: `AVAILABLE`},
		"vcn_id":         Representation{RepType: Optional, Create: `${oci_core_vcn.test_vcn.id}`},
		"filter":         RepresentationGroup{Required, routeTableDataSourceFilterRepresentation}}
	routeTableDataSourceFilterRepresentation = map[string]interface{}{
		"name":   Representation{RepType: Required, Create: `id`},
		"values": Representation{RepType: Required, Create: []string{`${oci_core_route_table.test_route_table.id}`}},
	}

	routeTableRepresentation = map[string]interface{}{
		"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id}`},
		"vcn_id":         Representation{RepType: Required, Create: `${oci_core_vcn.test_vcn.id}`},
		"defined_tags":   Representation{RepType: Optional, Create: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "value")}`, Update: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "updatedValue")}`},
		"display_name":   Representation{RepType: Optional, Create: `MyRouteTable`, Update: `displayName2`},
		"freeform_tags":  Representation{RepType: Optional, Create: map[string]string{"Department": "Finance"}, Update: map[string]string{"Department": "Accounting"}},
		"route_rules":    RepresentationGroup{Optional, routeTableRouteRulesRepresentation},
	}
	routeTableRouteRulesRepresentation = map[string]interface{}{
		"network_entity_id": Representation{RepType: Required, Create: `${oci_core_internet_gateway.test_internet_gateway.id}`},
		"description":       Representation{RepType: Optional, Create: `description`, Update: `description2`},
		"destination":       Representation{RepType: Optional, Create: `0.0.0.0/0`, Update: `10.0.0.0/8`},
		"destination_type":  Representation{RepType: Optional, Create: `CIDR_BLOCK`},
	}

	RouteTableResourceDependencies = GenerateResourceFromRepresentationMap("oci_core_internet_gateway", "test_internet_gateway", Required, Create, internetGatewayRepresentation) +
		GenerateResourceFromRepresentationMap("oci_core_vcn", "test_vcn", Required, Create, vcnRepresentation) +
		DefinedTagsDependencies
)

// issue-routing-tag: core/virtualNetwork
func TestCoreRouteTableResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestCoreRouteTableResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	compartmentIdU := getEnvSettingWithDefault("compartment_id_for_update", compartmentId)
	compartmentIdUVariableStr := fmt.Sprintf("variable \"compartment_id_for_update\" { default = \"%s\" }\n", compartmentIdU)

	resourceName := "oci_core_route_table.test_route_table"
	datasourceName := "data.oci_core_route_tables.test_route_tables"

	var resId, resId2 string
	// Save TF content to Create resource with optional properties. This has to be exactly the same as the config part in the "Create with optionals" step in the test.
	SaveConfigContent(config+compartmentIdVariableStr+RouteTableResourceDependencies+
		GenerateResourceFromRepresentationMap("oci_core_route_table", "test_route_table", Optional, Create, routeTableRepresentation), "core", "routeTable", t)

	ResourceTest(t, testAccCheckCoreRouteTableDestroy, []resource.TestStep{
		// verify Create
		{
			Config: config + compartmentIdVariableStr + RouteTableResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_core_route_table", "test_route_table", Required, Create, routeTableRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),

				func(s *terraform.State) (err error) {
					resId, err = FromInstanceState(s, resourceName, "id")
					return err
				},
			),
		},

		// delete before next Create
		{
			Config: config + compartmentIdVariableStr + RouteTableResourceDependencies,
		},
		// verify Create with optionals
		{
			Config: config + compartmentIdVariableStr + RouteTableResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_core_route_table", "test_route_table", Optional, Create, routeTableRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "display_name", "MyRouteTable"),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttr(resourceName, "route_rules.#", "1"),
				CheckResourceSetContainsElementWithProperties(resourceName, "route_rules", map[string]string{
					"description":      "description",
					"destination":      "0.0.0.0/0",
					"destination_type": "CIDR_BLOCK",
				},
					[]string{
						"network_entity_id",
					}),
				resource.TestCheckResourceAttrSet(resourceName, "state"),
				resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),

				func(s *terraform.State) (err error) {
					resId, err = FromInstanceState(s, resourceName, "id")
					if isEnableExportCompartment, _ := strconv.ParseBool(getEnvSettingWithDefault("enable_export_compartment", "true")); isEnableExportCompartment {
						if errExport := TestExportCompartmentWithResourceName(&resId, &compartmentId, resourceName); errExport != nil {
							return errExport
						}
					}
					return err
				},
			),
		},

		// verify Update to the compartment (the compartment will be switched back in the next step)
		{
			Config: config + compartmentIdVariableStr + compartmentIdUVariableStr + RouteTableResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_core_route_table", "test_route_table", Optional, Create,
					RepresentationCopyWithNewProperties(routeTableRepresentation, map[string]interface{}{
						"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id_for_update}`},
					})),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentIdU),
				resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "display_name", "MyRouteTable"),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttr(resourceName, "route_rules.#", "1"),
				CheckResourceSetContainsElementWithProperties(resourceName, "route_rules", map[string]string{
					"description":      "description",
					"destination":      "0.0.0.0/0",
					"destination_type": "CIDR_BLOCK",
				},
					[]string{
						"network_entity_id",
					}),
				resource.TestCheckResourceAttrSet(resourceName, "state"),
				resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),

				func(s *terraform.State) (err error) {
					resId2, err = FromInstanceState(s, resourceName, "id")
					if resId != resId2 {
						return fmt.Errorf("resource recreated when it was supposed to be updated")
					}
					return err
				},
			),
		},

		// verify updates to updatable parameters
		{
			Config: config + compartmentIdVariableStr + RouteTableResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_core_route_table", "test_route_table", Optional, Update, routeTableRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "display_name", "displayName2"),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttr(resourceName, "route_rules.#", "1"),
				CheckResourceSetContainsElementWithProperties(resourceName, "route_rules", map[string]string{
					"description":      "description2",
					"destination":      "10.0.0.0/8",
					"destination_type": "CIDR_BLOCK",
				},
					[]string{
						"network_entity_id",
					}),
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
		// verify datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_core_route_tables", "test_route_tables", Optional, Update, routeTableDataSourceRepresentation) +
				compartmentIdVariableStr + RouteTableResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_core_route_table", "test_route_table", Optional, Update, routeTableRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(datasourceName, "display_name", "displayName2"),
				resource.TestCheckResourceAttr(datasourceName, "state", "AVAILABLE"),
				resource.TestCheckResourceAttrSet(datasourceName, "vcn_id"),

				resource.TestCheckResourceAttr(datasourceName, "route_tables.#", "1"),
				resource.TestCheckResourceAttr(datasourceName, "route_tables.0.compartment_id", compartmentId),
				resource.TestCheckResourceAttr(datasourceName, "route_tables.0.defined_tags.%", "1"),
				resource.TestCheckResourceAttr(datasourceName, "route_tables.0.display_name", "displayName2"),
				resource.TestCheckResourceAttr(datasourceName, "route_tables.0.freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(datasourceName, "route_tables.0.id"),
				resource.TestCheckResourceAttr(datasourceName, "route_tables.0.route_rules.#", "1"),
				CheckResourceSetContainsElementWithProperties(datasourceName, "route_tables.0.route_rules", map[string]string{
					"description":      "description2",
					"destination":      "10.0.0.0/8",
					"destination_type": "CIDR_BLOCK",
				},
					[]string{
						"network_entity_id",
					}),
				resource.TestCheckResourceAttrSet(datasourceName, "route_tables.0.state"),
				resource.TestCheckResourceAttrSet(datasourceName, "route_tables.0.vcn_id"),
			),
		},
		// verify resource import
		{
			Config:                  config,
			ImportState:             true,
			ImportStateVerify:       true,
			ImportStateVerifyIgnore: []string{},
			ResourceName:            resourceName,
		},
	})
}

func testAccCheckCoreRouteTableDestroy(s *terraform.State) error {
	noResourceFound := true
	client := testAccProvider.Meta().(*OracleClients).virtualNetworkClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oci_core_route_table" {
			noResourceFound = false
			request := oci_core.GetRouteTableRequest{}

			tmp := rs.Primary.ID
			request.RtId = &tmp

			request.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "core")

			response, err := client.GetRouteTable(context.Background(), request)

			if err == nil {
				deletedLifecycleStates := map[string]bool{
					string(oci_core.RouteTableLifecycleStateTerminated): true,
				}
				if _, ok := deletedLifecycleStates[string(response.LifecycleState)]; !ok {
					//resource lifecycle state is not in expected deleted lifecycle states.
					return fmt.Errorf("resource lifecycle state: %s is not in expected deleted lifecycle states", response.LifecycleState)
				}
				//resource lifecycle state is in expected deleted lifecycle states. continue with next one.
				continue
			}

			//Verify that exception is for '404 not found'.
			if failure, isServiceError := common.IsServiceError(err); !isServiceError || failure.GetHTTPStatusCode() != 404 {
				return err
			}
		}
	}
	if noResourceFound {
		return fmt.Errorf("at least one resource was expected from the state file, but could not be found")
	}

	return nil
}

func init() {
	if DependencyGraph == nil {
		initDependencyGraph()
	}
	if !InSweeperExcludeList("CoreRouteTable") {
		resource.AddTestSweepers("CoreRouteTable", &resource.Sweeper{
			Name:         "CoreRouteTable",
			Dependencies: DependencyGraph["routeTable"],
			F:            sweepCoreRouteTableResource,
		})
	}
}

func sweepCoreRouteTableResource(compartment string) error {
	virtualNetworkClient := GetTestClients(&schema.ResourceData{}).virtualNetworkClient()
	routeTableIds, err := getRouteTableIds(compartment)
	if err != nil {
		return err
	}
	for _, routeTableId := range routeTableIds {
		if ok := SweeperDefaultResourceId[routeTableId]; !ok {
			deleteRouteTableRequest := oci_core.DeleteRouteTableRequest{}

			deleteRouteTableRequest.RtId = &routeTableId

			deleteRouteTableRequest.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "core")
			_, error := virtualNetworkClient.DeleteRouteTable(context.Background(), deleteRouteTableRequest)
			if error != nil {
				fmt.Printf("Error deleting RouteTable %s %s, It is possible that the resource is already deleted. Please verify manually \n", routeTableId, error)
				continue
			}
			WaitTillCondition(testAccProvider, &routeTableId, routeTableSweepWaitCondition, time.Duration(3*time.Minute),
				routeTableSweepResponseFetchOperation, "core", true)
		}
	}
	return nil
}

func getRouteTableIds(compartment string) ([]string, error) {
	ids := GetResourceIdsToSweep(compartment, "RouteTableId")
	if ids != nil {
		return ids, nil
	}
	var resourceIds []string
	compartmentId := compartment
	virtualNetworkClient := GetTestClients(&schema.ResourceData{}).virtualNetworkClient()

	listRouteTablesRequest := oci_core.ListRouteTablesRequest{}
	listRouteTablesRequest.CompartmentId = &compartmentId
	listRouteTablesRequest.LifecycleState = oci_core.RouteTableLifecycleStateAvailable
	listRouteTablesResponse, err := virtualNetworkClient.ListRouteTables(context.Background(), listRouteTablesRequest)

	if err != nil {
		return resourceIds, fmt.Errorf("Error getting RouteTable list for compartment id : %s , %s \n", compartmentId, err)
	}
	for _, routeTable := range listRouteTablesResponse.Items {
		id := *routeTable.Id
		resourceIds = append(resourceIds, id)
		AddResourceIdToSweeperResourceIdMap(compartmentId, "RouteTableId", id)
	}
	return resourceIds, nil
}

func routeTableSweepWaitCondition(response common.OCIOperationResponse) bool {
	// Only stop if the resource is available beyond 3 mins. As there could be an issue for the sweeper to delete the resource and manual intervention required.
	if routeTableResponse, ok := response.Response.(oci_core.GetRouteTableResponse); ok {
		return routeTableResponse.LifecycleState != oci_core.RouteTableLifecycleStateTerminated
	}
	return false
}

func routeTableSweepResponseFetchOperation(client *OracleClients, resourceId *string, retryPolicy *common.RetryPolicy) error {
	_, err := client.virtualNetworkClient().GetRouteTable(context.Background(), oci_core.GetRouteTableRequest{
		RtId: resourceId,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: retryPolicy,
		},
	})
	return err
}
