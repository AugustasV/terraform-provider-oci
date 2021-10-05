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
	"github.com/oracle/oci-go-sdk/v49/common"
	oci_dataflow "github.com/oracle/oci-go-sdk/v49/dataflow"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	DataFlowApplicationRequiredOnlyResource = dataFlowApplicationResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_dataflow_application", "test_application", Required, Create, dataFlowApplicationRepresentation)

	DataFlowApplicationResourceConfig = dataFlowApplicationResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_dataflow_application", "test_application", Optional, Update, dataFlowApplicationRepresentation)

	dataFlowApplicationSingularDataSourceRepresentation = map[string]interface{}{
		"application_id": Representation{RepType: Required, Create: `${oci_dataflow_application.test_application.id}`},
	}

	dataFlowApplicationDataSourceRepresentation = map[string]interface{}{
		"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id}`},
		"display_name":   Representation{RepType: Optional, Create: `test_wordcount_app`, Update: `displayName2`},
		"filter":         RepresentationGroup{Required, dataFlowApplicationDataSourceFilterRepresentation}}
	dataFlowApplicationDataSourceFilterRepresentation = map[string]interface{}{
		"name":   Representation{RepType: Required, Create: `id`},
		"values": Representation{RepType: Required, Create: []string{`${oci_dataflow_application.test_application.id}`}},
	}

	dataFlowApplicationRepresentation = map[string]interface{}{
		"compartment_id":       Representation{RepType: Required, Create: `${var.compartment_id}`},
		"display_name":         Representation{RepType: Required, Create: `test_wordcount_app`, Update: `displayName2`},
		"driver_shape":         Representation{RepType: Required, Create: `VM.Standard2.1`},
		"executor_shape":       Representation{RepType: Required, Create: `VM.Standard2.1`},
		"file_uri":             Representation{RepType: Required, Create: `${var.dataflow_file_uri}`, Update: `${var.dataflow_file_uri_updated}`},
		"language":             Representation{RepType: Required, Create: `PYTHON`, Update: `SCALA`},
		"num_executors":        Representation{RepType: Required, Create: `1`, Update: `2`},
		"spark_version":        Representation{RepType: Required, Create: `2.4`, Update: `2.4.4`},
		"archive_uri":          Representation{RepType: Optional, Create: `${var.dataflow_archive_uri}`},
		"arguments":            Representation{RepType: Optional, Create: []string{`arguments`}, Update: []string{`arguments2`}},
		"configuration":        Representation{RepType: Optional, Create: map[string]string{"spark.shuffle.io.maxRetries": "10"}, Update: map[string]string{"spark.shuffle.io.maxRetries": "11"}},
		"defined_tags":         Representation{RepType: Optional, Create: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "value")}`, Update: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "updatedValue")}`},
		"description":          Representation{RepType: Optional, Create: `description`, Update: `description2`},
		"freeform_tags":        Representation{RepType: Optional, Create: map[string]string{"Department": "Finance"}, Update: map[string]string{"Department": "Accounting"}},
		"logs_bucket_uri":      Representation{RepType: Optional, Create: `${var.dataflow_logs_bucket_uri}`},
		"metastore_id":         Representation{RepType: Optional, Create: `${var.metastore_id}`},
		"parameters":           RepresentationGroup{Optional, applicationParametersRepresentation},
		"private_endpoint_id":  Representation{RepType: Optional, Create: `${oci_dataflow_private_endpoint.test_private_endpoint.id}`},
		"warehouse_bucket_uri": Representation{RepType: Optional, Create: `${var.dataflow_warehouse_bucket_uri}`},
	}
	applicationParametersRepresentation = map[string]interface{}{
		"name":  Representation{RepType: Required, Create: `name`, Update: `name2`},
		"value": Representation{RepType: Required, Create: `value`, Update: `value2`},
	}

	dataFlowApplicationResourceDependencies = GenerateResourceFromRepresentationMap("oci_core_subnet", "test_subnet", Required, Create, subnetRepresentation) +
		GenerateResourceFromRepresentationMap("oci_core_vcn", "test_vcn", Required, Create, vcnRepresentation) +
		GenerateResourceFromRepresentationMap("oci_dataflow_private_endpoint", "test_private_endpoint", Required, Create, privateEndpointRepresentation) +
		DefinedTagsDependencies
)

// issue-routing-tag: dataflow/default
func TestDataflowApplicationResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestDataflowApplicationResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	compartmentIdU := getEnvSettingWithDefault("compartment_id_for_update", compartmentId)
	compartmentIdUVariableStr := fmt.Sprintf("variable \"compartment_id_for_update\" { default = \"%s\" }\n", compartmentIdU)
	fileUri := getEnvSettingWithBlankDefault("dataflow_file_uri")
	fileUriUpdated := getEnvSettingWithBlankDefault("dataflow_file_uri_updated")
	fileUriVariableStr := fmt.Sprintf("variable \"dataflow_file_uri\" { default = \"%s\" }\n", fileUri)
	fileUriVariableStrUpdated := fmt.Sprintf("variable \"dataflow_file_uri_updated\" { default = \"%s\" }\n", fileUriUpdated)
	archiveUri := getEnvSettingWithBlankDefault("dataflow_archive_uri")
	archiveUriVariableStr := fmt.Sprintf("variable \"dataflow_archive_uri\" { default = \"%s\" }\n", archiveUri)

	logsBucketUri := getEnvSettingWithBlankDefault("dataflow_logs_bucket_uri")
	logsBucketUriVariableStr := fmt.Sprintf("variable \"dataflow_logs_bucket_uri\" { default = \"%s\" }\n", logsBucketUri)
	warehouseBucketUri := getEnvSettingWithBlankDefault("dataflow_warehouse_bucket_uri")
	warehouseBucketUriVariableStr := fmt.Sprintf("variable \"dataflow_warehouse_bucket_uri\" { default = \"%s\" }\n", warehouseBucketUri)
	classNameUpdated := getEnvSettingWithBlankDefault("dataflow_class_name_updated")
	classNameStrUpdated := fmt.Sprintf("variable \"dataflow_class_name_updated\" { default = \"%s\" }\n", classNameUpdated)
	resourceName := "oci_dataflow_application.test_application"
	datasourceName := "data.oci_dataflow_applications.test_applications"
	singularDatasourceName := "data.oci_dataflow_application.test_application"

	metastoreId := getEnvSettingWithBlankDefault("metastore_id")
	metastoreIdVariableStr := fmt.Sprintf("variable \"metastore_id\" { default = \"%s\" }\n", metastoreId)

	var resId, resId2 string
	// Save TF content to Create resource with optional properties. This has to be exactly the same as the config part in the "Create with optionals" step in the test.
	SaveConfigContent(config+compartmentIdVariableStr+ApplicationResourceDependencies+
		GenerateResourceFromRepresentationMap("oci_dataflow_application", "test_application", Optional, Create, applicationRepresentation), "dataflow", "application", t)

	ResourceTest(t, testAccCheckDataflowApplicationDestroy, []resource.TestStep{
		// verify Create
		{
			Config: config + compartmentIdVariableStr + fileUriVariableStr + archiveUriVariableStr + logsBucketUriVariableStr + warehouseBucketUriVariableStr + metastoreIdVariableStr + dataFlowApplicationResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_dataflow_application", "test_application", Required, Create, dataFlowApplicationRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(resourceName, "display_name", "test_wordcount_app"),
				resource.TestCheckResourceAttr(resourceName, "driver_shape", "VM.Standard2.1"),
				resource.TestCheckResourceAttr(resourceName, "executor_shape", "VM.Standard2.1"),
				resource.TestCheckResourceAttr(resourceName, "file_uri", fileUri),
				resource.TestCheckResourceAttr(resourceName, "language", "PYTHON"),
				resource.TestCheckResourceAttr(resourceName, "num_executors", "1"),
				resource.TestCheckResourceAttr(resourceName, "spark_version", "2.4"),

				func(s *terraform.State) (err error) {
					resId, err = FromInstanceState(s, resourceName, "id")
					return err
				},
			),
		},

		// delete before next Create
		{
			Config: config + compartmentIdVariableStr + fileUriVariableStr + archiveUriVariableStr + logsBucketUriVariableStr + warehouseBucketUriVariableStr + metastoreIdVariableStr + dataFlowApplicationResourceDependencies,
		},
		// verify Create with optionals
		{
			Config: config + compartmentIdVariableStr + fileUriVariableStr + archiveUriVariableStr + logsBucketUriVariableStr + warehouseBucketUriVariableStr + metastoreIdVariableStr + dataFlowApplicationResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_dataflow_application", "test_application", Optional, Create, dataFlowApplicationRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(resourceName, "archive_uri"),
				resource.TestCheckResourceAttr(resourceName, "arguments.#", "1"),
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(resourceName, "configuration.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "description", "description"),
				resource.TestCheckResourceAttr(resourceName, "display_name", "test_wordcount_app"),
				resource.TestCheckResourceAttr(resourceName, "driver_shape", "VM.Standard2.1"),
				resource.TestCheckResourceAttr(resourceName, "executor_shape", "VM.Standard2.1"),
				resource.TestCheckResourceAttr(resourceName, "file_uri", fileUri),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttr(resourceName, "language", "PYTHON"),
				resource.TestCheckResourceAttr(resourceName, "logs_bucket_uri", logsBucketUri),
				resource.TestCheckResourceAttr(resourceName, "metastore_id", metastoreId),
				resource.TestCheckResourceAttr(resourceName, "num_executors", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "owner_principal_id"),
				resource.TestCheckResourceAttr(resourceName, "parameters.#", "1"),
				resource.TestCheckResourceAttr(resourceName, "parameters.0.name", "name"),
				resource.TestCheckResourceAttr(resourceName, "parameters.0.value", "value"),
				resource.TestCheckResourceAttrSet(resourceName, "private_endpoint_id"),
				resource.TestCheckResourceAttr(resourceName, "spark_version", "2.4"),
				resource.TestCheckResourceAttrSet(resourceName, "state"),
				resource.TestCheckResourceAttrSet(resourceName, "time_created"),
				resource.TestCheckResourceAttrSet(resourceName, "time_updated"),
				resource.TestCheckResourceAttr(resourceName, "warehouse_bucket_uri", warehouseBucketUri),

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
			Config: config + compartmentIdVariableStr + compartmentIdUVariableStr + dataFlowApplicationResourceDependencies + fileUriVariableStr + archiveUriVariableStr + warehouseBucketUriVariableStr + logsBucketUriVariableStr + metastoreIdVariableStr +
				GenerateResourceFromRepresentationMap("oci_dataflow_application", "test_application", Optional, Create,
					RepresentationCopyWithNewProperties(dataFlowApplicationRepresentation, map[string]interface{}{
						"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id_for_update}`},
					})),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(resourceName, "archive_uri"),
				resource.TestCheckResourceAttr(resourceName, "arguments.#", "1"),
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentIdU),
				resource.TestCheckResourceAttr(resourceName, "configuration.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "description", "description"),
				resource.TestCheckResourceAttr(resourceName, "display_name", "test_wordcount_app"),
				resource.TestCheckResourceAttr(resourceName, "driver_shape", "VM.Standard2.1"),
				resource.TestCheckResourceAttr(resourceName, "executor_shape", "VM.Standard2.1"),
				resource.TestCheckResourceAttrSet(resourceName, "file_uri"),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttr(resourceName, "language", "PYTHON"),
				resource.TestCheckResourceAttrSet(resourceName, "logs_bucket_uri"),
				resource.TestCheckResourceAttr(resourceName, "metastore_id", metastoreId),
				resource.TestCheckResourceAttr(resourceName, "num_executors", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "owner_principal_id"),
				resource.TestCheckResourceAttr(resourceName, "parameters.#", "1"),
				resource.TestCheckResourceAttr(resourceName, "parameters.0.name", "name"),
				resource.TestCheckResourceAttr(resourceName, "parameters.0.value", "value"),
				resource.TestCheckResourceAttrSet(resourceName, "private_endpoint_id"),
				resource.TestCheckResourceAttr(resourceName, "spark_version", "2.4"),
				resource.TestCheckResourceAttrSet(resourceName, "state"),
				resource.TestCheckResourceAttrSet(resourceName, "time_created"),
				resource.TestCheckResourceAttrSet(resourceName, "time_updated"),
				resource.TestCheckResourceAttrSet(resourceName, "warehouse_bucket_uri"),

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
			Config: config + compartmentIdVariableStr + fileUriVariableStr + classNameStrUpdated + fileUriVariableStrUpdated + archiveUriVariableStr + logsBucketUriVariableStr + warehouseBucketUriVariableStr + metastoreIdVariableStr + dataFlowApplicationResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_dataflow_application", "test_application", Optional, Update,
					RepresentationCopyWithNewProperties(dataFlowApplicationRepresentation, map[string]interface{}{
						"class_name": Representation{RepType: Optional, Create: `${var.dataflow_class_name_updated}`},
					})),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(resourceName, "archive_uri"),
				resource.TestCheckResourceAttr(resourceName, "arguments.#", "1"),
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(resourceName, "configuration.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "description", "description2"),
				resource.TestCheckResourceAttr(resourceName, "display_name", "displayName2"),
				resource.TestCheckResourceAttr(resourceName, "driver_shape", "VM.Standard2.1"),
				resource.TestCheckResourceAttr(resourceName, "executor_shape", "VM.Standard2.1"),
				resource.TestCheckResourceAttr(resourceName, "file_uri", fileUriUpdated),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttr(resourceName, "language", "SCALA"),
				resource.TestCheckResourceAttr(resourceName, "logs_bucket_uri", logsBucketUri),
				resource.TestCheckResourceAttr(resourceName, "metastore_id", metastoreId),
				resource.TestCheckResourceAttr(resourceName, "num_executors", "2"),
				resource.TestCheckResourceAttrSet(resourceName, "owner_principal_id"),
				resource.TestCheckResourceAttr(resourceName, "parameters.#", "1"),
				resource.TestCheckResourceAttr(resourceName, "parameters.0.name", "name2"),
				resource.TestCheckResourceAttr(resourceName, "parameters.0.value", "value2"),
				resource.TestCheckResourceAttrSet(resourceName, "private_endpoint_id"),
				resource.TestCheckResourceAttr(resourceName, "spark_version", "2.4.4"),
				resource.TestCheckResourceAttrSet(resourceName, "state"),
				resource.TestCheckResourceAttrSet(resourceName, "time_created"),
				resource.TestCheckResourceAttrSet(resourceName, "time_updated"),
				resource.TestCheckResourceAttr(resourceName, "warehouse_bucket_uri", warehouseBucketUri),
				resource.TestCheckResourceAttr(resourceName, "class_name", classNameUpdated),

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
				GenerateDataSourceFromRepresentationMap("oci_dataflow_applications", "test_applications", Optional, Update, dataFlowApplicationDataSourceRepresentation) +
				compartmentIdVariableStr + fileUriVariableStr + archiveUriVariableStr + fileUriVariableStrUpdated + logsBucketUriVariableStr + classNameStrUpdated + warehouseBucketUriVariableStr + metastoreIdVariableStr +
				dataFlowApplicationResourceDependencies + GenerateResourceFromRepresentationMap("oci_dataflow_application", "test_application", Optional, Update, RepresentationCopyWithNewProperties(dataFlowApplicationRepresentation, map[string]interface{}{
				"class_name": Representation{RepType: Optional, Create: `${var.dataflow_class_name_updated}`},
			})),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(datasourceName, "display_name", "displayName2"),
				resource.TestCheckResourceAttr(datasourceName, "applications.#", "1"),
				resource.TestCheckResourceAttr(datasourceName, "applications.0.compartment_id", compartmentId),
				resource.TestCheckResourceAttr(datasourceName, "applications.0.defined_tags.%", "1"),
				resource.TestCheckResourceAttr(datasourceName, "applications.0.display_name", "displayName2"),
				resource.TestCheckResourceAttr(datasourceName, "applications.0.freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(datasourceName, "applications.0.id"),
				resource.TestCheckResourceAttr(datasourceName, "applications.0.language", "SCALA"),
				resource.TestCheckResourceAttrSet(datasourceName, "applications.0.owner_principal_id"),
				resource.TestCheckResourceAttrSet(datasourceName, "applications.0.owner_user_name"),
				resource.TestCheckResourceAttr(datasourceName, "applications.0.spark_version", "2.4.4"),
				resource.TestCheckResourceAttrSet(datasourceName, "applications.0.state"),
				resource.TestCheckResourceAttrSet(datasourceName, "applications.0.time_created"),
				resource.TestCheckResourceAttrSet(datasourceName, "applications.0.time_updated"),
			),
		},
		// verify singular datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_dataflow_application", "test_application", Required, Create, dataFlowApplicationSingularDataSourceRepresentation) +
				compartmentIdVariableStr + fileUriVariableStr + archiveUriVariableStr + fileUriVariableStrUpdated + logsBucketUriVariableStr + classNameStrUpdated + warehouseBucketUriVariableStr + metastoreIdVariableStr + dataFlowApplicationResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_dataflow_application", "test_application", Optional, Update, RepresentationCopyWithNewProperties(dataFlowApplicationRepresentation, map[string]interface{}{
					"class_name": Representation{RepType: Optional, Create: `${var.dataflow_class_name_updated}`},
				})),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(singularDatasourceName, "application_id"),

				resource.TestCheckResourceAttrSet(singularDatasourceName, "archive_uri"),
				resource.TestCheckResourceAttr(singularDatasourceName, "arguments.#", "1"),
				resource.TestCheckResourceAttr(singularDatasourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(singularDatasourceName, "configuration.%", "1"),
				resource.TestCheckResourceAttr(singularDatasourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(singularDatasourceName, "description", "description2"),
				resource.TestCheckResourceAttr(singularDatasourceName, "display_name", "displayName2"),
				resource.TestCheckResourceAttr(singularDatasourceName, "driver_shape", "VM.Standard2.1"),
				resource.TestCheckResourceAttr(singularDatasourceName, "executor_shape", "VM.Standard2.1"),
				resource.TestCheckResourceAttr(singularDatasourceName, "file_uri", fileUriUpdated),
				resource.TestCheckResourceAttr(singularDatasourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "id"),
				resource.TestCheckResourceAttr(singularDatasourceName, "language", "SCALA"),
				resource.TestCheckResourceAttr(singularDatasourceName, "logs_bucket_uri", logsBucketUri),
				resource.TestCheckResourceAttr(singularDatasourceName, "metastore_id", metastoreId),
				resource.TestCheckResourceAttr(singularDatasourceName, "num_executors", "2"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "owner_user_name"),
				resource.TestCheckResourceAttr(singularDatasourceName, "parameters.#", "1"),
				resource.TestCheckResourceAttr(singularDatasourceName, "parameters.0.name", "name2"),
				resource.TestCheckResourceAttr(singularDatasourceName, "parameters.0.value", "value2"),
				resource.TestCheckResourceAttr(singularDatasourceName, "spark_version", "2.4.4"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "state"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "time_created"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "time_updated"),
				resource.TestCheckResourceAttr(singularDatasourceName, "warehouse_bucket_uri", warehouseBucketUri),
			),
		},
		// remove singular datasource from previous step so that it doesn't conflict with import tests
		{
			Config: config + compartmentIdVariableStr + fileUriVariableStr + archiveUriVariableStr + fileUriVariableStrUpdated + classNameStrUpdated + logsBucketUriVariableStr + warehouseBucketUriVariableStr + metastoreIdVariableStr + dataFlowApplicationResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_dataflow_application", "test_application", Optional, Update, RepresentationCopyWithNewProperties(dataFlowApplicationRepresentation, map[string]interface{}{
					"class_name": Representation{RepType: Optional, Create: `${var.dataflow_class_name_updated}`},
				})),
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

func testAccCheckDataflowApplicationDestroy(s *terraform.State) error {
	noResourceFound := true
	client := testAccProvider.Meta().(*OracleClients).dataFlowClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oci_dataflow_application" {
			noResourceFound = false
			request := oci_dataflow.GetApplicationRequest{}

			tmp := rs.Primary.ID
			request.ApplicationId = &tmp

			request.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "dataflow")

			response, err := client.GetApplication(context.Background(), request)

			if err == nil {
				deletedLifecycleStates := map[string]bool{
					string(oci_dataflow.ApplicationLifecycleStateDeleted): true,
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
	if !InSweeperExcludeList("DataflowApplication") {
		resource.AddTestSweepers("DataflowApplication", &resource.Sweeper{
			Name:         "DataflowApplication",
			Dependencies: DependencyGraph["application"],
			F:            sweepDataflowApplicationResource,
		})
	}
}

func sweepDataflowApplicationResource(compartment string) error {
	dataFlowClient := GetTestClients(&schema.ResourceData{}).dataFlowClient()
	applicationIds, err := dataFlowGetApplicationIds(compartment)
	if err != nil {
		return err
	}
	for _, applicationId := range applicationIds {
		if ok := SweeperDefaultResourceId[applicationId]; !ok {
			deleteApplicationRequest := oci_dataflow.DeleteApplicationRequest{}

			deleteApplicationRequest.ApplicationId = &applicationId

			deleteApplicationRequest.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "dataflow")
			_, error := dataFlowClient.DeleteApplication(context.Background(), deleteApplicationRequest)
			if error != nil {
				fmt.Printf("Error deleting Application %s %s, It is possible that the resource is already deleted. Please verify manually \n", applicationId, error)
				continue
			}
			WaitTillCondition(testAccProvider, &applicationId, dataFlowApplicationSweepWaitCondition, time.Duration(3*time.Minute),
				dataFlowApplicationSweepResponseFetchOperation, "dataflow", true)
		}
	}
	return nil
}

func dataFlowGetApplicationIds(compartment string) ([]string, error) {
	ids := GetResourceIdsToSweep(compartment, "ApplicationId")
	if ids != nil {
		return ids, nil
	}
	var resourceIds []string
	compartmentId := compartment
	dataFlowClient := GetTestClients(&schema.ResourceData{}).dataFlowClient()

	listApplicationsRequest := oci_dataflow.ListApplicationsRequest{}
	listApplicationsRequest.CompartmentId = &compartmentId
	listApplicationsResponse, err := dataFlowClient.ListApplications(context.Background(), listApplicationsRequest)

	if err != nil {
		return resourceIds, fmt.Errorf("Error getting Application list for compartment id : %s , %s \n", compartmentId, err)
	}
	for _, application := range listApplicationsResponse.Items {
		id := *application.Id
		resourceIds = append(resourceIds, id)
		AddResourceIdToSweeperResourceIdMap(compartmentId, "ApplicationId", id)
	}
	return resourceIds, nil
}

func dataFlowApplicationSweepWaitCondition(response common.OCIOperationResponse) bool {
	// Only stop if the resource is available beyond 3 mins. As there could be an issue for the sweeper to delete the resource and manual intervention required.
	if applicationResponse, ok := response.Response.(oci_dataflow.GetApplicationResponse); ok {
		return applicationResponse.LifecycleState != oci_dataflow.ApplicationLifecycleStateDeleted
	}
	return false
}

func dataFlowApplicationSweepResponseFetchOperation(client *OracleClients, resourceId *string, retryPolicy *common.RetryPolicy) error {
	_, err := client.dataFlowClient().GetApplication(context.Background(), oci_dataflow.GetApplicationRequest{
		ApplicationId: resourceId,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: retryPolicy,
		},
	})
	return err
}
