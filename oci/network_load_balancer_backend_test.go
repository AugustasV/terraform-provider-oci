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
	"github.com/oracle/oci-go-sdk/v52/common"
	oci_network_load_balancer "github.com/oracle/oci-go-sdk/v52/networkloadbalancer"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	NlbBackendRequiredOnlyResource = NlbBackendResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_network_load_balancer_backend", "test_backend", Required, Create, nlbBackendRepresentation)

	NlbBackendResourceConfig = NlbBackendResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_network_load_balancer_backend", "test_backend", Optional, Update, nlbBackendRepresentation)

	nlbBackendDataSourceRepresentation = map[string]interface{}{
		"backend_set_name":         Representation{RepType: Required, Create: `${oci_network_load_balancer_backend_set.test_backend_set.name}`},
		"network_load_balancer_id": Representation{RepType: Required, Create: `${oci_network_load_balancer_network_load_balancer.test_network_load_balancer.id}`},
	}

	nlbBackendRepresentation = map[string]interface{}{
		"backend_set_name":         Representation{RepType: Required, Create: `${oci_network_load_balancer_backend_set.test_backend_set.name}`},
		"network_load_balancer_id": Representation{RepType: Required, Create: `${oci_network_load_balancer_network_load_balancer.test_network_load_balancer.id}`},
		"port":                     Representation{RepType: Required, Create: `10`},
		"ip_address":               Representation{RepType: Required, Create: `10.0.0.3`},
		"is_backup":                Representation{RepType: Optional, Create: `false`, Update: `true`},
		"is_drain":                 Representation{RepType: Optional, Create: `false`, Update: `true`},
		"is_offline":               Representation{RepType: Optional, Create: `false`, Update: `true`},
		"name":                     Representation{RepType: Optional, Create: `name`},
		"weight":                   Representation{RepType: Required, Create: `10`, Update: `11`},
	}

	NlbBackendResourceDependencies = GenerateResourceFromRepresentationMap("oci_core_subnet", "test_subnet", Required, Create, subnetRepresentation) +
		GenerateResourceFromRepresentationMap("oci_core_vcn", "test_vcn", Required, Create, vcnRepresentation) +
		GenerateResourceFromRepresentationMap("oci_network_load_balancer_backend_set", "test_backend_set", Required, Create, nlbBackendSetRepresentation) +
		GenerateResourceFromRepresentationMap("oci_network_load_balancer_network_load_balancer", "test_network_load_balancer", Required, Create, networkLoadBalancerRepresentation)
)

// issue-routing-tag: network_load_balancer/default
func TestNetworkLoadBalancerBackendResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestNetworkLoadBalancerBackendResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	resourceName := "oci_network_load_balancer_backend.test_backend"
	datasourceName := "data.oci_network_load_balancer_backends.test_backends"

	var resId, resId2 string

	ResourceTest(t, testAccCheckNetworkLoadBalancerBackendDestroy, []resource.TestStep{
		// verify Create
		{
			Config: config + compartmentIdVariableStr + NlbBackendResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_network_load_balancer_backend", "test_backend", Required, Create, nlbBackendRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(resourceName, "backend_set_name"),
				resource.TestCheckResourceAttrSet(resourceName, "network_load_balancer_id"),
				resource.TestCheckResourceAttr(resourceName, "port", "10"),

				func(s *terraform.State) (err error) {
					resId, err = FromInstanceState(s, resourceName, "id")
					return err
				},
			),
		},

		// delete before next Create
		{
			Config: config + compartmentIdVariableStr + NlbBackendResourceDependencies,
		},
		// verify Create with optionals
		{
			Config: config + compartmentIdVariableStr + NlbBackendResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_network_load_balancer_backend", "test_backend", Optional, Create, nlbBackendRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(resourceName, "backend_set_name"),
				resource.TestCheckResourceAttr(resourceName, "ip_address", "10.0.0.3"),
				resource.TestCheckResourceAttr(resourceName, "is_backup", "false"),
				resource.TestCheckResourceAttr(resourceName, "is_drain", "false"),
				resource.TestCheckResourceAttr(resourceName, "is_offline", "false"),
				resource.TestCheckResourceAttr(resourceName, "name", "name"),
				resource.TestCheckResourceAttrSet(resourceName, "network_load_balancer_id"),
				resource.TestCheckResourceAttr(resourceName, "port", "10"),
				resource.TestCheckResourceAttr(resourceName, "weight", "10"),

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
			Config: config + compartmentIdVariableStr + NlbBackendResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_network_load_balancer_backend", "test_backend", Optional, Update, nlbBackendRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(resourceName, "backend_set_name"),
				resource.TestCheckResourceAttr(resourceName, "ip_address", "10.0.0.3"),
				resource.TestCheckResourceAttr(resourceName, "is_backup", "true"),
				resource.TestCheckResourceAttr(resourceName, "is_drain", "true"),
				resource.TestCheckResourceAttr(resourceName, "is_offline", "true"),
				resource.TestCheckResourceAttr(resourceName, "name", "name"),
				resource.TestCheckResourceAttrSet(resourceName, "network_load_balancer_id"),
				resource.TestCheckResourceAttr(resourceName, "port", "10"),
				resource.TestCheckResourceAttr(resourceName, "weight", "11"),

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
				GenerateDataSourceFromRepresentationMap("oci_network_load_balancer_backends", "test_backends", Optional, Update, nlbBackendDataSourceRepresentation) +
				compartmentIdVariableStr + NlbBackendResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_network_load_balancer_backend", "test_backend", Optional, Update, nlbBackendRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(datasourceName, "backend_set_name"),
				resource.TestCheckResourceAttrSet(datasourceName, "network_load_balancer_id"),

				resource.TestCheckResourceAttr(datasourceName, "backend_collection.#", "1"),
				resource.TestCheckResourceAttr(datasourceName, "backend_collection.0.items.#", "1"),
			),
		},
		// remove singular datasource from previous step so that it doesn't conflict with import tests
		{
			Config: config + compartmentIdVariableStr + NlbBackendResourceConfig,
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

func testAccCheckNetworkLoadBalancerBackendDestroy(s *terraform.State) error {
	noResourceFound := true
	client := testAccProvider.Meta().(*OracleClients).networkLoadBalancerClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oci_network_load_balancer_backend" {
			noResourceFound = false
			request := oci_network_load_balancer.GetBackendRequest{}

			if value, ok := rs.Primary.Attributes["name"]; ok {
				request.BackendName = &value
			}

			if value, ok := rs.Primary.Attributes["backend_set_name"]; ok {
				request.BackendSetName = &value
			}

			if value, ok := rs.Primary.Attributes["network_load_balancer_id"]; ok {
				request.NetworkLoadBalancerId = &value
			}

			request.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "network_load_balancer")

			_, err := client.GetBackend(context.Background(), request)

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
	if !InSweeperExcludeList("NetworkLoadBalancerBackend") {
		resource.AddTestSweepers("NetworkLoadBalancerBackend", &resource.Sweeper{
			Name:         "NetworkLoadBalancerBackend",
			Dependencies: DependencyGraph["backend"],
			F:            sweepNetworkLoadBalancerBackendResource,
		})
	}
}

func sweepNetworkLoadBalancerBackendResource(compartment string) error {
	networkLoadBalancerClient := GetTestClients(&schema.ResourceData{}).networkLoadBalancerClient()
	backendIds, err := getNetworkLoadBalancerBackendIds(compartment)
	if err != nil {
		return err
	}
	for _, backendId := range backendIds {
		if ok := SweeperDefaultResourceId[backendId]; !ok {
			deleteBackendRequest := oci_network_load_balancer.DeleteBackendRequest{}

			deleteBackendRequest.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "network_load_balancer")
			_, error := networkLoadBalancerClient.DeleteBackend(context.Background(), deleteBackendRequest)
			if error != nil {
				fmt.Printf("Error deleting Backend %s %s, It is possible that the resource is already deleted. Please verify manually \n", backendId, error)
				continue
			}
		}
	}
	return nil
}

func getNetworkLoadBalancerBackendIds(compartment string) ([]string, error) {
	ids := GetResourceIdsToSweep(compartment, "BackendId")
	if ids != nil {
		return ids, nil
	}
	var resourceIds []string
	compartmentId := compartment
	networkLoadBalancerClient := GetTestClients(&schema.ResourceData{}).networkLoadBalancerClient()

	listBackendsRequest := oci_network_load_balancer.ListBackendsRequest{}

	backendSetNames, error := getNetworkLoadBalancerBackendSetIds(compartment)
	if error != nil {
		return resourceIds, fmt.Errorf("Error getting backendSetName required for Backend resource requests \n")
	}
	for _, backendSetName := range backendSetNames {
		listBackendsRequest.BackendSetName = &backendSetName

		networkLoadBalancerIds, error := getNetworkLoadBalancerIds(compartment)
		if error != nil {
			return resourceIds, fmt.Errorf("Error getting networkLoadBalancerId required for Backend resource requests \n")
		}
		for _, networkLoadBalancerId := range networkLoadBalancerIds {
			listBackendsRequest.NetworkLoadBalancerId = &networkLoadBalancerId

			listBackendsResponse, err := networkLoadBalancerClient.ListBackends(context.Background(), listBackendsRequest)

			if err != nil {
				return resourceIds, fmt.Errorf("Error getting Backend list for compartment id : %s , %s \n", compartmentId, err)
			}
			for _, backend := range listBackendsResponse.Items {
				id := *backend.Name
				resourceIds = append(resourceIds, id)
				AddResourceIdToSweeperResourceIdMap(compartmentId, "BackendId", id)
			}

		}
	}
	return resourceIds, nil
}
