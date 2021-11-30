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
	oci_database_management "github.com/oracle/oci-go-sdk/v52/databasemanagement"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	DbManagementPrivateEndpointRequiredOnlyResource = DbManagementPrivateEndpointResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_database_management_db_management_private_endpoint", "test_db_management_private_endpoint", Required, Create, dbManagementPrivateEndpointRepresentation)

	DbManagementPrivateEndpointResourceConfig = DbManagementPrivateEndpointResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_database_management_db_management_private_endpoint", "test_db_management_private_endpoint", Optional, Update, dbManagementPrivateEndpointRepresentation)

	dbManagementPrivateEndpointSingularDataSourceRepresentation = map[string]interface{}{
		"db_management_private_endpoint_id": Representation{RepType: Required, Create: `${oci_database_management_db_management_private_endpoint.test_db_management_private_endpoint.id}`},
	}

	dbManagementPrivateEndpointDataSourceRepresentation = map[string]interface{}{
		"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id}`},
		"is_cluster":     Representation{RepType: Optional, Create: `false`},
		"name":           Representation{RepType: Optional, Create: `name`, Update: `name2`},
		"state":          Representation{RepType: Optional, Create: `ACTIVE`},
		"vcn_id":         Representation{RepType: Optional, Create: `${oci_core_vcn.test_vcn.id}`},
		"filter":         RepresentationGroup{Required, dbManagementPrivateEndpointDataSourceFilterRepresentation}}
	dbManagementPrivateEndpointDataSourceFilterRepresentation = map[string]interface{}{
		"name":   Representation{RepType: Required, Create: `id`},
		"values": Representation{RepType: Required, Create: []string{`${oci_database_management_db_management_private_endpoint.test_db_management_private_endpoint.id}`}},
	}

	dbManagementPrivateEndpointRepresentation = map[string]interface{}{
		"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id}`},
		"name":           Representation{RepType: Required, Create: `name`, Update: `name2`},
		"subnet_id":      Representation{RepType: Required, Create: `${oci_core_subnet.test_subnet.id}`},
		"description":    Representation{RepType: Optional, Create: `description`, Update: `description2`},
		"is_cluster":     Representation{RepType: Optional, Create: `false`},
		"nsg_ids":        Representation{RepType: Optional, Create: []string{`${oci_core_network_security_group.test_network_security_group.id}`}, Update: []string{}},
	}

	DbManagementPrivateEndpointResourceDependencies = GenerateResourceFromRepresentationMap("oci_core_network_security_group", "test_network_security_group", Required, Create, networkSecurityGroupRepresentation) +
		GenerateResourceFromRepresentationMap("oci_core_subnet", "test_subnet", Required, Create, subnetRepresentation) +
		GenerateResourceFromRepresentationMap("oci_core_vcn", "test_vcn", Required, Create, vcnRepresentation)
)

// issue-routing-tag: database_management/default
func TestDatabaseManagementDbManagementPrivateEndpointResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestDatabaseManagementDbManagementPrivateEndpointResource_basic")
	defer httpreplay.SaveScenario()

	provider := testAccProvider
	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	compartmentIdU := getEnvSettingWithDefault("compartment_id_for_update", compartmentId)
	compartmentIdUVariableStr := fmt.Sprintf("variable \"compartment_id_for_update\" { default = \"%s\" }\n", compartmentIdU)

	resourceName := "oci_database_management_db_management_private_endpoint.test_db_management_private_endpoint"
	datasourceName := "data.oci_database_management_db_management_private_endpoints.test_db_management_private_endpoints"
	singularDatasourceName := "data.oci_database_management_db_management_private_endpoint.test_db_management_private_endpoint"

	var resId, resId2 string
	// Save TF content to Create resource with optional properties. This has to be exactly the same as the config part in the "Create with optionals" step in the test.
	SaveConfigContent(config+compartmentIdVariableStr+DbManagementPrivateEndpointResourceDependencies+
		GenerateResourceFromRepresentationMap("oci_database_management_db_management_private_endpoint", "test_db_management_private_endpoint", Optional, Create, dbManagementPrivateEndpointRepresentation), "databasemanagement", "dbManagementPrivateEndpoint", t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{
			"oci": provider,
		},
		CheckDestroy: testAccCheckDatabaseManagementDbManagementPrivateEndpointDestroy,
		Steps: []resource.TestStep{
			// verify Create
			{
				Config: config + compartmentIdVariableStr + DbManagementPrivateEndpointResourceDependencies +
					GenerateResourceFromRepresentationMap("oci_database_management_db_management_private_endpoint", "test_db_management_private_endpoint", Required, Create, dbManagementPrivateEndpointRepresentation),
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "name", "name"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
					resource.TestCheckResourceAttr(resourceName, "is_cluster", "false"),

					func(s *terraform.State) (err error) {
						resId, err = FromInstanceState(s, resourceName, "id")
						return err
					},
				),
			},

			// delete before next Create
			{
				Config: config + compartmentIdVariableStr + DbManagementPrivateEndpointResourceDependencies,
			},
			// verify Create with optionals
			{
				Config: config + compartmentIdVariableStr + DbManagementPrivateEndpointResourceDependencies +
					GenerateResourceFromRepresentationMap("oci_database_management_db_management_private_endpoint", "test_db_management_private_endpoint", Optional, Create, dbManagementPrivateEndpointRepresentation),
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "name"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),
					resource.TestCheckResourceAttr(resourceName, "nsg_ids.#", "1"),

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
				Config: config + compartmentIdVariableStr + compartmentIdUVariableStr + DbManagementPrivateEndpointResourceDependencies +
					GenerateResourceFromRepresentationMap("oci_database_management_db_management_private_endpoint", "test_db_management_private_endpoint", Optional, Create,
						RepresentationCopyWithNewProperties(dbManagementPrivateEndpointRepresentation, map[string]interface{}{
							"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id_for_update}`},
						})),
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentIdU),
					resource.TestCheckResourceAttr(resourceName, "description", "description"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "name"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),
					resource.TestCheckResourceAttr(resourceName, "is_cluster", "false"),

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
				Config: config + compartmentIdVariableStr + DbManagementPrivateEndpointResourceDependencies +
					GenerateResourceFromRepresentationMap("oci_database_management_db_management_private_endpoint", "test_db_management_private_endpoint", Optional, Update, dbManagementPrivateEndpointRepresentation),
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(resourceName, "description", "description2"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "name", "name2"),
					resource.TestCheckResourceAttrSet(resourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(resourceName, "subnet_id"),
					resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),
					resource.TestCheckResourceAttr(resourceName, "nsg_ids.#", "0"),
					resource.TestCheckResourceAttr(resourceName, "is_cluster", "false"),

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
					GenerateDataSourceFromRepresentationMap("oci_database_management_db_management_private_endpoints", "test_db_management_private_endpoints", Optional, Update, dbManagementPrivateEndpointDataSourceRepresentation) +
					compartmentIdVariableStr + DbManagementPrivateEndpointResourceDependencies +
					GenerateResourceFromRepresentationMap("oci_database_management_db_management_private_endpoint", "test_db_management_private_endpoint", Optional, Update, dbManagementPrivateEndpointRepresentation),
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(datasourceName, "name", "name2"),
					resource.TestCheckResourceAttr(datasourceName, "state", "ACTIVE"),
					resource.TestCheckResourceAttrSet(datasourceName, "vcn_id"),
					resource.TestCheckResourceAttr(datasourceName, "db_management_private_endpoint_collection.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "db_management_private_endpoint_collection.0.items.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "is_cluster", "false"),
				),
			},
			// verify singular datasource
			{
				Config: config +
					GenerateDataSourceFromRepresentationMap("oci_database_management_db_management_private_endpoint", "test_db_management_private_endpoint", Required, Create, dbManagementPrivateEndpointSingularDataSourceRepresentation) +
					compartmentIdVariableStr + DbManagementPrivateEndpointResourceConfig,
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttrSet(singularDatasourceName, "db_management_private_endpoint_id"),
					resource.TestCheckResourceAttr(singularDatasourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttr(singularDatasourceName, "description", "description2"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "id"),
					resource.TestCheckResourceAttr(singularDatasourceName, "name", "name2"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "private_ip"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "state"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "time_created"),
					resource.TestCheckResourceAttr(singularDatasourceName, "is_cluster", "false"),
				),
			},
			// remove singular datasource from previous step so that it doesn't conflict with import tests
			{
				Config: config + compartmentIdVariableStr + DbManagementPrivateEndpointResourceConfig,
			},
			// verify resource import
			{
				Config:                  config,
				ImportState:             true,
				ImportStateVerify:       true,
				ImportStateVerifyIgnore: []string{},
				ResourceName:            resourceName,
			},
		},
	})
}

func testAccCheckDatabaseManagementDbManagementPrivateEndpointDestroy(s *terraform.State) error {
	noResourceFound := true
	client := testAccProvider.Meta().(*OracleClients).dbManagementClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oci_database_management_db_management_private_endpoint" {
			noResourceFound = false
			request := oci_database_management.GetDbManagementPrivateEndpointRequest{}

			tmp := rs.Primary.ID
			request.DbManagementPrivateEndpointId = &tmp

			request.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "database_management")

			response, err := client.GetDbManagementPrivateEndpoint(context.Background(), request)

			if err == nil {
				deletedLifecycleStates := map[string]bool{
					string(oci_database_management.LifecycleStatesDeleted): true,
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
	if !InSweeperExcludeList("DatabaseManagementDbManagementPrivateEndpoint") {
		resource.AddTestSweepers("DatabaseManagementDbManagementPrivateEndpoint", &resource.Sweeper{
			Name:         "DatabaseManagementDbManagementPrivateEndpoint",
			Dependencies: DependencyGraph["dbManagementPrivateEndpoint"],
			F:            sweepDatabaseManagementDbManagementPrivateEndpointResource,
		})
	}
}

func sweepDatabaseManagementDbManagementPrivateEndpointResource(compartment string) error {
	dbManagementClient := GetTestClients(&schema.ResourceData{}).dbManagementClient()
	dbManagementPrivateEndpointIds, err := getDbManagementPrivateEndpointIds(compartment)
	if err != nil {
		return err
	}
	for _, dbManagementPrivateEndpointId := range dbManagementPrivateEndpointIds {
		if ok := SweeperDefaultResourceId[dbManagementPrivateEndpointId]; !ok {
			deleteDbManagementPrivateEndpointRequest := oci_database_management.DeleteDbManagementPrivateEndpointRequest{}

			deleteDbManagementPrivateEndpointRequest.DbManagementPrivateEndpointId = &dbManagementPrivateEndpointId

			deleteDbManagementPrivateEndpointRequest.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "database_management")
			_, error := dbManagementClient.DeleteDbManagementPrivateEndpoint(context.Background(), deleteDbManagementPrivateEndpointRequest)
			if error != nil {
				fmt.Printf("Error deleting DbManagementPrivateEndpoint %s %s, It is possible that the resource is already deleted. Please verify manually \n", dbManagementPrivateEndpointId, error)
				continue
			}
			WaitTillCondition(testAccProvider, &dbManagementPrivateEndpointId, dbManagementPrivateEndpointSweepWaitCondition, time.Duration(3*time.Minute),
				dbManagementPrivateEndpointSweepResponseFetchOperation, "database_management", true)
		}
	}
	return nil
}

func getDbManagementPrivateEndpointIds(compartment string) ([]string, error) {
	ids := GetResourceIdsToSweep(compartment, "DbManagementPrivateEndpointId")
	if ids != nil {
		return ids, nil
	}
	var resourceIds []string
	compartmentId := compartment
	dbManagementClient := GetTestClients(&schema.ResourceData{}).dbManagementClient()

	listDbManagementPrivateEndpointsRequest := oci_database_management.ListDbManagementPrivateEndpointsRequest{}
	listDbManagementPrivateEndpointsRequest.CompartmentId = &compartmentId
	listDbManagementPrivateEndpointsRequest.LifecycleState = oci_database_management.ListDbManagementPrivateEndpointsLifecycleStateActive
	listDbManagementPrivateEndpointsResponse, err := dbManagementClient.ListDbManagementPrivateEndpoints(context.Background(), listDbManagementPrivateEndpointsRequest)

	if err != nil {
		return resourceIds, fmt.Errorf("Error getting DbManagementPrivateEndpoint list for compartment id : %s , %s \n", compartmentId, err)
	}
	for _, dbManagementPrivateEndpoint := range listDbManagementPrivateEndpointsResponse.Items {
		id := *dbManagementPrivateEndpoint.Id
		resourceIds = append(resourceIds, id)
		AddResourceIdToSweeperResourceIdMap(compartmentId, "DbManagementPrivateEndpointId", id)
	}
	return resourceIds, nil
}

func dbManagementPrivateEndpointSweepWaitCondition(response common.OCIOperationResponse) bool {
	// Only stop if the resource is available beyond 3 mins. As there could be an issue for the sweeper to delete the resource and manual intervention required.
	if dbManagementPrivateEndpointResponse, ok := response.Response.(oci_database_management.GetDbManagementPrivateEndpointResponse); ok {
		return dbManagementPrivateEndpointResponse.LifecycleState != oci_database_management.LifecycleStatesDeleted
	}
	return false
}

func dbManagementPrivateEndpointSweepResponseFetchOperation(client *OracleClients, resourceId *string, retryPolicy *common.RetryPolicy) error {
	_, err := client.dbManagementClient().GetDbManagementPrivateEndpoint(context.Background(), oci_database_management.GetDbManagementPrivateEndpointRequest{
		DbManagementPrivateEndpointId: resourceId,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: retryPolicy,
		},
	})
	return err
}
