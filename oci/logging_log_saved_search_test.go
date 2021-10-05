// Copyright (c) 2017, 2021, Oracle and/or its affiliates. All rights reserved.
// Licensed under the Mozilla Public License v2.0

package oci

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/oracle/oci-go-sdk/v49/common"
	oci_logging "github.com/oracle/oci-go-sdk/v49/logging"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	LogSavedSearchRequiredOnlyResource = LogSavedSearchResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_logging_log_saved_search", "test_log_saved_search", Required, Create, logSavedSearchRepresentation)

	LogSavedSearchResourceConfig = LogSavedSearchResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_logging_log_saved_search", "test_log_saved_search", Optional, Update, logSavedSearchRepresentation)

	logSavedSearchSingularDataSourceRepresentation = map[string]interface{}{
		"log_saved_search_id": Representation{RepType: Required, Create: `${oci_logging_log_saved_search.test_log_saved_search.id}`},
	}

	logSavedSearchDataSourceRepresentation = map[string]interface{}{
		"compartment_id":      Representation{RepType: Required, Create: `${var.compartment_id}`},
		"log_saved_search_id": Representation{RepType: Optional, Create: `${oci_logging_log_saved_search.test_log_saved_search.id}`},
		"name":                Representation{RepType: Optional, Create: `name`, Update: `name2`},
		"filter":              RepresentationGroup{Required, logSavedSearchDataSourceFilterRepresentation}}

	logSavedSearchDataSourceFilterRepresentation = map[string]interface{}{
		"name":   Representation{RepType: Required, Create: `id`},
		"values": Representation{RepType: Required, Create: []string{`${oci_logging_log_saved_search.test_log_saved_search.id}`}},
	}

	logSavedSearchRepresentation = map[string]interface{}{
		"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id}`},
		"name":           Representation{RepType: Required, Create: `name`, Update: `name2`},
		"query":          Representation{RepType: Required, Create: `query`, Update: `query2`},
		"defined_tags":   Representation{RepType: Optional, Create: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "value")}`, Update: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "updatedValue")}`},
		"description":    Representation{RepType: Optional, Create: `description`, Update: `description2`},
		"freeform_tags":  Representation{RepType: Optional, Create: map[string]string{"Department": "Finance"}, Update: map[string]string{"Department": "Accounting"}},
	}

	LogSavedSearchResourceDependencies = DefinedTagsDependencies
)

// issue-routing-tag: logging/default
func TestLoggingLogSavedSearchResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestLoggingLogSavedSearchResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	compartmentIdU := getEnvSettingWithDefault("compartment_id_for_update", compartmentId)
	compartmentIdUVariableStr := fmt.Sprintf("variable \"compartment_id_for_update\" { default = \"%s\" }\n", compartmentIdU)

	resourceName := "oci_logging_log_saved_search.test_log_saved_search"
	datasourceName := "data.oci_logging_log_saved_searches.test_log_saved_searches"
	singularDatasourceName := "data.oci_logging_log_saved_search.test_log_saved_search"

	var resId, resId2 string
	// Save TF content to Create resource with optional properties. This has to be exactly the same as the config part in the "Create with optionals" step in the test.
	SaveConfigContent(config+compartmentIdVariableStr+LogSavedSearchResourceDependencies+
		GenerateResourceFromRepresentationMap("oci_logging_log_saved_search", "test_log_saved_search", Optional, Create, logSavedSearchRepresentation), "logging", "logSavedSearch", t)

	ResourceTest(t, testAccCheckLoggingLogSavedSearchDestroy, []resource.TestStep{
		// verify Create
		{
			Config: config + compartmentIdVariableStr + LogSavedSearchResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_logging_log_saved_search", "test_log_saved_search", Required, Create, logSavedSearchRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(resourceName, "name", "name"),
				resource.TestCheckResourceAttr(resourceName, "query", "query"),

				func(s *terraform.State) (err error) {
					resId, err = FromInstanceState(s, resourceName, "id")
					return err
				},
			),
		},

		// delete before next Create
		{
			Config: config + compartmentIdVariableStr + LogSavedSearchResourceDependencies,
		},
		// verify Create with optionals
		{
			Config: config + compartmentIdVariableStr + LogSavedSearchResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_logging_log_saved_search", "test_log_saved_search", Optional, Create, logSavedSearchRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "description", "description"),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttr(resourceName, "name", "name"),
				resource.TestCheckResourceAttr(resourceName, "query", "query"),

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
			Config: config + compartmentIdVariableStr + compartmentIdUVariableStr + LogSavedSearchResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_logging_log_saved_search", "test_log_saved_search", Optional, Create,
					RepresentationCopyWithNewProperties(logSavedSearchRepresentation, map[string]interface{}{
						"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id_for_update}`},
					})),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentIdU),
				resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "description", "description"),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttr(resourceName, "name", "name"),
				resource.TestCheckResourceAttr(resourceName, "query", "query"),

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
			Config: config + compartmentIdVariableStr + LogSavedSearchResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_logging_log_saved_search", "test_log_saved_search", Optional, Update, logSavedSearchRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "description", "description2"),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttr(resourceName, "name", "name2"),
				resource.TestCheckResourceAttr(resourceName, "query", "query2"),

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
				GenerateDataSourceFromRepresentationMap("oci_logging_log_saved_searches", "test_log_saved_searches", Optional, Update, logSavedSearchDataSourceRepresentation) +
				compartmentIdVariableStr + LogSavedSearchResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_logging_log_saved_search", "test_log_saved_search", Optional, Update, logSavedSearchRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttrSet(datasourceName, "log_saved_search_id"),
				resource.TestCheckResourceAttr(datasourceName, "name", "name2"),

				resource.TestCheckResourceAttr(datasourceName, "log_saved_search_summary_collection.#", "1"),
				resource.TestCheckResourceAttr(datasourceName, "log_saved_search_summary_collection.0.items.#", "1"),
			),
		},
		// verify singular datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_logging_log_saved_search", "test_log_saved_search", Required, Create, logSavedSearchSingularDataSourceRepresentation) +
				compartmentIdVariableStr + LogSavedSearchResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(singularDatasourceName, "log_saved_search_id"),

				resource.TestCheckResourceAttr(singularDatasourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(singularDatasourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(singularDatasourceName, "description", "description2"),
				resource.TestCheckResourceAttr(singularDatasourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "id"),
				resource.TestCheckResourceAttr(singularDatasourceName, "name", "name2"),
				resource.TestCheckResourceAttr(singularDatasourceName, "query", "query2"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "state"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "time_created"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "time_last_modified"),
			),
		},
		// remove singular datasource from previous step so that it doesn't conflict with import tests
		{
			Config: config + compartmentIdVariableStr + LogSavedSearchResourceConfig,
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

func testAccCheckLoggingLogSavedSearchDestroy(s *terraform.State) error {
	noResourceFound := true
	client := testAccProvider.Meta().(*OracleClients).loggingManagementClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oci_logging_log_saved_search" {
			noResourceFound = false
			request := oci_logging.GetLogSavedSearchRequest{}

			tmp := rs.Primary.ID
			request.LogSavedSearchId = &tmp

			request.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "logging")

			_, err := client.GetLogSavedSearch(context.Background(), request)

			if err == nil {
				return fmt.Errorf("resource still exists")
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
	if !InSweeperExcludeList("LoggingLogSavedSearch") {
		resource.AddTestSweepers("LoggingLogSavedSearch", &resource.Sweeper{
			Name:         "LoggingLogSavedSearch",
			Dependencies: DependencyGraph["logSavedSearch"],
			F:            sweepLoggingLogSavedSearchResource,
		})
	}
}

func sweepLoggingLogSavedSearchResource(compartment string) error {
	loggingManagementClient := GetTestClients(&schema.ResourceData{}).loggingManagementClient()
	logSavedSearchIds, err := getLogSavedSearchIds(compartment)
	if err != nil {
		return err
	}
	for _, logSavedSearchId := range logSavedSearchIds {
		if ok := SweeperDefaultResourceId[logSavedSearchId]; !ok {
			deleteLogSavedSearchRequest := oci_logging.DeleteLogSavedSearchRequest{}

			deleteLogSavedSearchRequest.LogSavedSearchId = &logSavedSearchId

			deleteLogSavedSearchRequest.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "logging")
			_, error := loggingManagementClient.DeleteLogSavedSearch(context.Background(), deleteLogSavedSearchRequest)
			if error != nil {
				fmt.Printf("Error deleting LogSavedSearch %s %s, It is possible that the resource is already deleted. Please verify manually \n", logSavedSearchId, error)
				continue
			}
		}
	}
	return nil
}

func getLogSavedSearchIds(compartment string) ([]string, error) {
	ids := GetResourceIdsToSweep(compartment, "LogSavedSearchId")
	if ids != nil {
		return ids, nil
	}
	var resourceIds []string
	compartmentId := compartment
	loggingManagementClient := GetTestClients(&schema.ResourceData{}).loggingManagementClient()

	listLogSavedSearchesRequest := oci_logging.ListLogSavedSearchesRequest{}
	listLogSavedSearchesRequest.CompartmentId = &compartmentId
	listLogSavedSearchesResponse, err := loggingManagementClient.ListLogSavedSearches(context.Background(), listLogSavedSearchesRequest)

	if err != nil {
		return resourceIds, fmt.Errorf("Error getting LogSavedSearch list for compartment id : %s , %s \n", compartmentId, err)
	}
	for _, logSavedSearch := range listLogSavedSearchesResponse.Items {
		id := *logSavedSearch.Id
		resourceIds = append(resourceIds, id)
		AddResourceIdToSweeperResourceIdMap(compartmentId, "LogSavedSearchId", id)
	}
	return resourceIds, nil
}
