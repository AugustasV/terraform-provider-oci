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
	"github.com/oracle/oci-go-sdk/v54/common"
	oci_core "github.com/oracle/oci-go-sdk/v54/core"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	DrgRouteDistributionRequiredOnlyResource = DrgRouteDistributionResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_core_drg_route_distribution", "test_drg_route_distribution", Required, Create, drgRouteDistributionRepresentation)

	DrgRouteDistributionResourceConfig = DrgRouteDistributionResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_core_drg_route_distribution", "test_drg_route_distribution", Optional, Update, drgRouteDistributionRepresentation)

	drgRouteDistributionSingularDataSourceRepresentation = map[string]interface{}{
		"drg_route_distribution_id": Representation{RepType: Required, Create: `${oci_core_drg_route_distribution.test_drg_route_distribution.id}`},
	}

	drgRouteDistributionDataSourceRepresentation = map[string]interface{}{
		"drg_id":       Representation{RepType: Required, Create: `${oci_core_drg.test_drg.id}`},
		"display_name": Representation{RepType: Optional, Create: `displayName`, Update: `displayName2`},
		"state":        Representation{RepType: Optional, Create: `AVAILABLE`},
		"filter":       RepresentationGroup{Required, drgRouteDistributionDataSourceFilterRepresentation}}
	drgRouteDistributionDataSourceFilterRepresentation = map[string]interface{}{
		"name":   Representation{RepType: Required, Create: `id`},
		"values": Representation{RepType: Required, Create: []string{`${oci_core_drg_route_distribution.test_drg_route_distribution.id}`}},
	}

	drgRouteDistributionRepresentation = map[string]interface{}{
		"distribution_type": Representation{RepType: Required, Create: `IMPORT`},
		"drg_id":            Representation{RepType: Required, Create: `${oci_core_drg.test_drg.id}`},
		"defined_tags":      Representation{RepType: Optional, Create: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "value")}`, Update: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "updatedValue")}`},
		"display_name":      Representation{RepType: Optional, Create: `displayName`, Update: `displayName2`},
		"freeform_tags":     Representation{RepType: Optional, Create: map[string]string{"Department": "Finance"}, Update: map[string]string{"Department": "Accounting"}},
		"lifecycle":         RepresentationGroup{Required, ignoreChangesLBRepresentation},
	}

	DrgRouteDistributionResourceDependencies = GenerateResourceFromRepresentationMap("oci_core_drg", "test_drg", Required, Create, drgRepresentation) +
		DefinedTagsDependencies
)

// issue-routing-tag: core/pnp
func TestCoreDrgRouteDistributionResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestCoreDrgRouteDistributionResource_basic")
	defer httpreplay.SaveScenario()

	config := ProviderTestConfig()

	compartmentId := GetEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	resourceName := "oci_core_drg_route_distribution.test_drg_route_distribution"
	datasourceName := "data.oci_core_drg_route_distributions.test_drg_route_distributions"
	singularDatasourceName := "data.oci_core_drg_route_distribution.test_drg_route_distribution"

	var resId, resId2 string
	// Save TF content to Create resource with optional properties. This has to be exactly the same as the config part in the "Create with optionals" step in the test.
	SaveConfigContent(config+compartmentIdVariableStr+DrgRouteDistributionResourceDependencies+
		GenerateResourceFromRepresentationMap("oci_core_drg_route_distribution", "test_drg_route_distribution", Optional, Create, drgRouteDistributionRepresentation), "core", "drgRouteDistribution", t)

	ResourceTest(t, testAccCheckCoreDrgRouteDistributionDestroy, []resource.TestStep{
		// verify Create
		{
			Config: config + compartmentIdVariableStr + DrgRouteDistributionResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_core_drg_route_distribution", "test_drg_route_distribution", Required, Create, drgRouteDistributionRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "distribution_type", "IMPORT"),
				resource.TestCheckResourceAttrSet(resourceName, "drg_id"),

				func(s *terraform.State) (err error) {
					resId, err = FromInstanceState(s, resourceName, "id")
					return err
				},
			),
		},

		// delete before next Create
		{
			Config: config + compartmentIdVariableStr + DrgRouteDistributionResourceDependencies,
		},
		// verify Create with optionals
		{
			Config: config + compartmentIdVariableStr + DrgRouteDistributionResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_core_drg_route_distribution", "test_drg_route_distribution", Optional, Create, drgRouteDistributionRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(resourceName, "compartment_id"),
				resource.TestCheckResourceAttr(resourceName, "display_name", "displayName"),
				resource.TestCheckResourceAttr(resourceName, "distribution_type", "IMPORT"),
				resource.TestCheckResourceAttrSet(resourceName, "drg_id"),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
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

		// verify updates to updatable parameters
		{
			Config: config + compartmentIdVariableStr + DrgRouteDistributionResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_core_drg_route_distribution", "test_drg_route_distribution", Optional, Update, drgRouteDistributionRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(resourceName, "compartment_id"),
				resource.TestCheckResourceAttr(resourceName, "display_name", "displayName2"),
				resource.TestCheckResourceAttr(resourceName, "distribution_type", "IMPORT"),
				resource.TestCheckResourceAttrSet(resourceName, "drg_id"),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
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
				GenerateDataSourceFromRepresentationMap("oci_core_drg_route_distributions", "test_drg_route_distributions", Optional, Update, drgRouteDistributionDataSourceRepresentation) +
				compartmentIdVariableStr + DrgRouteDistributionResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_core_drg_route_distribution", "test_drg_route_distribution", Optional, Update, drgRouteDistributionRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(datasourceName, "display_name", "displayName2"),
				resource.TestCheckResourceAttrSet(datasourceName, "drg_id"),
				resource.TestCheckResourceAttr(datasourceName, "state", "AVAILABLE"),

				resource.TestCheckResourceAttr(datasourceName, "drg_route_distributions.#", "1"),
				resource.TestCheckResourceAttrSet(datasourceName, "drg_route_distributions.0.compartment_id"),
				resource.TestCheckResourceAttr(datasourceName, "drg_route_distributions.0.display_name", "displayName2"),
				resource.TestCheckResourceAttr(datasourceName, "drg_route_distributions.0.distribution_type", "IMPORT"),
				resource.TestCheckResourceAttrSet(datasourceName, "drg_route_distributions.0.drg_id"),
				resource.TestCheckResourceAttr(datasourceName, "drg_route_distributions.0.freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(datasourceName, "drg_route_distributions.0.id"),
				resource.TestCheckResourceAttrSet(datasourceName, "drg_route_distributions.0.state"),
				resource.TestCheckResourceAttrSet(datasourceName, "drg_route_distributions.0.time_created"),
			),
		},
		// verify singular datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_core_drg_route_distribution", "test_drg_route_distribution", Required, Create, drgRouteDistributionSingularDataSourceRepresentation) +
				compartmentIdVariableStr + DrgRouteDistributionResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(singularDatasourceName, "drg_route_distribution_id"),

				resource.TestCheckResourceAttrSet(singularDatasourceName, "compartment_id"),
				resource.TestCheckResourceAttr(singularDatasourceName, "display_name", "displayName2"),
				resource.TestCheckResourceAttr(singularDatasourceName, "distribution_type", "IMPORT"),
				resource.TestCheckResourceAttr(singularDatasourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "id"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "state"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "time_created"),
			),
		},
		// remove singular datasource from previous step so that it doesn't conflict with import tests
		{
			Config: config + compartmentIdVariableStr + DrgRouteDistributionResourceConfig,
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

func testAccCheckCoreDrgRouteDistributionDestroy(s *terraform.State) error {
	noResourceFound := true
	client := TestAccProvider.Meta().(*OracleClients).virtualNetworkClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oci_core_drg_route_distribution" {
			noResourceFound = false
			request := oci_core.GetDrgRouteDistributionRequest{}

			tmp := rs.Primary.ID
			request.DrgRouteDistributionId = &tmp

			request.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "core")

			response, err := client.GetDrgRouteDistribution(context.Background(), request)

			if err == nil {
				deletedLifecycleStates := map[string]bool{
					string(oci_core.DrgRouteDistributionLifecycleStateTerminated): true,
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
	if !InSweeperExcludeList("CoreDrgRouteDistribution") {
		resource.AddTestSweepers("CoreDrgRouteDistribution", &resource.Sweeper{
			Name:         "CoreDrgRouteDistribution",
			Dependencies: DependencyGraph["drgRouteDistribution"],
			F:            sweepCoreDrgRouteDistributionResource,
		})
	}
}

func sweepCoreDrgRouteDistributionResource(compartment string) error {
	virtualNetworkClient := GetTestClients(&schema.ResourceData{}).virtualNetworkClient()
	drgRouteDistributionIds, err := getDrgRouteDistributionIds(compartment)
	if err != nil {
		return err
	}
	for _, drgRouteDistributionId := range drgRouteDistributionIds {
		if ok := SweeperDefaultResourceId[drgRouteDistributionId]; !ok {
			deleteDrgRouteDistributionRequest := oci_core.DeleteDrgRouteDistributionRequest{}

			deleteDrgRouteDistributionRequest.DrgRouteDistributionId = &drgRouteDistributionId

			deleteDrgRouteDistributionRequest.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "core")
			_, error := virtualNetworkClient.DeleteDrgRouteDistribution(context.Background(), deleteDrgRouteDistributionRequest)
			if error != nil {
				fmt.Printf("Error deleting DrgRouteDistribution %s %s, It is possible that the resource is already deleted. Please verify manually \n", drgRouteDistributionId, error)
				continue
			}
			WaitTillCondition(TestAccProvider, &drgRouteDistributionId, drgRouteDistributionSweepWaitCondition, time.Duration(3*time.Minute),
				drgRouteDistributionSweepResponseFetchOperation, "core", true)
		}
	}
	return nil
}

func getDrgRouteDistributionIds(compartment string) ([]string, error) {
	ids := GetResourceIdsToSweep(compartment, "DrgRouteDistributionId")
	if ids != nil {
		return ids, nil
	}
	var resourceIds []string
	compartmentId := compartment
	virtualNetworkClient := GetTestClients(&schema.ResourceData{}).virtualNetworkClient()

	listDrgRouteDistributionsRequest := oci_core.ListDrgRouteDistributionsRequest{}
	// listDrgRouteDistributionsRequest.CompartmentId = &compartmentId

	drgIds, error := getDrgIds(compartment)
	if error != nil {
		return resourceIds, fmt.Errorf("Error getting drgId required for DrgRouteDistribution resource requests \n")
	}
	for _, drgId := range drgIds {
		listDrgRouteDistributionsRequest.DrgId = &drgId

		listDrgRouteDistributionsRequest.LifecycleState = oci_core.DrgRouteDistributionLifecycleStateAvailable
		listDrgRouteDistributionsResponse, err := virtualNetworkClient.ListDrgRouteDistributions(context.Background(), listDrgRouteDistributionsRequest)

		if err != nil {
			return resourceIds, fmt.Errorf("Error getting DrgRouteDistribution list for compartment id : %s , %s \n", compartmentId, err)
		}
		for _, drgRouteDistribution := range listDrgRouteDistributionsResponse.Items {
			id := *drgRouteDistribution.Id
			resourceIds = append(resourceIds, id)
			AddResourceIdToSweeperResourceIdMap(compartmentId, "DrgRouteDistributionId", id)
		}

	}
	return resourceIds, nil
}

func drgRouteDistributionSweepWaitCondition(response common.OCIOperationResponse) bool {
	// Only stop if the resource is available beyond 3 mins. As there could be an issue for the sweeper to delete the resource and manual intervention required.
	if drgRouteDistributionResponse, ok := response.Response.(oci_core.GetDrgRouteDistributionResponse); ok {
		return drgRouteDistributionResponse.LifecycleState != oci_core.DrgRouteDistributionLifecycleStateTerminated
	}
	return false
}

func drgRouteDistributionSweepResponseFetchOperation(client *OracleClients, resourceId *string, retryPolicy *common.RetryPolicy) error {
	_, err := client.virtualNetworkClient().GetDrgRouteDistribution(context.Background(), oci_core.GetDrgRouteDistributionRequest{
		DrgRouteDistributionId: resourceId,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: retryPolicy,
		},
	})
	return err
}
