// Copyright (c) 2017, 2021, Oracle and/or its affiliates. All rights reserved.
// Licensed under the Mozilla Public License v2.0

package oci

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/oracle/oci-go-sdk/v55/common"
	oci_nosql "github.com/oracle/oci-go-sdk/v55/nosql"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	IndexRequiredOnlyResource = IndexResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_nosql_index", "test_index", Required, Create, indexRepresentation)

	IndexResourceConfig = IndexResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_nosql_index", "test_index", Optional, Update, indexRepresentation)

	indexSingularDataSourceRepresentation = map[string]interface{}{
		"index_name":       Representation{RepType: Required, Create: `${oci_nosql_index.test_index.id}`},
		"table_name_or_id": Representation{RepType: Required, Create: `${oci_nosql_table.test_table.id}`},
		"compartment_id":   Representation{RepType: Required, Create: `${var.compartment_id}`},
	}

	indexDataSourceRepresentation = map[string]interface{}{
		"table_name_or_id": Representation{RepType: Required, Create: `${oci_nosql_table.test_table.id}`},
		"compartment_id":   Representation{RepType: Optional, Create: `${var.compartment_id}`},
		"name":             Representation{RepType: Optional, Create: `test_index`},
		"state":            Representation{RepType: Optional, Create: `ACTIVE`},
		"filter":           RepresentationGroup{Required, indexDataSourceFilterRepresentation}}
	indexDataSourceFilterRepresentation = map[string]interface{}{
		"name":   Representation{RepType: Required, Create: `name`},
		"values": Representation{RepType: Required, Create: []string{`${oci_nosql_index.test_index.name}`}},
	}

	indexRepresentation = map[string]interface{}{
		"keys":             RepresentationGroup{Required, indexKeysRepresentation},
		"name":             Representation{RepType: Required, Create: `test_index`},
		"table_name_or_id": Representation{RepType: Required, Create: `${oci_nosql_table.test_table.id}`},
	}
	indexKeysRepresentation = map[string]interface{}{
		"column_name": Representation{RepType: Required, Create: `name`},
	}

	indexOptionalRepresentation = map[string]interface{}{
		"keys":             RepresentationGroup{Required, indexKeyWithJsonRepresentation},
		"name":             Representation{RepType: Required, Create: `test_index`},
		"table_name_or_id": Representation{RepType: Required, Create: `${oci_nosql_table.test_table.id}`},
		"compartment_id":   Representation{RepType: Optional, Create: `${var.compartment_id}`},
		"is_if_not_exists": Representation{RepType: Optional, Create: `false`},
	}
	indexKeyWithJsonRepresentation = map[string]interface{}{
		"column_name":     Representation{RepType: Required, Create: `info`},
		"json_field_type": Representation{RepType: Optional, Create: `STRING`},
		"json_path":       Representation{RepType: Optional, Create: `info`},
	}

	IndexResourceDependencies = GenerateResourceFromRepresentationMap("oci_nosql_table", "test_table", Required, Create, tableRepresentation)
)

// issue-routing-tag: nosql/default
func TestNosqlIndexResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestNosqlIndexResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	resourceName := "oci_nosql_index.test_index"

	datasourceName := "data.oci_nosql_indexes.test_indexes"
	singularDatasourceName := "data.oci_nosql_index.test_index"

	var compositeId string

	// Save TF content to Create resource with optional properties. This has to be exactly the same as the config part in the "Create with optionals" step in the test.
	SaveConfigContent(config+compartmentIdVariableStr+IndexResourceDependencies+
		GenerateResourceFromRepresentationMap("oci_nosql_index", "test_index", Optional, Create, indexRepresentation), "nosql", "index", t)

	ResourceTest(t, testAccCheckNosqlIndexDestroy, []resource.TestStep{
		// verify Create
		{
			Config: config + compartmentIdVariableStr + IndexResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_nosql_index", "test_index", Required, Create, indexRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "keys.#", "1"),
				resource.TestCheckResourceAttr(resourceName, "keys.0.column_name", "name"),
				resource.TestCheckResourceAttr(resourceName, "name", "test_index"),
				resource.TestCheckResourceAttrSet(resourceName, "table_name_or_id"),
			),
		},

		// delete before next Create
		{
			Config: config + compartmentIdVariableStr + IndexResourceDependencies,
		},
		// verify Create with optionals
		{
			Config: config + compartmentIdVariableStr + IndexResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_nosql_index", "test_index", Optional, Create, indexOptionalRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(resourceName, "is_if_not_exists", "false"),
				resource.TestCheckResourceAttr(resourceName, "keys.#", "1"),
				resource.TestCheckResourceAttr(resourceName, "keys.0.column_name", "info"),
				resource.TestCheckResourceAttr(resourceName, "keys.0.json_field_type", "STRING"),
				resource.TestCheckResourceAttr(resourceName, "keys.0.json_path", "info"),
				resource.TestCheckResourceAttr(resourceName, "name", "test_index"),
				resource.TestCheckResourceAttrSet(resourceName, "table_name_or_id"),

				func(s *terraform.State) (err error) {
					indexName, err := FromInstanceState(s, resourceName, "id")
					tableName, _ := FromInstanceState(s, resourceName, "table_name_or_id")
					compositeId = "tables/" + tableName + "/indexes/" + indexName
					log.Printf("[DEBUG] Composite ID to import: %s", compositeId)
					if isEnableExportCompartment, _ := strconv.ParseBool(getEnvSettingWithDefault("enable_export_compartment", "true")); isEnableExportCompartment {
						if errExport := TestExportCompartmentWithResourceName(&compositeId, &compartmentId, resourceName); errExport != nil {
							return errExport
						}
					}
					return err
				},
			),
		},

		// verify datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_nosql_indexes", "test_indexes", Optional, Update, indexDataSourceRepresentation) +
				compartmentIdVariableStr + IndexResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_nosql_index", "test_index", Optional, Update, indexRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(datasourceName, "name", "test_index"),
				resource.TestCheckResourceAttr(datasourceName, "state", "ACTIVE"),
				resource.TestCheckResourceAttrSet(datasourceName, "table_name_or_id"),

				resource.TestCheckResourceAttr(datasourceName, "index_collection.#", "1"),
			),
		},

		// verify singular datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_nosql_index", "test_index", Required, Create, indexSingularDataSourceRepresentation) +
				compartmentIdVariableStr + IndexResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(singularDatasourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "index_name"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "table_name_or_id"),

				resource.TestCheckResourceAttr(singularDatasourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(singularDatasourceName, "keys.#", "1"),
				resource.TestCheckResourceAttr(singularDatasourceName, "keys.0.column_name", "name"),
				resource.TestCheckResourceAttr(singularDatasourceName, "name", "test_index"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "state"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "table_id"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "table_name"),
			),
		},
	})
}

func testAccCheckNosqlIndexDestroy(s *terraform.State) error {
	noResourceFound := true
	client := testAccProvider.Meta().(*OracleClients).nosqlClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oci_nosql_index" {
			noResourceFound = false
			request := oci_nosql.GetIndexRequest{}

			if value, ok := rs.Primary.Attributes["compartment_id"]; ok {
				request.CompartmentId = &value
			}

			if value, ok := rs.Primary.Attributes["name"]; ok {
				request.IndexName = &value
			}

			if value, ok := rs.Primary.Attributes["table_name_or_id"]; ok {
				request.TableNameOrId = &value
			}

			request.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "nosql")

			response, err := client.GetIndex(context.Background(), request)

			if err == nil {
				deletedLifecycleStates := map[string]bool{
					string(oci_nosql.IndexLifecycleStateDeleted): true,
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
	if !InSweeperExcludeList("NosqlIndex") {
		resource.AddTestSweepers("NosqlIndex", &resource.Sweeper{
			Name:         "NosqlIndex",
			Dependencies: DependencyGraph["index"],
			F:            sweepNosqlIndexResource,
		})
	}
}

func sweepNosqlIndexResource(compartment string) error {
	nosqlClient := GetTestClients(&schema.ResourceData{}).nosqlClient()
	indexIds, err := getIndexIds(compartment)
	if err != nil {
		return err
	}
	for _, indexId := range indexIds {
		if ok := SweeperDefaultResourceId[indexId]; !ok {
			deleteIndexRequest := oci_nosql.DeleteIndexRequest{}

			deleteIndexRequest.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "nosql")
			_, error := nosqlClient.DeleteIndex(context.Background(), deleteIndexRequest)
			if error != nil {
				fmt.Printf("Error deleting Index %s %s, It is possible that the resource is already deleted. Please verify manually \n", indexId, error)
				continue
			}
			WaitTillCondition(testAccProvider, &indexId, indexSweepWaitCondition, time.Duration(3*time.Minute),
				indexSweepResponseFetchOperation, "nosql", true)
		}
	}
	return nil
}

func getIndexIds(compartment string) ([]string, error) {
	ids := GetResourceIdsToSweep(compartment, "IndexId")
	if ids != nil {
		return ids, nil
	}
	var resourceIds []string
	compartmentId := compartment
	nosqlClient := GetTestClients(&schema.ResourceData{}).nosqlClient()

	listIndexesRequest := oci_nosql.ListIndexesRequest{}
	listIndexesRequest.CompartmentId = &compartmentId

	tableNameOrIds, error := getTableIds(compartment)
	if error != nil {
		return resourceIds, fmt.Errorf("Error getting tableNameOrId required for Index resource requests \n")
	}
	for _, tableNameOrId := range tableNameOrIds {
		listIndexesRequest.TableNameOrId = &tableNameOrId

		listIndexesRequest.LifecycleState = oci_nosql.ListIndexesLifecycleStateActive
		listIndexesResponse, err := nosqlClient.ListIndexes(context.Background(), listIndexesRequest)

		if err != nil {
			return resourceIds, fmt.Errorf("Error getting Index list for compartment id : %s , %s \n", compartmentId, err)
		}
		for _, index := range listIndexesResponse.Items {
			id := *index.Name
			resourceIds = append(resourceIds, id)
			AddResourceIdToSweeperResourceIdMap(compartmentId, "IndexId", id)
		}

	}
	return resourceIds, nil
}

func indexSweepWaitCondition(response common.OCIOperationResponse) bool {
	// Only stop if the resource is available beyond 3 mins. As there could be an issue for the sweeper to delete the resource and manual intervention required.
	if indexResponse, ok := response.Response.(oci_nosql.GetIndexResponse); ok {
		return indexResponse.LifecycleState != oci_nosql.IndexLifecycleStateDeleted
	}
	return false
}

func indexSweepResponseFetchOperation(client *OracleClients, resourceId *string, retryPolicy *common.RetryPolicy) error {
	_, err := client.nosqlClient().GetIndex(context.Background(), oci_nosql.GetIndexRequest{RequestMetadata: common.RequestMetadata{
		RetryPolicy: retryPolicy,
	},
	})
	return err
}
