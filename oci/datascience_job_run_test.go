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
	"github.com/oracle/oci-go-sdk/v51/common"
	oci_datascience "github.com/oracle/oci-go-sdk/v51/datascience"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	JobRunRequiredOnlyResource = JobRunResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_datascience_job_run", "test_job_run", Required, Create, jobRunRepresentation)

	JobRunResourceConfig = JobRunResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_datascience_job_run", "test_job_run", Optional, Update, jobRunRepresentation)

	jobRunSingularDataSourceRepresentation = map[string]interface{}{
		"job_run_id": Representation{RepType: Required, Create: `${oci_datascience_job_run.test_job_run.id}`},
	}

	jobRunDataSourceRepresentation = map[string]interface{}{
		"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id}`},
		"created_by":     Representation{RepType: Optional, Create: `${oci_datascience_job_run.test_job_run.created_by}`},
		"display_name":   Representation{RepType: Optional, Create: `displayName`, Update: `displayName2`},
		"id":             Representation{RepType: Optional, Create: `${oci_datascience_job_run.test_job_run.id}`},
		"job_id":         Representation{RepType: Optional, Create: `${oci_datascience_job.test_job.id}`},
		"state":          Representation{RepType: Optional, Create: `SUCCEEDED`},
		"filter":         RepresentationGroup{Required, jobRunDataSourceFilterRepresentation},
	}

	jobRunDataSourceFilterRepresentation = map[string]interface{}{
		"name":   Representation{RepType: Required, Create: `id`},
		"values": Representation{RepType: Required, Create: []string{`${oci_datascience_job_run.test_job_run.id}`}},
	}

	jobRunRepresentation = map[string]interface{}{
		"compartment_id":                     Representation{RepType: Required, Create: `${var.compartment_id}`},
		"job_id":                             Representation{RepType: Required, Create: `${oci_datascience_job.test_job.id}`},
		"project_id":                         Representation{RepType: Required, Create: `${oci_datascience_project.test_project.id}`},
		"defined_tags":                       Representation{RepType: Optional, Create: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "value")}`, Update: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "updatedValue")}`},
		"display_name":                       Representation{RepType: Optional, Create: `displayName`, Update: `displayName2`},
		"freeform_tags":                      Representation{RepType: Optional, Create: map[string]string{"Department": "Finance"}, Update: map[string]string{"Department": "Accounting"}},
		"asynchronous":                       Representation{RepType: Required, Create: `false`},
		"job_configuration_override_details": RepresentationGroup{Required, jobRunJobConfigurationOverrideDetailsRepresentation},
		"lifecycle":                          RepresentationGroup{Required, ignoreJobRunDefinedTagsChangesRepresentation},
	}
	jobRunJobConfigurationOverrideDetailsRepresentation = map[string]interface{}{
		"job_type":                   Representation{RepType: Required, Create: `DEFAULT`},
		"command_line_arguments":     Representation{RepType: Optional, Create: `commandLineArguments`},
		"environment_variables":      Representation{RepType: Required, Create: map[string]string{"environmentVariables": "environmentVariables"}},
		"maximum_runtime_in_minutes": Representation{RepType: Optional, Create: `10`},
	}

	ignoreJobRunDefinedTagsChangesRepresentation = map[string]interface{}{
		"ignore_changes": Representation{RepType: Required, Create: []string{`defined_tags`}},
	}

	JobRunResourceDependencies = GenerateDataSourceFromRepresentationMap("oci_core_shapes", "test_shapes", Required, Create, shapeDataSourceRepresentation) +
		GenerateResourceFromRepresentationMap("oci_core_subnet", "test_subnet", Required, Create, subnetRepresentation) +
		GenerateResourceFromRepresentationMap("oci_core_vcn", "test_vcn", Required, Create, vcnRepresentation) +
		GenerateResourceFromRepresentationMap("oci_datascience_job", "test_job", Required, Create, mlJobWithArtifactNoLogging) +
		GenerateResourceFromRepresentationMap("oci_datascience_project", "test_project", Required, Create, projectRepresentation) +
		DefinedTagsDependencies
)

// issue-routing-tag: datascience/default
func TestDatascienceJobRunResource_basic(t *testing.T) {
	t.Skip("Skip this test until service fixes it")
	httpreplay.SetScenario("TestDatascienceJobRunResource_basic")
	defer httpreplay.SaveScenario()

	provider := testAccProvider
	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	resourceName := "oci_datascience_job_run.test_job_run"
	datasourceName := "data.oci_datascience_job_runs.test_job_runs"
	singularDatasourceName := "data.oci_datascience_job_run.test_job_run"

	var resId, resId2 string
	// Save TF content to Create resource with optional properties. This has to be exactly the same as the config part in the "Create with optionals" step in the test.
	SaveConfigContent(config+compartmentIdVariableStr+JobRunResourceDependencies+
		GenerateResourceFromRepresentationMap("oci_datascience_job_run", "test_job_run", Optional, Create, jobRunRepresentation), "datascience", "jobRun", t)

	resource.Test(t, resource.TestCase{
		PreCheck: func() { testAccPreCheck(t) },
		Providers: map[string]terraform.ResourceProvider{
			"oci": provider,
		},
		CheckDestroy: testAccCheckDatascienceJobRunDestroy,
		Steps: []resource.TestStep{
			// verify Create
			{
				Config: config + compartmentIdVariableStr + JobRunResourceDependencies +
					GenerateResourceFromRepresentationMap("oci_datascience_job_run", "test_job_run", Required, Create, jobRunRepresentation),
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttrSet(resourceName, "job_id"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),

					func(s *terraform.State) (err error) {
						resId, err = FromInstanceState(s, resourceName, "id")
						return err
					},
				),
			},

			// delete before next Create
			{
				Config: config + compartmentIdVariableStr + JobRunResourceDependencies,
			},
			// verify Create with optionals
			{
				Config: config + compartmentIdVariableStr + JobRunResourceDependencies +
					GenerateResourceFromRepresentationMap("oci_datascience_job_run", "test_job_run", Optional, Create, jobRunRepresentation),
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
					resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "displayName"),
					resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "job_configuration_override_details.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "job_configuration_override_details.0.command_line_arguments", "commandLineArguments"),
					resource.TestCheckResourceAttr(resourceName, "job_configuration_override_details.0.environment_variables.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "job_configuration_override_details.0.job_type", "DEFAULT"),
					resource.TestCheckResourceAttr(resourceName, "job_configuration_override_details.0.maximum_runtime_in_minutes", "10"),
					resource.TestCheckResourceAttrSet(resourceName, "job_id"),
					resource.TestCheckResourceAttr(resourceName, "job_infrastructure_configuration_details.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "job_log_configuration_override_details.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "time_accepted"),

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

			// verify updates to updatable parameters
			{
				Config: config + compartmentIdVariableStr + JobRunResourceDependencies +
					GenerateResourceFromRepresentationMap("oci_datascience_job_run", "test_job_run", Optional, Update, jobRunRepresentation),
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttrSet(resourceName, "created_by"),
					resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "3"),
					resource.TestCheckResourceAttr(resourceName, "display_name", "displayName2"),
					resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
					resource.TestCheckResourceAttr(resourceName, "job_configuration_override_details.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "job_configuration_override_details.0.command_line_arguments", "commandLineArguments"),
					resource.TestCheckResourceAttr(resourceName, "job_configuration_override_details.0.environment_variables.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "job_configuration_override_details.0.job_type", "DEFAULT"),
					resource.TestCheckResourceAttr(resourceName, "job_configuration_override_details.0.maximum_runtime_in_minutes", "10"),
					resource.TestCheckResourceAttrSet(resourceName, "job_id"),
					resource.TestCheckResourceAttr(resourceName, "job_infrastructure_configuration_details.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "job_log_configuration_override_details.#", "0"),
					resource.TestCheckResourceAttrSet(resourceName, "project_id"),
					resource.TestCheckResourceAttrSet(resourceName, "state"),
					resource.TestCheckResourceAttrSet(resourceName, "time_accepted"),

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
					GenerateDataSourceFromRepresentationMap("oci_datascience_job_runs", "test_job_runs", Optional, Update, jobRunDataSourceRepresentation) +
					compartmentIdVariableStr + JobRunResourceDependencies +
					GenerateResourceFromRepresentationMap("oci_datascience_job_run", "test_job_run", Optional, Update, jobRunRepresentation),
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttrSet(datasourceName, "created_by"),
					resource.TestCheckResourceAttr(datasourceName, "display_name", "displayName2"),
					resource.TestCheckResourceAttrSet(datasourceName, "id"),
					resource.TestCheckResourceAttrSet(datasourceName, "job_id"),
					resource.TestCheckResourceAttr(datasourceName, "state", "SUCCEEDED"),

					resource.TestCheckResourceAttr(datasourceName, "job_runs.#", "1"),
					resource.TestCheckResourceAttr(datasourceName, "job_runs.0.compartment_id", compartmentId),
					resource.TestCheckResourceAttrSet(datasourceName, "job_runs.0.created_by"),
					resource.TestCheckResourceAttr(datasourceName, "job_runs.0.defined_tags.%", "3"),
					resource.TestCheckResourceAttr(datasourceName, "job_runs.0.display_name", "displayName2"),
					resource.TestCheckResourceAttr(datasourceName, "job_runs.0.freeform_tags.%", "1"),
					resource.TestCheckResourceAttrSet(datasourceName, "job_runs.0.id"),
					resource.TestCheckResourceAttrSet(datasourceName, "job_runs.0.job_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "job_runs.0.project_id"),
					resource.TestCheckResourceAttrSet(datasourceName, "job_runs.0.state"),
					resource.TestCheckResourceAttrSet(datasourceName, "job_runs.0.time_accepted"),
				),
			},
			// verify singular datasource
			{
				Config: config +
					GenerateDataSourceFromRepresentationMap("oci_datascience_job_run", "test_job_run", Required, Create, jobRunSingularDataSourceRepresentation) +
					compartmentIdVariableStr + JobRunResourceConfig,
				Check: ComposeAggregateTestCheckFuncWrapper(
					resource.TestCheckResourceAttrSet(singularDatasourceName, "job_run_id"),

					resource.TestCheckResourceAttr(singularDatasourceName, "compartment_id", compartmentId),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "created_by"),
					resource.TestCheckResourceAttr(singularDatasourceName, "defined_tags.%", "3"),
					resource.TestCheckResourceAttr(singularDatasourceName, "display_name", "displayName2"),
					resource.TestCheckResourceAttr(singularDatasourceName, "freeform_tags.%", "1"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "id"),
					resource.TestCheckResourceAttr(singularDatasourceName, "job_configuration_override_details.#", "1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "job_configuration_override_details.0.command_line_arguments", "commandLineArguments"),
					resource.TestCheckResourceAttr(singularDatasourceName, "job_configuration_override_details.0.environment_variables.%", "1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "job_configuration_override_details.0.job_type", "DEFAULT"),
					resource.TestCheckResourceAttr(singularDatasourceName, "job_configuration_override_details.0.maximum_runtime_in_minutes", "10"),
					resource.TestCheckResourceAttr(singularDatasourceName, "job_infrastructure_configuration_details.#", "1"),
					resource.TestCheckResourceAttr(singularDatasourceName, "job_log_configuration_override_details.#", "0"),
					resource.TestCheckResourceAttr(singularDatasourceName, "log_details.#", "0"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "state"),
					resource.TestCheckResourceAttrSet(singularDatasourceName, "time_accepted"),
				),
			},
			// remove singular datasource from previous step so that it doesn't conflict with import tests
			{
				Config: config + compartmentIdVariableStr + JobRunResourceConfig,
			},
			// verify resource import
			{
				Config:            config,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateVerifyIgnore: []string{
					"asynchronous",
				},
				ResourceName: resourceName,
			},
		},
	})
}

func testAccCheckDatascienceJobRunDestroy(s *terraform.State) error {
	noResourceFound := true
	client := testAccProvider.Meta().(*OracleClients).dataScienceClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oci_datascience_job_run" {
			noResourceFound = false
			request := oci_datascience.GetJobRunRequest{}

			tmp := rs.Primary.ID
			request.JobRunId = &tmp

			request.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "datascience")

			response, err := client.GetJobRun(context.Background(), request)

			if err == nil {
				deletedLifecycleStates := map[string]bool{
					string(oci_datascience.JobRunLifecycleStateDeleted): true,
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
	if !InSweeperExcludeList("DatascienceJobRun") {
		resource.AddTestSweepers("DatascienceJobRun", &resource.Sweeper{
			Name:         "DatascienceJobRun",
			Dependencies: DependencyGraph["jobRun"],
			F:            sweepDatascienceJobRunResource,
		})
	}
}

func sweepDatascienceJobRunResource(compartment string) error {
	dataScienceClient := GetTestClients(&schema.ResourceData{}).dataScienceClient()
	jobRunIds, err := getJobRunIds(compartment)
	if err != nil {
		return err
	}
	for _, jobRunId := range jobRunIds {
		if ok := SweeperDefaultResourceId[jobRunId]; !ok {
			deleteJobRunRequest := oci_datascience.DeleteJobRunRequest{}

			deleteJobRunRequest.JobRunId = &jobRunId

			deleteJobRunRequest.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "datascience")
			_, error := dataScienceClient.DeleteJobRun(context.Background(), deleteJobRunRequest)
			if error != nil {
				fmt.Printf("Error deleting JobRun %s %s, It is possible that the resource is already deleted. Please verify manually \n", jobRunId, error)
				continue
			}
			WaitTillCondition(testAccProvider, &jobRunId, jobRunSweepWaitCondition, time.Duration(3*time.Minute),
				jobRunSweepResponseFetchOperation, "datascience", true)
		}
	}
	return nil
}

func getJobRunIds(compartment string) ([]string, error) {
	ids := GetResourceIdsToSweep(compartment, "JobRunId")
	if ids != nil {
		return ids, nil
	}
	var resourceIds []string
	compartmentId := compartment
	dataScienceClient := GetTestClients(&schema.ResourceData{}).dataScienceClient()

	listJobRunsRequest := oci_datascience.ListJobRunsRequest{}
	listJobRunsRequest.CompartmentId = &compartmentId
	listJobRunsRequest.LifecycleState = oci_datascience.ListJobRunsLifecycleStateSucceeded
	listJobRunsResponse, err := dataScienceClient.ListJobRuns(context.Background(), listJobRunsRequest)

	if err != nil {
		return resourceIds, fmt.Errorf("Error getting JobRun list for compartment id : %s , %s \n", compartmentId, err)
	}
	for _, jobRun := range listJobRunsResponse.Items {
		id := *jobRun.Id
		resourceIds = append(resourceIds, id)
		AddResourceIdToSweeperResourceIdMap(compartmentId, "JobRunId", id)
	}
	return resourceIds, nil
}

func jobRunSweepWaitCondition(response common.OCIOperationResponse) bool {
	// Only stop if the resource is available beyond 3 mins. As there could be an issue for the sweeper to delete the resource and manual intervention required.
	if jobRunResponse, ok := response.Response.(oci_datascience.GetJobRunResponse); ok {
		return jobRunResponse.LifecycleState != oci_datascience.JobRunLifecycleStateDeleted
	}
	return false
}

func jobRunSweepResponseFetchOperation(client *OracleClients, resourceId *string, retryPolicy *common.RetryPolicy) error {
	_, err := client.dataScienceClient().GetJobRun(context.Background(), oci_datascience.GetJobRunRequest{
		JobRunId: resourceId,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: retryPolicy,
		},
	})
	return err
}
