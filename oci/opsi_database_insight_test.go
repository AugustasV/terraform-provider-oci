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
	"github.com/oracle/oci-go-sdk/v55/common"
	oci_opsi "github.com/oracle/oci-go-sdk/v55/opsi"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	DatabaseInsightRequiredOnlyResource = DatabaseInsightResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_opsi_database_insight", "test_database_insight", Required, Create, databaseInsightRepresentation)

	DatabaseInsightResourceConfig = DatabaseInsightResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_opsi_database_insight", "test_database_insight", Optional, Update, databaseInsightRepresentation)

	databaseInsightSingularDataSourceRepresentation = map[string]interface{}{
		"database_insight_id": Representation{RepType: Required, Create: `${oci_opsi_database_insight.test_database_insight.id}`},
	}

	databaseInsightDataSourceRepresentation = map[string]interface{}{
		"compartment_id":               Representation{RepType: Optional, Create: `${var.compartment_id}`},
		"compartment_id_in_subtree":    Representation{RepType: Optional, Create: `false`},
		"database_type":                Representation{RepType: Optional, Create: []string{`EXTERNAL-NONCDB`}},
		"enterprise_manager_bridge_id": Representation{RepType: Optional, Create: `${var.enterprise_manager_bridge_id}`},
		"fields":                       Representation{RepType: Optional, Create: []string{`databaseName`, `databaseType`, `compartmentId`, `databaseDisplayName`, `freeformTags`, `definedTags`, `systemTags`}},
		"id":                           Representation{RepType: Optional, Create: `${oci_opsi_database_insight.test_database_insight.id}`},
		"state":                        Representation{RepType: Optional, Create: []string{`ACTIVE`}},
		"status":                       Representation{RepType: Optional, Create: []string{`ENABLED`}, Update: []string{`DISABLED`}},
		"filter":                       RepresentationGroup{Required, databaseInsightDataSourceFilterRepresentation},
	}

	databaseInsightDataSourceFilterRepresentation = map[string]interface{}{
		"name":   Representation{RepType: Required, Create: `id`},
		"values": Representation{RepType: Required, Create: []string{`${oci_opsi_database_insight.test_database_insight.id}`}},
	}

	databaseInsightRepresentation = map[string]interface{}{
		"compartment_id":                       Representation{RepType: Required, Create: `${var.compartment_id}`},
		"enterprise_manager_bridge_id":         Representation{RepType: Required, Create: `${var.enterprise_manager_bridge_id}`},
		"enterprise_manager_entity_identifier": Representation{RepType: Required, Create: `${var.enterprise_manager_entity_id}`},
		"enterprise_manager_identifier":        Representation{RepType: Required, Create: `${var.enterprise_manager_id}`},
		"status":                               Representation{RepType: Optional, Create: `ENABLED`, Update: `DISABLED`},
		"entity_source":                        Representation{RepType: Required, Create: `EM_MANAGED_EXTERNAL_DATABASE`, Update: `EM_MANAGED_EXTERNAL_DATABASE`},
		"defined_tags":                         Representation{RepType: Optional, Create: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "value")}`, Update: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "updatedValue")}`},
		"freeform_tags":                        Representation{RepType: Optional, Create: map[string]string{"bar-key": "value"}, Update: map[string]string{"Department": "Accounting"}},
		"lifecycle":                            RepresentationGroup{Required, ignoreChangesdatabaseInsightRepresentation},
	}

	ignoreChangesdatabaseInsightRepresentation = map[string]interface{}{
		"ignore_changes": Representation{RepType: Required, Create: []string{`defined_tags`}},
	}

	DatabaseInsightResourceDependencies = DefinedTagsDependencies
)

// issue-routing-tag: opsi/controlPlane
func TestOpsiDatabaseInsightResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestOpsiDatabaseInsightResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	compartmentIdU := getEnvSettingWithDefault("compartment_id_for_update", compartmentId)
	compartmentIdUVariableStr := fmt.Sprintf("variable \"compartment_id_for_update\" { default = \"%s\" }\n", compartmentIdU)

	emBridgeId := getEnvSettingWithBlankDefault("enterprise_manager_bridge_ocid")
	emBridgeIdVariableStr := fmt.Sprintf("variable \"enterprise_manager_bridge_id\" { default = \"%s\" }\n", emBridgeId)

	enterpriseManagerId := getEnvSettingWithBlankDefault("enterprise_manager_id")
	enterpriseManagerIdVariableStr := fmt.Sprintf("variable \"enterprise_manager_id\" { default = \"%s\" }\n", enterpriseManagerId)

	enterpriseManagerEntityId := getEnvSettingWithBlankDefault("enterprise_manager_entity_id")
	enterpriseManagerEntityIdVariableStr := fmt.Sprintf("variable \"enterprise_manager_entity_id\" { default = \"%s\" }\n", enterpriseManagerEntityId)

	resourceName := "oci_opsi_database_insight.test_database_insight"
	datasourceName := "data.oci_opsi_database_insights.test_database_insights"
	singularDatasourceName := "data.oci_opsi_database_insight.test_database_insight"

	var resId, resId2 string
	// Save TF content to create resource with optional properties. This has to be exactly the same as the config part in the "create with optionals" step in the test.
	SaveConfigContent(config+compartmentIdVariableStr+emBridgeIdVariableStr+enterpriseManagerIdVariableStr+enterpriseManagerEntityIdVariableStr+DatabaseInsightResourceDependencies+
		GenerateResourceFromRepresentationMap("oci_opsi_database_insight", "test_database_insight", Optional, Create, databaseInsightRepresentation), "opsi", "databaseInsight", t)

	ResourceTest(t, testAccCheckOpsiDatabaseInsightDestroy, []resource.TestStep{
		// verify create with optional
		{
			Config: config + compartmentIdVariableStr + emBridgeIdVariableStr + enterpriseManagerIdVariableStr + enterpriseManagerEntityIdVariableStr + DatabaseInsightResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_opsi_database_insight", "test_database_insight", Optional, Create, databaseInsightRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				//resource.TestCheckResourceAttrSet(resourceName, "database_id"), // Won't be available for EM managed databases
				//resource.TestCheckResourceAttrSet(resourceName, "database_name"),
				//resource.TestCheckResourceAttrSet(resourceName, "database_resource_type"),
				resource.TestCheckResourceAttrSet(resourceName, "enterprise_manager_bridge_id"),
				resource.TestCheckResourceAttr(resourceName, "enterprise_manager_entity_identifier", enterpriseManagerEntityId),
				resource.TestCheckResourceAttrSet(resourceName, "enterprise_manager_entity_name"),
				resource.TestCheckResourceAttrSet(resourceName, "enterprise_manager_entity_type"),
				resource.TestCheckResourceAttr(resourceName, "enterprise_manager_identifier", enterpriseManagerId),
				resource.TestCheckResourceAttr(resourceName, "entity_source", "EM_MANAGED_EXTERNAL_DATABASE"),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttrSet(resourceName, "state"),
				resource.TestCheckResourceAttrSet(resourceName, "status"),
				resource.TestCheckResourceAttrSet(resourceName, "time_created"),

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
		// verify update to the compartment (the compartment will be switched back in the next step)
		{
			Config: config + compartmentIdVariableStr + emBridgeIdVariableStr + enterpriseManagerIdVariableStr + enterpriseManagerEntityIdVariableStr + compartmentIdUVariableStr + DatabaseInsightResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_opsi_database_insight", "test_database_insight", Optional, Create,
					RepresentationCopyWithNewProperties(databaseInsightRepresentation, map[string]interface{}{
						"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id_for_update}`},
					})),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentIdU),
				//resource.TestCheckResourceAttrSet(resourceName, "database_id"), // Won't be available for EM managed databases
				//resource.TestCheckResourceAttrSet(resourceName, "database_name"),
				//resource.TestCheckResourceAttrSet(resourceName, "database_resource_type"),
				resource.TestCheckResourceAttrSet(resourceName, "enterprise_manager_bridge_id"),
				resource.TestCheckResourceAttr(resourceName, "enterprise_manager_entity_identifier", enterpriseManagerEntityId),
				resource.TestCheckResourceAttrSet(resourceName, "enterprise_manager_entity_name"),
				resource.TestCheckResourceAttrSet(resourceName, "enterprise_manager_entity_type"),
				resource.TestCheckResourceAttr(resourceName, "enterprise_manager_identifier", enterpriseManagerId),
				resource.TestCheckResourceAttr(resourceName, "entity_source", "EM_MANAGED_EXTERNAL_DATABASE"),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttrSet(resourceName, "state"),
				resource.TestCheckResourceAttrSet(resourceName, "status"),
				resource.TestCheckResourceAttrSet(resourceName, "time_created"),

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
			Config: config + compartmentIdVariableStr + emBridgeIdVariableStr + enterpriseManagerIdVariableStr + enterpriseManagerEntityIdVariableStr + DatabaseInsightResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_opsi_database_insight", "test_database_insight", Optional, Update, databaseInsightRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				//resource.TestCheckResourceAttrSet(resourceName, "database_id"), // Won't be available for EM managed databases
				//resource.TestCheckResourceAttrSet(resourceName, "database_name"),
				//resource.TestCheckResourceAttrSet(resourceName, "database_resource_type"),
				resource.TestCheckResourceAttrSet(resourceName, "enterprise_manager_bridge_id"),
				resource.TestCheckResourceAttr(resourceName, "enterprise_manager_entity_identifier", enterpriseManagerEntityId),
				resource.TestCheckResourceAttrSet(resourceName, "enterprise_manager_entity_name"),
				resource.TestCheckResourceAttrSet(resourceName, "enterprise_manager_entity_type"),
				resource.TestCheckResourceAttr(resourceName, "enterprise_manager_identifier", enterpriseManagerId),
				resource.TestCheckResourceAttr(resourceName, "entity_source", "EM_MANAGED_EXTERNAL_DATABASE"),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttrSet(resourceName, "state"),
				resource.TestCheckResourceAttrSet(resourceName, "status"),
				resource.TestCheckResourceAttr(resourceName, "status", "DISABLED"),
				resource.TestCheckResourceAttrSet(resourceName, "time_created"),

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
				GenerateDataSourceFromRepresentationMap("oci_opsi_database_insights", "test_database_insights", Optional, Update, databaseInsightDataSourceRepresentation) +
				compartmentIdVariableStr + emBridgeIdVariableStr + enterpriseManagerIdVariableStr + enterpriseManagerEntityIdVariableStr + DatabaseInsightResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_opsi_database_insight", "test_database_insight", Optional, Update, databaseInsightRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(datasourceName, "compartment_id_in_subtree", "false"),
				//resource.TestCheckResourceAttr(datasourceName, "database_id.#", "1"), // Won't be available for EM managed databases
				resource.TestCheckResourceAttr(datasourceName, "database_type.#", "1"),
				resource.TestCheckResourceAttrSet(datasourceName, "enterprise_manager_bridge_id"),
				resource.TestCheckResourceAttr(datasourceName, "fields.#", "7"),
				//resource.TestCheckResourceAttr(datasourceName, "id.#", "1"), // id is no more list. It is a string
				resource.TestCheckResourceAttr(datasourceName, "state.#", "1"),
				resource.TestCheckResourceAttr(datasourceName, "status.#", "1"),

				resource.TestCheckResourceAttr(datasourceName, "database_insights_collection.#", "1"),
				resource.TestCheckResourceAttr(datasourceName, "database_insights_collection.0.items.#", "1"),
			),
		},
		// verify singular datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_opsi_database_insight", "test_database_insight", Required, Create, databaseInsightSingularDataSourceRepresentation) +
				compartmentIdVariableStr + emBridgeIdVariableStr + enterpriseManagerIdVariableStr + enterpriseManagerEntityIdVariableStr + DatabaseInsightResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(singularDatasourceName, "compartment_id", compartmentId),
				//resource.TestCheckResourceAttr(singularDatasourceName, "connection_credential_details.#", "1"), //Won't be available for EM managed databses
				//resource.TestCheckResourceAttr(singularDatasourceName, "connection_details.#", "1"),
				//resource.TestCheckResourceAttrSet(singularDatasourceName, "connector_id"),
				//resource.TestCheckResourceAttrSet(singularDatasourceName, "database_display_name"),
				//resource.TestCheckResourceAttrSet(singularDatasourceName, "database_name"),
				//resource.TestCheckResourceAttrSet(singularDatasourceName, "database_resource_type"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "database_type"),
				//resource.TestCheckResourceAttrSet(singularDatasourceName, "database_version"),
				//resource.TestCheckResourceAttrSet(singularDatasourceName, "db_additional_details"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "enterprise_manager_entity_display_name"),
				resource.TestCheckResourceAttr(resourceName, "enterprise_manager_entity_identifier", enterpriseManagerEntityId),
				resource.TestCheckResourceAttrSet(resourceName, "enterprise_manager_entity_name"),
				resource.TestCheckResourceAttrSet(resourceName, "enterprise_manager_entity_type"),
				resource.TestCheckResourceAttr(resourceName, "enterprise_manager_identifier", enterpriseManagerId),
				resource.TestCheckResourceAttr(singularDatasourceName, "entity_source", "EM_MANAGED_EXTERNAL_DATABASE"),
				resource.TestCheckResourceAttr(singularDatasourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "id"),
				//resource.TestCheckResourceAttrSet(singularDatasourceName, "management_agent_id"),
				//resource.TestCheckResourceAttrSet(singularDatasourceName, "processor_count"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "state"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "status"),
				resource.TestCheckResourceAttr(singularDatasourceName, "status", "DISABLED"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "time_created"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "time_updated"),
			),
		},
		// remove singular datasource from previous step so that it doesn't conflict with import tests
		{
			Config: config + compartmentIdVariableStr + emBridgeIdVariableStr + enterpriseManagerIdVariableStr + enterpriseManagerEntityIdVariableStr + DatabaseInsightResourceConfig,
		},
		// verify enable
		{
			Config: config + compartmentIdVariableStr + emBridgeIdVariableStr + enterpriseManagerIdVariableStr + enterpriseManagerEntityIdVariableStr + DatabaseInsightResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_opsi_database_insight", "test_database_insight", Optional, Update,
					RepresentationCopyWithNewProperties(databaseInsightRepresentation, map[string]interface{}{
						"status": Representation{RepType: Required, Update: `ENABLED`},
					})),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "status", "ENABLED"),

				func(s *terraform.State) (err error) {
					resId2, err = FromInstanceState(s, resourceName, "id")
					if resId != resId2 {
						return fmt.Errorf("resource recreated when it was supposed to be updated")
					}
					return err
				},
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

func testAccCheckOpsiDatabaseInsightDestroy(s *terraform.State) error {
	noResourceFound := true
	client := testAccProvider.Meta().(*OracleClients).operationsInsightsClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oci_opsi_database_insight" {
			noResourceFound = false
			request := oci_opsi.GetDatabaseInsightRequest{}

			tmp := rs.Primary.ID
			request.DatabaseInsightId = &tmp

			request.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "opsi")

			response, err := client.GetDatabaseInsight(context.Background(), request)

			if err == nil {
				deletedLifecycleStates := map[string]bool{
					string(oci_opsi.LifecycleStateDeleted): true,
				}
				if _, ok := deletedLifecycleStates[string(response.GetLifecycleState())]; !ok {
					//resource lifecycle state is not in expected deleted lifecycle states.
					return fmt.Errorf("resource lifecycle state: %s is not in expected deleted lifecycle states", response.GetLifecycleState())
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
	if !InSweeperExcludeList("OpsiDatabaseInsight") {
		resource.AddTestSweepers("OpsiDatabaseInsight", &resource.Sweeper{
			Name:         "OpsiDatabaseInsight",
			Dependencies: DependencyGraph["databaseInsight"],
			F:            sweepOpsiDatabaseInsightResource,
		})
	}
}

func sweepOpsiDatabaseInsightResource(compartment string) error {
	operationsInsightsClient := GetTestClients(&schema.ResourceData{}).operationsInsightsClient()
	databaseInsightIds, err := getDatabaseInsightIds(compartment)
	if err != nil {
		return err
	}
	for _, databaseInsightId := range databaseInsightIds {
		if ok := SweeperDefaultResourceId[databaseInsightId]; !ok {
			deleteDatabaseInsightRequest := oci_opsi.DeleteDatabaseInsightRequest{}

			deleteDatabaseInsightRequest.DatabaseInsightId = &databaseInsightId

			deleteDatabaseInsightRequest.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "opsi")
			_, error := operationsInsightsClient.DeleteDatabaseInsight(context.Background(), deleteDatabaseInsightRequest)
			if error != nil {
				fmt.Printf("Error deleting DatabaseInsight %s %s, It is possible that the resource is already deleted. Please verify manually \n", databaseInsightId, error)
				continue
			}
			WaitTillCondition(testAccProvider, &databaseInsightId, databaseInsightSweepWaitCondition, time.Duration(3*time.Minute),
				databaseInsightSweepResponseFetchOperation, "opsi", true)
		}
	}
	return nil
}

func getDatabaseInsightIds(compartment string) ([]string, error) {
	ids := GetResourceIdsToSweep(compartment, "DatabaseInsightId")
	if ids != nil {
		return ids, nil
	}
	var resourceIds []string
	compartmentId := compartment
	operationsInsightsClient := GetTestClients(&schema.ResourceData{}).operationsInsightsClient()

	listDatabaseInsightsRequest := oci_opsi.ListDatabaseInsightsRequest{}
	listDatabaseInsightsRequest.CompartmentId = &compartmentId
	listDatabaseInsightsRequest.LifecycleState = []oci_opsi.LifecycleStateEnum{oci_opsi.LifecycleStateActive}
	listDatabaseInsightsResponse, err := operationsInsightsClient.ListDatabaseInsights(context.Background(), listDatabaseInsightsRequest)

	if err != nil {
		return resourceIds, fmt.Errorf("Error getting DatabaseInsight list for compartment id : %s , %s \n", compartmentId, err)
	}
	for _, databaseInsight := range listDatabaseInsightsResponse.Items {
		id := *databaseInsight.GetId()
		resourceIds = append(resourceIds, id)
		AddResourceIdToSweeperResourceIdMap(compartmentId, "DatabaseInsightId", id)
	}
	return resourceIds, nil
}

func databaseInsightSweepWaitCondition(response common.OCIOperationResponse) bool {
	// Only stop if the resource is available beyond 3 mins. As there could be an issue for the sweeper to delete the resource and manual intervention required.
	if databaseInsightResponse, ok := response.Response.(oci_opsi.GetDatabaseInsightResponse); ok {
		return databaseInsightResponse.GetLifecycleState() != oci_opsi.LifecycleStateDeleted
	}
	return false
}

func databaseInsightSweepResponseFetchOperation(client *OracleClients, resourceId *string, retryPolicy *common.RetryPolicy) error {
	_, err := client.operationsInsightsClient().GetDatabaseInsight(context.Background(), oci_opsi.GetDatabaseInsightRequest{
		DatabaseInsightId: resourceId,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: retryPolicy,
		},
	})
	return err
}
