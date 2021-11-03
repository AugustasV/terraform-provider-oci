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
	oci_artifacts "github.com/oracle/oci-go-sdk/v49/artifacts"
	"github.com/oracle/oci-go-sdk/v49/common"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	ContainerRepositoryRequiredOnlyResource = ContainerRepositoryResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_artifacts_container_repository", "test_container_repository", Required, Create, containerRepositoryRepresentation)

	ContainerRepositoryResourceConfig = ContainerRepositoryResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_artifacts_container_repository", "test_container_repository", Optional, Update, containerRepositoryRepresentation)

	containerRepositorySingularDataSourceRepresentation = map[string]interface{}{
		"repository_id": Representation{RepType: Required, Create: `${oci_artifacts_container_repository.test_container_repository.id}`},
	}

	containerRepositoryDataSourceRepresentation = map[string]interface{}{
		"compartment_id":            Representation{RepType: Required, Create: `${var.compartment_id}`},
		"compartment_id_in_subtree": Representation{RepType: Optional, Create: `false`},
		"is_public":                 Representation{RepType: Optional, Create: `false`, Update: `true`},
		"repository_id":             Representation{RepType: Optional, Create: `${oci_artifacts_container_repository.test_container_repository.id}`},
		"state":                     Representation{RepType: Optional, Create: `AVAILABLE`},
		"filter":                    RepresentationGroup{Required, containerRepositoryDataSourceFilterRepresentation}}
	containerRepositoryDataSourceFilterRepresentation = map[string]interface{}{
		"name":   Representation{RepType: Required, Create: `id`},
		"values": Representation{RepType: Required, Create: []string{`${oci_artifacts_container_repository.test_container_repository.id}`}},
	}

	containerRepositoryRepresentation = map[string]interface{}{
		"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id}`},
		"display_name":   Representation{RepType: Required, Create: RandomString(15, CharsetLowerCaseWithoutDigits)},
		"is_immutable":   Representation{RepType: Optional, Create: `false`, Update: `false`},
		"is_public":      Representation{RepType: Optional, Create: `false`, Update: `true`},
		"readme":         RepresentationGroup{Optional, containerRepositoryReadmeRepresentation},
	}
	containerRepositoryReadmeRepresentation = map[string]interface{}{
		"content": Representation{RepType: Required, Create: `content`, Update: `content2`},
		"format":  Representation{RepType: Required, Create: `TEXT_MARKDOWN`, Update: `TEXT_PLAIN`},
	}

	ContainerRepositoryResourceDependencies = GenerateResourceFromRepresentationMap("oci_artifacts_repository", "test_repository", Required, Create, repositoryRepresentation)
)

// issue-routing-tag: artifacts/default
func TestArtifactsContainerRepositoryResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestArtifactsContainerRepositoryResource_basic")
	defer httpreplay.SaveScenario()

	config := ProviderTestConfig()

	compartmentId := GetEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	compartmentIdU := GetEnvSettingWithDefault("compartment_id_for_update", compartmentId)
	compartmentIdUVariableStr := fmt.Sprintf("variable \"compartment_id_for_update\" { default = \"%s\" }\n", compartmentIdU)

	resourceName := "oci_artifacts_container_repository.test_container_repository"
	datasourceName := "data.oci_artifacts_container_repositories.test_container_repositories"
	singularDatasourceName := "data.oci_artifacts_container_repository.test_container_repository"

	var resId, resId2 string
	// Save TF content to Create resource with optional properties. This has to be exactly the same as the config part in the "Create with optionals" step in the test.
	SaveConfigContent(config+compartmentIdVariableStr+ContainerRepositoryResourceDependencies+
		GenerateResourceFromRepresentationMap("oci_artifacts_container_repository", "test_container_repository", Optional, Create, containerRepositoryRepresentation), "artifacts", "containerRepository", t)

	ResourceTest(t, testAccCheckArtifactsContainerRepositoryDestroy, []resource.TestStep{
		// verify Create
		{
			Config: config + compartmentIdVariableStr + ContainerRepositoryResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_artifacts_container_repository", "test_container_repository", Required, Create, containerRepositoryRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttrSet(resourceName, "display_name"),

				func(s *terraform.State) (err error) {
					resId, err = FromInstanceState(s, resourceName, "id")
					return err
				},
			),
		},

		// delete before next Create
		{
			Config: config + compartmentIdVariableStr + ContainerRepositoryResourceDependencies,
		},
		// verify Create with optionals
		{
			Config: config + compartmentIdVariableStr + ContainerRepositoryResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_artifacts_container_repository", "test_container_repository", Optional, Create, containerRepositoryRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(resourceName, "billable_size_in_gbs"),
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttrSet(resourceName, "created_by"),
				resource.TestCheckResourceAttrSet(resourceName, "display_name"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttrSet(resourceName, "image_count"),
				resource.TestCheckResourceAttr(resourceName, "is_immutable", "false"),
				resource.TestCheckResourceAttr(resourceName, "is_public", "false"),
				resource.TestCheckResourceAttrSet(resourceName, "layer_count"),
				resource.TestCheckResourceAttrSet(resourceName, "layers_size_in_bytes"),
				resource.TestCheckResourceAttr(resourceName, "readme.#", "1"),
				resource.TestCheckResourceAttr(resourceName, "readme.0.content", "content"),
				resource.TestCheckResourceAttr(resourceName, "readme.0.format", "TEXT_MARKDOWN"),
				resource.TestCheckResourceAttrSet(resourceName, "state"),
				resource.TestCheckResourceAttrSet(resourceName, "time_created"),

				func(s *terraform.State) (err error) {
					resId, err = FromInstanceState(s, resourceName, "id")
					if isEnableExportCompartment, _ := strconv.ParseBool(GetEnvSettingWithDefault("enable_export_compartment", "true")); isEnableExportCompartment {
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
			Config: config + compartmentIdVariableStr + compartmentIdUVariableStr + ContainerRepositoryResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_artifacts_container_repository", "test_container_repository", Optional, Create,
					RepresentationCopyWithNewProperties(containerRepositoryRepresentation, map[string]interface{}{
						"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id_for_update}`},
					})),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(resourceName, "billable_size_in_gbs"),
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentIdU),
				resource.TestCheckResourceAttrSet(resourceName, "created_by"),
				resource.TestCheckResourceAttrSet(resourceName, "display_name"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttrSet(resourceName, "image_count"),
				resource.TestCheckResourceAttr(resourceName, "is_immutable", "false"),
				resource.TestCheckResourceAttr(resourceName, "is_public", "false"),
				resource.TestCheckResourceAttrSet(resourceName, "layer_count"),
				resource.TestCheckResourceAttrSet(resourceName, "layers_size_in_bytes"),
				resource.TestCheckResourceAttr(resourceName, "readme.#", "1"),
				resource.TestCheckResourceAttr(resourceName, "readme.0.content", "content"),
				resource.TestCheckResourceAttr(resourceName, "readme.0.format", "TEXT_MARKDOWN"),
				resource.TestCheckResourceAttrSet(resourceName, "state"),
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
			Config: config + compartmentIdVariableStr + ContainerRepositoryResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_artifacts_container_repository", "test_container_repository", Optional, Update, containerRepositoryRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(resourceName, "billable_size_in_gbs"),
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttrSet(resourceName, "created_by"),
				resource.TestCheckResourceAttrSet(resourceName, "display_name"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttrSet(resourceName, "image_count"),
				resource.TestCheckResourceAttr(resourceName, "is_immutable", "false"),
				resource.TestCheckResourceAttr(resourceName, "is_public", "true"),
				resource.TestCheckResourceAttrSet(resourceName, "layer_count"),
				resource.TestCheckResourceAttrSet(resourceName, "layers_size_in_bytes"),
				resource.TestCheckResourceAttr(resourceName, "readme.#", "1"),
				resource.TestCheckResourceAttr(resourceName, "readme.0.content", "content2"),
				resource.TestCheckResourceAttr(resourceName, "readme.0.format", "TEXT_PLAIN"),
				resource.TestCheckResourceAttrSet(resourceName, "state"),
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
				GenerateDataSourceFromRepresentationMap("oci_artifacts_container_repositories", "test_container_repositories", Optional, Update, containerRepositoryDataSourceRepresentation) +
				compartmentIdVariableStr + ContainerRepositoryResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_artifacts_container_repository", "test_container_repository", Optional, Update, containerRepositoryRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(datasourceName, "compartment_id_in_subtree", "false"),
				resource.TestCheckResourceAttrSet(resourceName, "display_name"),
				resource.TestCheckResourceAttr(datasourceName, "is_public", "true"),
				resource.TestCheckResourceAttrSet(datasourceName, "repository_id"),
				resource.TestCheckResourceAttr(datasourceName, "state", "AVAILABLE"),

				resource.TestCheckResourceAttr(datasourceName, "container_repository_collection.#", "1"),
				resource.TestCheckResourceAttr(datasourceName, "container_repository_collection.0.items.#", "1"),
				resource.TestCheckResourceAttrSet(datasourceName, "container_repository_collection.0.image_count"),
				resource.TestCheckResourceAttrSet(datasourceName, "container_repository_collection.0.layer_count"),
				resource.TestCheckResourceAttrSet(datasourceName, "container_repository_collection.0.layers_size_in_bytes"),
			),
		},
		// verify singular datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_artifacts_container_repository", "test_container_repository", Required, Create, containerRepositorySingularDataSourceRepresentation) +
				compartmentIdVariableStr + ContainerRepositoryResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(singularDatasourceName, "repository_id"),

				resource.TestCheckResourceAttrSet(singularDatasourceName, "billable_size_in_gbs"),
				resource.TestCheckResourceAttr(singularDatasourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "created_by"),
				resource.TestCheckResourceAttrSet(resourceName, "display_name"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "id"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "image_count"),
				resource.TestCheckResourceAttr(singularDatasourceName, "is_immutable", "false"),
				resource.TestCheckResourceAttr(singularDatasourceName, "is_public", "true"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "layer_count"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "layers_size_in_bytes"),
				resource.TestCheckResourceAttr(singularDatasourceName, "readme.#", "1"),
				resource.TestCheckResourceAttr(singularDatasourceName, "readme.0.content", "content2"),
				resource.TestCheckResourceAttr(singularDatasourceName, "readme.0.format", "TEXT_PLAIN"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "state"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "time_created"),
			),
		},
		// remove singular datasource from previous step so that it doesn't conflict with import tests
		{
			Config: config + compartmentIdVariableStr + ContainerRepositoryResourceConfig,
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

func testAccCheckArtifactsContainerRepositoryDestroy(s *terraform.State) error {
	noResourceFound := true
	client := TestAccProvider.Meta().(*OracleClients).artifactsClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oci_artifacts_container_repository" {
			noResourceFound = false
			request := oci_artifacts.GetContainerRepositoryRequest{}

			if value, ok := rs.Primary.Attributes["id"]; ok {
				request.RepositoryId = &value
			}

			request.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "artifacts")

			response, err := client.GetContainerRepository(context.Background(), request)

			if err == nil {
				deletedLifecycleStates := map[string]bool{
					string(oci_artifacts.ContainerRepositoryLifecycleStateDeleted): true,
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
		InitDependencyGraph()
	}
	if !InSweeperExcludeList("ArtifactsContainerRepository") {
		resource.AddTestSweepers("ArtifactsContainerRepository", &resource.Sweeper{
			Name:         "ArtifactsContainerRepository",
			Dependencies: DependencyGraph["containerRepository"],
			F:            sweepArtifactsContainerRepositoryResource,
		})
	}
}

func sweepArtifactsContainerRepositoryResource(compartment string) error {
	artifactsClient := GetTestClients(&schema.ResourceData{}).artifactsClient()
	containerRepositoryIds, err := getContainerRepositoryIds(compartment)
	if err != nil {
		return err
	}
	for _, containerRepositoryId := range containerRepositoryIds {
		if ok := SweeperDefaultResourceId[containerRepositoryId]; !ok {
			deleteContainerRepositoryRequest := oci_artifacts.DeleteContainerRepositoryRequest{}

			deleteContainerRepositoryRequest.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "artifacts")
			_, error := artifactsClient.DeleteContainerRepository(context.Background(), deleteContainerRepositoryRequest)
			if error != nil {
				fmt.Printf("Error deleting ContainerRepository %s %s, It is possible that the resource is already deleted. Please verify manually \n", containerRepositoryId, error)
				continue
			}
			WaitTillCondition(TestAccProvider, &containerRepositoryId, containerRepositorySweepWaitCondition, time.Duration(3*time.Minute),
				containerRepositorySweepResponseFetchOperation, "artifacts", true)
		}
	}
	return nil
}

func getContainerRepositoryIds(compartment string) ([]string, error) {
	ids := GetResourceIdsToSweep(compartment, "ContainerRepositoryId")
	if ids != nil {
		return ids, nil
	}
	var resourceIds []string
	compartmentId := compartment
	artifactsClient := GetTestClients(&schema.ResourceData{}).artifactsClient()

	listContainerRepositoriesRequest := oci_artifacts.ListContainerRepositoriesRequest{}
	listContainerRepositoriesRequest.CompartmentId = &compartmentId
	state := string(oci_artifacts.ContainerRepositoryLifecycleStateAvailable)
	listContainerRepositoriesRequest.LifecycleState = &state
	listContainerRepositoriesResponse, err := artifactsClient.ListContainerRepositories(context.Background(), listContainerRepositoriesRequest)

	if err != nil {
		return resourceIds, fmt.Errorf("Error getting ContainerRepository list for compartment id : %s , %s \n", compartmentId, err)
	}
	for _, containerRepository := range listContainerRepositoriesResponse.Items {
		id := *containerRepository.Id
		resourceIds = append(resourceIds, id)
		AddResourceIdToSweeperResourceIdMap(compartmentId, "ContainerRepositoryId", id)
	}
	return resourceIds, nil
}

func containerRepositorySweepWaitCondition(response common.OCIOperationResponse) bool {
	// Only stop if the resource is available beyond 3 mins. As there could be an issue for the sweeper to delete the resource and manual intervention required.
	if containerRepositoryResponse, ok := response.Response.(oci_artifacts.GetContainerRepositoryResponse); ok {
		return containerRepositoryResponse.LifecycleState != oci_artifacts.ContainerRepositoryLifecycleStateDeleted
	}
	return false
}

func containerRepositorySweepResponseFetchOperation(client *OracleClients, resourceId *string, retryPolicy *common.RetryPolicy) error {
	_, err := client.artifactsClient().GetContainerRepository(context.Background(), oci_artifacts.GetContainerRepositoryRequest{RequestMetadata: common.RequestMetadata{
		RetryPolicy: retryPolicy,
	},
	})
	return err
}
