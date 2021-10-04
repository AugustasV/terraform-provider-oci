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
	"github.com/oracle/oci-go-sdk/v48/common"
	oci_core "github.com/oracle/oci-go-sdk/v48/core"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	NatGatewayRequiredOnlyResource = NatGatewayResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_core_nat_gateway", "test_nat_gateway", Required, Create, natGatewayRepresentation)

	NatGatewayResourceConfig = NatGatewayResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_core_nat_gateway", "test_nat_gateway", Optional, Update, natGatewayRepresentation)

	natGatewaySingularDataSourceRepresentation = map[string]interface{}{
		"nat_gateway_id": Representation{RepType: Required, Create: `${oci_core_nat_gateway.test_nat_gateway.id}`},
	}

	natGatewayDataSourceRepresentation = map[string]interface{}{
		"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id}`},
		"display_name":   Representation{RepType: Optional, Create: `displayName`, Update: `displayName2`},
		"state":          Representation{RepType: Optional, Create: `AVAILABLE`},
		"vcn_id":         Representation{RepType: Optional, Create: `${oci_core_vcn.test_vcn.id}`},
		"filter":         RepresentationGroup{Required, natGatewayDataSourceFilterRepresentation}}
	natGatewayDataSourceFilterRepresentation = map[string]interface{}{
		"name":   Representation{RepType: Required, Create: `id`},
		"values": Representation{RepType: Required, Create: []string{`${oci_core_nat_gateway.test_nat_gateway.id}`}},
	}

	natGatewayRepresentation = map[string]interface{}{
		"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id}`},
		"vcn_id":         Representation{RepType: Required, Create: `${oci_core_vcn.test_vcn.id}`},
		"block_traffic":  Representation{RepType: Optional, Create: `false`, Update: `true`},
		"defined_tags":   Representation{RepType: Optional, Create: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "value")}`, Update: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "updatedValue")}`},
		"display_name":   Representation{RepType: Optional, Create: `displayName`, Update: `displayName2`},
		"freeform_tags":  Representation{RepType: Optional, Create: map[string]string{"Department": "Finance"}, Update: map[string]string{"Department": "Accounting"}},
		"public_ip_id":   Representation{RepType: Optional, Create: `${oci_core_public_ip.test_public_ip.id}`},
	}

	NatGatewayResourceDependencies = GenerateResourceFromRepresentationMap("oci_core_public_ip", "test_public_ip", Required, Create,
		RepresentationCopyWithNewProperties(publicIpRepresentation, map[string]interface{}{
			"public_ip_pool_id": Representation{RepType: Required, Create: `${oci_core_public_ip_pool.test_public_ip_pool.id}`},
		})) +
		GenerateResourceFromRepresentationMap("oci_core_public_ip_pool_capacity", "test_public_ip_pool_capacity", Required, Create, publicIpPoolCapacityRepresentation) +
		GenerateResourceFromRepresentationMap("oci_core_public_ip_pool", "test_public_ip_pool", Required, Create, publicIpPoolRepresentation) +
		GenerateResourceFromRepresentationMap("oci_core_vcn", "test_vcn", Required, Create, vcnRepresentation) +
		DefinedTagsDependencies + byoipRangeIdVariableStr + publicIpPoolCidrBlockVariableStr
)

// issue-routing-tag: core/pnp
func TestCoreNatGatewayResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestCoreNatGatewayResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	compartmentIdU := getEnvSettingWithDefault("compartment_id_for_update", compartmentId)
	compartmentIdUVariableStr := fmt.Sprintf("variable \"compartment_id_for_update\" { default = \"%s\" }\n", compartmentIdU)

	resourceName := "oci_core_nat_gateway.test_nat_gateway"
	datasourceName := "data.oci_core_nat_gateways.test_nat_gateways"
	singularDatasourceName := "data.oci_core_nat_gateway.test_nat_gateway"

	var resId, resId2 string
	// Save TF content to Create resource with optional properties. This has to be exactly the same as the config part in the "Create with optionals" step in the test.
	SaveConfigContent(config+compartmentIdVariableStr+NatGatewayResourceDependencies+
		GenerateResourceFromRepresentationMap("oci_core_nat_gateway", "test_nat_gateway", Optional, Create, natGatewayRepresentation), "core", "natGateway", t)

	ResourceTest(t, testAccCheckCoreNatGatewayDestroy, []resource.TestStep{
		// verify Create
		{
			Config: config + compartmentIdVariableStr + NatGatewayResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_core_nat_gateway", "test_nat_gateway", Required, Create, natGatewayRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),

				func(s *terraform.State) (err error) {
					resId, err = FromInstanceState(s, resourceName, "id")
					return err
				},
			),
		},

		// delete before next Create
		{
			Config: config + compartmentIdVariableStr + NatGatewayResourceDependencies,
		},
		// verify Create with optionals
		{
			Config: config + compartmentIdVariableStr + NatGatewayResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_core_nat_gateway", "test_nat_gateway", Optional, Create, natGatewayRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "block_traffic", "false"),
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "display_name", "displayName"),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttrSet(resourceName, "nat_ip"),
				resource.TestCheckResourceAttrSet(resourceName, "public_ip_id"),
				resource.TestCheckResourceAttrSet(resourceName, "state"),
				resource.TestCheckResourceAttrSet(resourceName, "time_created"),
				resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),

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
			Config: config + compartmentIdVariableStr + compartmentIdUVariableStr + NatGatewayResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_core_nat_gateway", "test_nat_gateway", Optional, Create,
					RepresentationCopyWithNewProperties(natGatewayRepresentation, map[string]interface{}{
						"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id_for_update}`},
					})),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "block_traffic", "false"),
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentIdU),
				resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "display_name", "displayName"),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttrSet(resourceName, "nat_ip"),
				resource.TestCheckResourceAttrSet(resourceName, "public_ip_id"),
				resource.TestCheckResourceAttrSet(resourceName, "state"),
				resource.TestCheckResourceAttrSet(resourceName, "time_created"),
				resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),

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
			Config: config + compartmentIdVariableStr + NatGatewayResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_core_nat_gateway", "test_nat_gateway", Optional, Update, natGatewayRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "block_traffic", "true"),
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "display_name", "displayName2"),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttrSet(resourceName, "nat_ip"),
				resource.TestCheckResourceAttrSet(resourceName, "public_ip_id"),
				resource.TestCheckResourceAttrSet(resourceName, "state"),
				resource.TestCheckResourceAttrSet(resourceName, "time_created"),
				resource.TestCheckResourceAttrSet(resourceName, "vcn_id"),

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
				GenerateDataSourceFromRepresentationMap("oci_core_nat_gateways", "test_nat_gateways", Optional, Update, natGatewayDataSourceRepresentation) +
				compartmentIdVariableStr + NatGatewayResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_core_nat_gateway", "test_nat_gateway", Optional, Update, natGatewayRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(datasourceName, "display_name", "displayName2"),
				resource.TestCheckResourceAttr(datasourceName, "state", "AVAILABLE"),
				resource.TestCheckResourceAttrSet(datasourceName, "vcn_id"),

				resource.TestCheckResourceAttr(datasourceName, "nat_gateways.#", "1"),
				resource.TestCheckResourceAttr(datasourceName, "nat_gateways.0.block_traffic", "true"),
				resource.TestCheckResourceAttr(datasourceName, "nat_gateways.0.compartment_id", compartmentId),
				resource.TestCheckResourceAttr(datasourceName, "nat_gateways.0.defined_tags.%", "1"),
				resource.TestCheckResourceAttr(datasourceName, "nat_gateways.0.display_name", "displayName2"),
				resource.TestCheckResourceAttr(datasourceName, "nat_gateways.0.freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(datasourceName, "nat_gateways.0.id"),
				resource.TestCheckResourceAttrSet(datasourceName, "nat_gateways.0.nat_ip"),
				resource.TestCheckResourceAttrSet(datasourceName, "nat_gateways.0.public_ip_id"),
				resource.TestCheckResourceAttrSet(datasourceName, "nat_gateways.0.state"),
				resource.TestCheckResourceAttrSet(datasourceName, "nat_gateways.0.time_created"),
				resource.TestCheckResourceAttrSet(datasourceName, "nat_gateways.0.vcn_id"),
			),
		},
		// verify singular datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_core_nat_gateway", "test_nat_gateway", Required, Create, natGatewaySingularDataSourceRepresentation) +
				compartmentIdVariableStr + NatGatewayResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(singularDatasourceName, "nat_gateway_id"),

				resource.TestCheckResourceAttr(singularDatasourceName, "block_traffic", "true"),
				resource.TestCheckResourceAttr(singularDatasourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(singularDatasourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(singularDatasourceName, "display_name", "displayName2"),
				resource.TestCheckResourceAttr(singularDatasourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "id"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "nat_ip"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "state"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "time_created"),
			),
		},
		// remove singular datasource from previous step so that it doesn't conflict with import tests
		{
			Config: config + compartmentIdVariableStr + NatGatewayResourceConfig,
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

func testAccCheckCoreNatGatewayDestroy(s *terraform.State) error {
	noResourceFound := true
	client := testAccProvider.Meta().(*OracleClients).virtualNetworkClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oci_core_nat_gateway" {
			noResourceFound = false
			request := oci_core.GetNatGatewayRequest{}

			tmp := rs.Primary.ID
			request.NatGatewayId = &tmp

			request.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "core")

			response, err := client.GetNatGateway(context.Background(), request)

			if err == nil {
				deletedLifecycleStates := map[string]bool{
					string(oci_core.NatGatewayLifecycleStateTerminated): true,
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
	if !InSweeperExcludeList("CoreNatGateway") {
		resource.AddTestSweepers("CoreNatGateway", &resource.Sweeper{
			Name:         "CoreNatGateway",
			Dependencies: DependencyGraph["natGateway"],
			F:            sweepCoreNatGatewayResource,
		})
	}
}

func sweepCoreNatGatewayResource(compartment string) error {
	virtualNetworkClient := GetTestClients(&schema.ResourceData{}).virtualNetworkClient()
	natGatewayIds, err := getNatGatewayIds(compartment)
	if err != nil {
		return err
	}
	for _, natGatewayId := range natGatewayIds {
		if ok := SweeperDefaultResourceId[natGatewayId]; !ok {
			deleteNatGatewayRequest := oci_core.DeleteNatGatewayRequest{}

			deleteNatGatewayRequest.NatGatewayId = &natGatewayId

			deleteNatGatewayRequest.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "core")
			_, error := virtualNetworkClient.DeleteNatGateway(context.Background(), deleteNatGatewayRequest)
			if error != nil {
				fmt.Printf("Error deleting NatGateway %s %s, It is possible that the resource is already deleted. Please verify manually \n", natGatewayId, error)
				continue
			}
			WaitTillCondition(testAccProvider, &natGatewayId, natGatewaySweepWaitCondition, time.Duration(3*time.Minute),
				natGatewaySweepResponseFetchOperation, "core", true)
		}
	}
	return nil
}

func getNatGatewayIds(compartment string) ([]string, error) {
	ids := GetResourceIdsToSweep(compartment, "NatGatewayId")
	if ids != nil {
		return ids, nil
	}
	var resourceIds []string
	compartmentId := compartment
	virtualNetworkClient := GetTestClients(&schema.ResourceData{}).virtualNetworkClient()

	listNatGatewaysRequest := oci_core.ListNatGatewaysRequest{}
	listNatGatewaysRequest.CompartmentId = &compartmentId
	listNatGatewaysRequest.LifecycleState = oci_core.NatGatewayLifecycleStateAvailable
	listNatGatewaysResponse, err := virtualNetworkClient.ListNatGateways(context.Background(), listNatGatewaysRequest)

	if err != nil {
		return resourceIds, fmt.Errorf("Error getting NatGateway list for compartment id : %s , %s \n", compartmentId, err)
	}
	for _, natGateway := range listNatGatewaysResponse.Items {
		id := *natGateway.Id
		resourceIds = append(resourceIds, id)
		AddResourceIdToSweeperResourceIdMap(compartmentId, "NatGatewayId", id)
	}
	return resourceIds, nil
}

func natGatewaySweepWaitCondition(response common.OCIOperationResponse) bool {
	// Only stop if the resource is available beyond 3 mins. As there could be an issue for the sweeper to delete the resource and manual intervention required.
	if natGatewayResponse, ok := response.Response.(oci_core.GetNatGatewayResponse); ok {
		return natGatewayResponse.LifecycleState != oci_core.NatGatewayLifecycleStateTerminated
	}
	return false
}

func natGatewaySweepResponseFetchOperation(client *OracleClients, resourceId *string, retryPolicy *common.RetryPolicy) error {
	_, err := client.virtualNetworkClient().GetNatGateway(context.Background(), oci_core.GetNatGatewayRequest{
		NatGatewayId: resourceId,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: retryPolicy,
		},
	})
	return err
}
