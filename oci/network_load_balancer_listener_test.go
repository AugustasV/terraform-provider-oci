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
	"github.com/oracle/oci-go-sdk/v50/common"
	oci_network_load_balancer "github.com/oracle/oci-go-sdk/v50/networkloadbalancer"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	NlbListenerResourceConfig = NlbListenerResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_network_load_balancer_listener", "test_listener", Optional, Update, nlbListenerRepresentation)

	nlbListenerDataSourceRepresentation = map[string]interface{}{
		"network_load_balancer_id": Representation{RepType: Required, Create: `${oci_network_load_balancer_network_load_balancer.test_network_load_balancer.id}`},
	}

	nlbListenerSingularDataSourceRepresentation = map[string]interface{}{
		"listener_name":            Representation{RepType: Required, Create: `${oci_network_load_balancer_listener.test_listener.name}`},
		"network_load_balancer_id": Representation{RepType: Required, Create: `${oci_network_load_balancer_network_load_balancer.test_network_load_balancer.id}`},
	}

	nlbListenerRepresentation = map[string]interface{}{
		"default_backend_set_name": Representation{RepType: Required, Create: `${oci_network_load_balancer_backend_set.test_backend_set.name}`},
		"name":                     Representation{RepType: Required, Create: `example_listener`},
		"network_load_balancer_id": Representation{RepType: Required, Create: `${oci_network_load_balancer_network_load_balancer.test_network_load_balancer.id}`},
		"port":                     Representation{RepType: Required, Create: `10`, Update: `11`},
		"protocol":                 Representation{RepType: Required, Create: `UDP`, Update: `TCP`},
	}

	NlbListenerResourceDependencies = GenerateResourceFromRepresentationMap("oci_core_subnet", "test_subnet", Required, Create, subnetRepresentation) +
		GenerateResourceFromRepresentationMap("oci_core_vcn", "test_vcn", Required, Create, vcnRepresentation) +
		GenerateResourceFromRepresentationMap("oci_network_load_balancer_backend_set", "test_backend_set", Required, Create, nlbBackendSetRepresentation) +
		GenerateResourceFromRepresentationMap("oci_network_load_balancer_network_load_balancer", "test_network_load_balancer", Required, Create, networkLoadBalancerRepresentation)
)

// issue-routing-tag: network_load_balancer/default
func TestNetworkLoadBalancerListenerResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestNetworkLoadBalancerListenerResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	resourceName := "oci_network_load_balancer_listener.test_listener"
	datasourceName := "data.oci_network_load_balancer_listeners.test_listeners"
	singularDatasourceName := "data.oci_network_load_balancer_listener.test_listener"

	var resId, resId2 string

	ResourceTest(t, testAccCheckNetworkLoadBalancerListenerDestroy, []resource.TestStep{
		// verify Create
		{
			Config: config + compartmentIdVariableStr + NlbListenerResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_network_load_balancer_listener", "test_listener", Required, Create, nlbListenerRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(resourceName, "default_backend_set_name"),
				resource.TestCheckResourceAttr(resourceName, "name", "example_listener"),
				resource.TestCheckResourceAttrSet(resourceName, "network_load_balancer_id"),
				resource.TestCheckResourceAttr(resourceName, "port", "10"),
				resource.TestCheckResourceAttr(resourceName, "protocol", "UDP"),

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
			Config: config + compartmentIdVariableStr + NlbListenerResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_network_load_balancer_listener", "test_listener", Optional, Update, nlbListenerRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(resourceName, "default_backend_set_name"),
				resource.TestCheckResourceAttr(resourceName, "name", "example_listener"),
				resource.TestCheckResourceAttrSet(resourceName, "network_load_balancer_id"),
				resource.TestCheckResourceAttr(resourceName, "port", "11"),
				resource.TestCheckResourceAttr(resourceName, "protocol", "TCP"),

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
				GenerateDataSourceFromRepresentationMap("oci_network_load_balancer_listeners", "test_listeners", Optional, Update, nlbListenerDataSourceRepresentation) +
				compartmentIdVariableStr + NlbListenerResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(datasourceName, "network_load_balancer_id"),

				resource.TestCheckResourceAttr(datasourceName, "listener_collection.#", "1"),
				resource.TestCheckResourceAttr(datasourceName, "listener_collection.0.items.#", "1"),
			),
		},
		// verify singular datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_network_load_balancer_listener", "test_listener", Required, Create, nlbListenerSingularDataSourceRepresentation) +
				compartmentIdVariableStr + NlbListenerResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(singularDatasourceName, "listener_name"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "network_load_balancer_id"),

				resource.TestCheckResourceAttr(singularDatasourceName, "name", "example_listener"),
				resource.TestCheckResourceAttr(singularDatasourceName, "port", "11"),
				resource.TestCheckResourceAttr(singularDatasourceName, "protocol", "TCP"),
			),
		},
		// remove singular datasource from previous step so that it doesn't conflict with import tests
		{
			Config: config + compartmentIdVariableStr + NlbListenerResourceConfig,
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

func testAccCheckNetworkLoadBalancerListenerDestroy(s *terraform.State) error {
	noResourceFound := true
	client := testAccProvider.Meta().(*OracleClients).networkLoadBalancerClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oci_network_load_balancer_listener" {
			noResourceFound = false
			request := oci_network_load_balancer.GetListenerRequest{}

			if value, ok := rs.Primary.Attributes["name"]; ok {
				request.ListenerName = &value
			}

			if value, ok := rs.Primary.Attributes["network_load_balancer_id"]; ok {
				request.NetworkLoadBalancerId = &value
			}

			request.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "network_load_balancer")

			_, err := client.GetListener(context.Background(), request)

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
	if !InSweeperExcludeList("NetworkLoadBalancerListener") {
		resource.AddTestSweepers("NetworkLoadBalancerListener", &resource.Sweeper{
			Name:         "NetworkLoadBalancerListener",
			Dependencies: DependencyGraph["listener"],
			F:            sweepNetworkLoadBalancerListenerResource,
		})
	}
}

func sweepNetworkLoadBalancerListenerResource(compartment string) error {
	networkLoadBalancerClient := GetTestClients(&schema.ResourceData{}).networkLoadBalancerClient()
	listenerIds, err := getNetworkLoadBalancerListenerIds(compartment)
	if err != nil {
		return err
	}
	for _, listenerId := range listenerIds {
		if ok := SweeperDefaultResourceId[listenerId]; !ok {
			deleteListenerRequest := oci_network_load_balancer.DeleteListenerRequest{}

			deleteListenerRequest.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "network_load_balancer")
			_, error := networkLoadBalancerClient.DeleteListener(context.Background(), deleteListenerRequest)
			if error != nil {
				fmt.Printf("Error deleting Listener %s %s, It is possible that the resource is already deleted. Please verify manually \n", listenerId, error)
				continue
			}
		}
	}
	return nil
}

func getNetworkLoadBalancerListenerIds(compartment string) ([]string, error) {
	ids := GetResourceIdsToSweep(compartment, "ListenerId")
	if ids != nil {
		return ids, nil
	}
	var resourceIds []string
	compartmentId := compartment
	networkLoadBalancerClient := GetTestClients(&schema.ResourceData{}).networkLoadBalancerClient()

	listListenersRequest := oci_network_load_balancer.ListListenersRequest{}

	networkLoadBalancerIds, error := getNetworkLoadBalancerIds(compartment)
	if error != nil {
		return resourceIds, fmt.Errorf("Error getting networkLoadBalancerId required for Listener resource requests \n")
	}
	for _, networkLoadBalancerId := range networkLoadBalancerIds {
		listListenersRequest.NetworkLoadBalancerId = &networkLoadBalancerId

		listListenersResponse, err := networkLoadBalancerClient.ListListeners(context.Background(), listListenersRequest)

		if err != nil {
			return resourceIds, fmt.Errorf("Error getting Listener list for NLB id : %s , %s \n", networkLoadBalancerId, err)
		}
		for _, listener := range listListenersResponse.Items {
			id := *listener.Name
			resourceIds = append(resourceIds, id)
			AddResourceIdToSweeperResourceIdMap(compartmentId, "ListenerId", id)
			SweeperDefaultResourceId[*listener.Name] = true
		}

	}
	return resourceIds, nil
}
