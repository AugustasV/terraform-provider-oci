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
	"github.com/oracle/oci-go-sdk/v54/common"
	oci_load_balancer "github.com/oracle/oci-go-sdk/v54/loadbalancer"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	LoadBalancerRoutingPolicyResourceConfig = LoadBalancerRoutingPolicyResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_load_balancer_load_balancer_routing_policy", "test_load_balancer_routing_policy", Optional, Update, loadBalancerRoutingPolicyRepresentation)

	loadBalancerRoutingPolicySingularDataSourceRepresentation = map[string]interface{}{
		"load_balancer_id":    Representation{RepType: Required, Create: `${oci_load_balancer_load_balancer.test_load_balancer.id}`},
		"routing_policy_name": Representation{RepType: Required, Create: `example_routing_rules`},
	}

	loadBalancerRoutingPolicyDataSourceRepresentation = map[string]interface{}{
		"load_balancer_id": Representation{RepType: Required, Create: `${oci_load_balancer_load_balancer.test_load_balancer.id}`},
		"filter":           RepresentationGroup{Required, loadBalancerRoutingPolicyDataSourceFilterRepresentation}}
	loadBalancerRoutingPolicyDataSourceFilterRepresentation = map[string]interface{}{
		"name":   Representation{RepType: Required, Create: `name`},
		"values": Representation{RepType: Required, Create: []string{`${oci_load_balancer_load_balancer_routing_policy.test_load_balancer_routing_policy.name}`}},
	}

	loadBalancerRoutingPolicyRepresentation = map[string]interface{}{
		"condition_language_version": Representation{RepType: Required, Create: `V1`},
		"load_balancer_id":           Representation{RepType: Required, Create: `${oci_load_balancer_load_balancer.test_load_balancer.id}`},
		"name":                       Representation{RepType: Required, Create: `example_routing_rules`},
		"rules":                      RepresentationGroup{Required, loadBalancerRoutingPolicyRulesRepresentation},
	}
	loadBalancerRoutingPolicyRulesRepresentation = map[string]interface{}{
		"actions":   RepresentationGroup{Required, loadBalancerRoutingPolicyRulesActionsRepresentation},
		"condition": Representation{RepType: Required, Create: `all(http.request.url.path eq (i ''))`},
		"name":      Representation{RepType: Required, Create: `example_routing_rules`, Update: `name2`},
	}
	loadBalancerRoutingPolicyRulesActionsRepresentation = map[string]interface{}{
		"name":             Representation{RepType: Required, Create: `FORWARD_TO_BACKENDSET`, Update: `FORWARD_TO_BACKENDSET`},
		"backend_set_name": Representation{RepType: Required, Create: `${oci_load_balancer_backend_set.test_backend_set.name}`},
	}

	LoadBalancerRoutingPolicyResourceDependencies = GenerateResourceFromRepresentationMap("oci_load_balancer_backend_set", "test_backend_set", Required, Create, backendSetRepresentation) +
		GenerateResourceFromRepresentationMap("oci_load_balancer_load_balancer", "test_load_balancer", Required, Create, loadBalancerRepresentation) +
		LoadBalancerSubnetDependencies
)

// issue-routing-tag: load_balancer/default
func TestLoadBalancerLoadBalancerRoutingPolicyResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestLoadBalancerLoadBalancerRoutingPolicyResource_basic")
	defer httpreplay.SaveScenario()

	config := ProviderTestConfig()

	compartmentId := GetEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	resourceName := "oci_load_balancer_load_balancer_routing_policy.test_load_balancer_routing_policy"
	datasourceName := "data.oci_load_balancer_load_balancer_routing_policies.test_load_balancer_routing_policies"
	singularDatasourceName := "data.oci_load_balancer_load_balancer_routing_policy.test_load_balancer_routing_policy"

	var resId, resId2 string

	ResourceTest(t, testAccCheckLoadBalancerLoadBalancerRoutingPolicyDestroy, []resource.TestStep{
		// verify Create
		{
			Config: config + compartmentIdVariableStr + LoadBalancerRoutingPolicyResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_load_balancer_load_balancer_routing_policy", "test_load_balancer_routing_policy", Required, Create, loadBalancerRoutingPolicyRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "condition_language_version", "V1"),
				resource.TestCheckResourceAttrSet(resourceName, "load_balancer_id"),
				resource.TestCheckResourceAttr(resourceName, "name", "example_routing_rules"),
				resource.TestCheckResourceAttr(resourceName, "rules.#", "1"),
				resource.TestCheckResourceAttr(resourceName, "rules.0.actions.#", "1"),
				resource.TestCheckResourceAttr(resourceName, "rules.0.actions.0.name", "FORWARD_TO_BACKENDSET"),
				resource.TestCheckResourceAttr(resourceName, "rules.0.condition", "all(http.request.url.path eq (i ''))"),
				resource.TestCheckResourceAttr(resourceName, "rules.0.name", "example_routing_rules"),

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
			Config: config + compartmentIdVariableStr + LoadBalancerRoutingPolicyResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_load_balancer_load_balancer_routing_policy", "test_load_balancer_routing_policy", Optional, Update, loadBalancerRoutingPolicyRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "condition_language_version", "V1"),
				resource.TestCheckResourceAttrSet(resourceName, "load_balancer_id"),
				resource.TestCheckResourceAttr(resourceName, "name", "example_routing_rules"),
				resource.TestCheckResourceAttr(resourceName, "rules.#", "1"),
				resource.TestCheckResourceAttr(resourceName, "rules.0.actions.#", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "rules.0.actions.0.backend_set_name"),
				resource.TestCheckResourceAttr(resourceName, "rules.0.actions.0.name", "FORWARD_TO_BACKENDSET"),
				resource.TestCheckResourceAttr(resourceName, "rules.0.condition", "all(http.request.url.path eq (i ''))"),
				resource.TestCheckResourceAttr(resourceName, "rules.0.name", "name2"),

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
				GenerateDataSourceFromRepresentationMap("oci_load_balancer_load_balancer_routing_policies", "test_load_balancer_routing_policies", Optional, Update, loadBalancerRoutingPolicyDataSourceRepresentation) +
				compartmentIdVariableStr + LoadBalancerRoutingPolicyResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_load_balancer_load_balancer_routing_policy", "test_load_balancer_routing_policy", Optional, Update, loadBalancerRoutingPolicyRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(datasourceName, "load_balancer_id"),

				resource.TestCheckResourceAttr(datasourceName, "routing_policies.#", "1"),
				resource.TestCheckResourceAttr(datasourceName, "routing_policies.0.condition_language_version", "V1"),
				resource.TestCheckResourceAttr(datasourceName, "routing_policies.0.name", "example_routing_rules"),
				resource.TestCheckResourceAttr(datasourceName, "routing_policies.0.rules.#", "1"),
				resource.TestCheckResourceAttr(datasourceName, "routing_policies.0.rules.0.actions.#", "1"),
				resource.TestCheckResourceAttrSet(datasourceName, "routing_policies.0.rules.0.actions.0.backend_set_name"),
				resource.TestCheckResourceAttr(datasourceName, "routing_policies.0.rules.0.actions.0.name", "FORWARD_TO_BACKENDSET"),
				resource.TestCheckResourceAttr(datasourceName, "routing_policies.0.rules.0.condition", "all(http.request.url.path eq (i ''))"),
				resource.TestCheckResourceAttr(datasourceName, "routing_policies.0.rules.0.name", "name2"),
			),
		},
		// verify singular datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_load_balancer_load_balancer_routing_policy", "test_load_balancer_routing_policy", Required, Create, loadBalancerRoutingPolicySingularDataSourceRepresentation) +
				compartmentIdVariableStr + LoadBalancerRoutingPolicyResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(singularDatasourceName, "load_balancer_id"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "routing_policy_name"),

				resource.TestCheckResourceAttr(singularDatasourceName, "condition_language_version", "V1"),
				resource.TestCheckResourceAttr(singularDatasourceName, "name", "example_routing_rules"),
				resource.TestCheckResourceAttr(singularDatasourceName, "rules.#", "1"),
				resource.TestCheckResourceAttr(singularDatasourceName, "rules.0.actions.#", "1"),
				resource.TestCheckResourceAttr(singularDatasourceName, "rules.0.actions.0.name", "FORWARD_TO_BACKENDSET"),
				resource.TestCheckResourceAttr(singularDatasourceName, "rules.0.condition", "all(http.request.url.path eq (i ''))"),
				resource.TestCheckResourceAttr(singularDatasourceName, "rules.0.name", "name2"),
			),
		},
		// remove singular datasource from previous step so that it doesn't conflict with import tests
		{
			Config: config + compartmentIdVariableStr + LoadBalancerRoutingPolicyResourceConfig,
		},
		// verify resource import
		{
			Config:            config,
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				"state",
			},
			ResourceName: resourceName,
		},
	})
}

func testAccCheckLoadBalancerLoadBalancerRoutingPolicyDestroy(s *terraform.State) error {
	noResourceFound := true
	client := TestAccProvider.Meta().(*OracleClients).loadBalancerClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oci_load_balancer_load_balancer_routing_policy" {
			noResourceFound = false
			request := oci_load_balancer.GetRoutingPolicyRequest{}

			if value, ok := rs.Primary.Attributes["load_balancer_id"]; ok {
				request.LoadBalancerId = &value
			}

			if value, ok := rs.Primary.Attributes["name"]; ok {
				request.RoutingPolicyName = &value
			}

			request.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "load_balancer")

			_, err := client.GetRoutingPolicy(context.Background(), request)

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
		InitDependencyGraph()
	}
	if !InSweeperExcludeList("LoadBalancerLoadBalancerRoutingPolicy") {
		resource.AddTestSweepers("LoadBalancerLoadBalancerRoutingPolicy", &resource.Sweeper{
			Name:         "LoadBalancerLoadBalancerRoutingPolicy",
			Dependencies: DependencyGraph["loadBalancerRoutingPolicy"],
			F:            sweepLoadBalancerLoadBalancerRoutingPolicyResource,
		})
	}
}

func sweepLoadBalancerLoadBalancerRoutingPolicyResource(compartment string) error {
	loadBalancerClient := GetTestClients(&schema.ResourceData{}).loadBalancerClient()
	loadBalancerRoutingPolicyIds, err := getLoadBalancerRoutingPolicyIds(compartment)
	if err != nil {
		return err
	}
	for _, loadBalancerRoutingPolicyId := range loadBalancerRoutingPolicyIds {
		if ok := SweeperDefaultResourceId[loadBalancerRoutingPolicyId]; !ok {
			deleteRoutingPolicyRequest := oci_load_balancer.DeleteRoutingPolicyRequest{}

			deleteRoutingPolicyRequest.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "load_balancer")
			_, error := loadBalancerClient.DeleteRoutingPolicy(context.Background(), deleteRoutingPolicyRequest)
			if error != nil {
				fmt.Printf("Error deleting LoadBalancerRoutingPolicy %s %s, It is possible that the resource is already deleted. Please verify manually \n", loadBalancerRoutingPolicyId, error)
				continue
			}
		}
	}
	return nil
}

func getLoadBalancerRoutingPolicyIds(compartment string) ([]string, error) {
	ids := GetResourceIdsToSweep(compartment, "LoadBalancerRoutingPolicyId")
	if ids != nil {
		return ids, nil
	}
	var resourceIds []string
	compartmentId := compartment
	loadBalancerClient := GetTestClients(&schema.ResourceData{}).loadBalancerClient()

	listRoutingPoliciesRequest := oci_load_balancer.ListRoutingPoliciesRequest{}

	loadBalancerIds, error := getLoadBalancerIds(compartment)
	if error != nil {
		return resourceIds, fmt.Errorf("Error getting loadBalancerId required for LoadBalancerRoutingPolicy resource requests \n")
	}
	for _, loadBalancerId := range loadBalancerIds {
		listRoutingPoliciesRequest.LoadBalancerId = &loadBalancerId

		listRoutingPoliciesResponse, err := loadBalancerClient.ListRoutingPolicies(context.Background(), listRoutingPoliciesRequest)

		if err != nil {
			return resourceIds, fmt.Errorf("Error getting LoadBalancerRoutingPolicy list for compartment id : %s , %s \n", compartmentId, err)
		}
		for _, loadBalancerRoutingPolicy := range listRoutingPoliciesResponse.Items {
			id := *loadBalancerRoutingPolicy.Name
			resourceIds = append(resourceIds, id)
			AddResourceIdToSweeperResourceIdMap(compartmentId, "LoadBalancerRoutingPolicyId", id)
		}

	}
	return resourceIds, nil
}
