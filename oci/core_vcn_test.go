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
	"github.com/oracle/oci-go-sdk/v53/common"
	oci_core "github.com/oracle/oci-go-sdk/v53/core"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	VcnRequiredOnlyResource = VcnRequiredOnlyResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_core_vcn", "test_vcn", Required, Create, vcnRepresentation)

	VcnResourceConfig = GenerateResourceFromRepresentationMap("oci_core_vcn", "test_vcn", Optional, Create, vcnRepresentation)

	vcnSingularDataSourceRepresentation = map[string]interface{}{
		"vcn_id": Representation{RepType: Required, Create: `${oci_core_vcn.test_vcn.id}`},
	}

	vcnDataSourceRepresentation = map[string]interface{}{
		"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id}`},
		"display_name":   Representation{RepType: Optional, Create: `displayName`, Update: `displayName2`},
		"state":          Representation{RepType: Optional, Create: `AVAILABLE`},
		"filter":         RepresentationGroup{Required, vcnDataSourceFilterRepresentation}}
	vcnDataSourceFilterRepresentation = map[string]interface{}{
		"name":   Representation{RepType: Required, Create: `id`},
		"values": Representation{RepType: Required, Create: []string{`${oci_core_vcn.test_vcn.id}`}},
	}

	vcnRepresentation = map[string]interface{}{
		"cidr_block":     Representation{RepType: Required, Create: `10.0.0.0/16`},
		"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id}`},
		"defined_tags":   Representation{RepType: Optional, Create: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "value")}`, Update: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "updatedValue")}`},
		"display_name":   Representation{RepType: Optional, Create: `displayName`, Update: `displayName2`},
		"dns_label":      Representation{RepType: Optional, Create: `dnslabel`},
		"freeform_tags":  Representation{RepType: Optional, Create: map[string]string{"Department": "Finance"}, Update: map[string]string{"Department": "Accounting"}},
		"lifecycle":      RepresentationGroup{Required, ignoreDefinedTagsChangesRep},
	}

	VcnRequiredOnlyResourceDependencies = ``
	VcnResourceDependencies             = DefinedTagsDependencies
)

// issue-routing-tag: core/virtualNetwork
func TestCoreVcnResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestCoreVcnResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	compartmentIdU := getEnvSettingWithDefault("compartment_id_for_update", compartmentId)
	compartmentIdUVariableStr := fmt.Sprintf("variable \"compartment_id_for_update\" { default = \"%s\" }\n", compartmentIdU)

	resourceName := "oci_core_vcn.test_vcn"
	datasourceName := "data.oci_core_vcns.test_vcns"
	singularDatasourceName := "data.oci_core_vcn.test_vcn"

	var resId, resId2 string
	// Save TF content to Create resource with optional properties. This has to be exactly the same as the config part in the "Create with optionals" step in the test.
	SaveConfigContent(config+compartmentIdVariableStr+VcnResourceDependencies+
		GenerateResourceFromRepresentationMap("oci_core_vcn", "test_vcn", Optional, Create, vcnRepresentation), "core", "vcn", t)

	ResourceTest(t, testAccCheckCoreVcnDestroy, []resource.TestStep{
		// verify Create
		{
			Config: config + compartmentIdVariableStr + VcnResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_core_vcn", "test_vcn", Required, Create, vcnRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "cidr_block", "10.0.0.0/16"),
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),

				func(s *terraform.State) (err error) {
					resId, err = FromInstanceState(s, resourceName, "id")
					return err
				},
			),
		},

		// delete before next Create
		{
			Config: config + compartmentIdVariableStr + VcnResourceDependencies,
		},
		// verify Create with optionals
		{
			Config: config + compartmentIdVariableStr + VcnResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_core_vcn", "test_vcn", Optional, Create,
					RepresentationCopyWithNewProperties(RepresentationCopyWithRemovedProperties(vcnRepresentation, []string{"cidr_blocks"}), map[string]interface{}{
						"is_ipv6enabled": Representation{RepType: Optional, Create: `true`},
					})),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "cidr_block", "10.0.0.0/16"),
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "display_name", "displayName"),
				resource.TestCheckResourceAttr(resourceName, "dns_label", "dnslabel"),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "ipv6cidr_blocks.#", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttrSet(resourceName, "state"),

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
			Config: config + compartmentIdVariableStr + compartmentIdUVariableStr + VcnResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_core_vcn", "test_vcn", Optional, Create,
					RepresentationCopyWithNewProperties(vcnRepresentation, map[string]interface{}{
						"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id_for_update}`},
					})),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "cidr_block", "10.0.0.0/16"),
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentIdU),
				resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "display_name", "displayName"),
				resource.TestCheckResourceAttr(resourceName, "dns_label", "dnslabel"),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "ipv6cidr_blocks.#", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttrSet(resourceName, "state"),

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
			Config: config + compartmentIdVariableStr + VcnResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_core_vcn", "test_vcn", Optional, Update, RepresentationCopyWithNewProperties(vcnRepresentation, map[string]interface{}{
					"is_ipv6enabled": Representation{RepType: Optional, Update: `true`},
				})),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "cidr_block", "10.0.0.0/16"),
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "display_name", "displayName2"),
				resource.TestCheckResourceAttr(resourceName, "dns_label", "dnslabel"),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "ipv6cidr_blocks.#", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttrSet(resourceName, "state"),

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
				GenerateDataSourceFromRepresentationMap("oci_core_vcns", "test_vcns", Optional, Update, vcnDataSourceRepresentation) +
				compartmentIdVariableStr + VcnResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_core_vcn", "test_vcn", Optional, Update, RepresentationCopyWithNewProperties(vcnRepresentation, map[string]interface{}{
					"is_ipv6enabled": Representation{RepType: Optional, Update: `true`},
				})),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(datasourceName, "display_name", "displayName2"),
				resource.TestCheckResourceAttr(datasourceName, "state", "AVAILABLE"),

				resource.TestCheckResourceAttr(datasourceName, "virtual_networks.#", "1"),
				resource.TestCheckResourceAttr(datasourceName, "virtual_networks.0.cidr_block", "10.0.0.0/16"),
				resource.TestCheckResourceAttr(datasourceName, "virtual_networks.0.compartment_id", compartmentId),
				resource.TestCheckResourceAttrSet(datasourceName, "virtual_networks.0.default_dhcp_options_id"),
				resource.TestCheckResourceAttrSet(datasourceName, "virtual_networks.0.default_route_table_id"),
				resource.TestCheckResourceAttrSet(datasourceName, "virtual_networks.0.default_security_list_id"),
				resource.TestCheckResourceAttr(datasourceName, "virtual_networks.0.defined_tags.%", "1"),
				resource.TestCheckResourceAttr(datasourceName, "virtual_networks.0.display_name", "displayName2"),
				resource.TestCheckResourceAttr(datasourceName, "virtual_networks.0.dns_label", "dnslabel"),
				resource.TestCheckResourceAttr(datasourceName, "virtual_networks.0.freeform_tags.%", "1"),
				resource.TestCheckResourceAttr(datasourceName, "virtual_networks.0.ipv6cidr_blocks.#", "1"),
				resource.TestCheckResourceAttrSet(datasourceName, "virtual_networks.0.id"),
				resource.TestCheckResourceAttrSet(datasourceName, "virtual_networks.0.state"),
				resource.TestCheckResourceAttrSet(datasourceName, "virtual_networks.0.time_created"),
				resource.TestCheckResourceAttrSet(datasourceName, "virtual_networks.0.vcn_domain_name"),
			),
		},
		// verify singular datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_core_vcn", "test_vcn", Required, Create, vcnSingularDataSourceRepresentation) +
				compartmentIdVariableStr + VcnResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_core_vcn", "test_vcn", Optional, Update, RepresentationCopyWithNewProperties(vcnRepresentation, map[string]interface{}{
					"is_ipv6enabled": Representation{RepType: Optional, Create: `true`},
				})),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(singularDatasourceName, "vcn_id"),

				resource.TestCheckResourceAttr(singularDatasourceName, "cidr_block", "10.0.0.0/16"),
				resource.TestCheckResourceAttr(singularDatasourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "default_dhcp_options_id"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "default_route_table_id"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "default_security_list_id"),
				resource.TestCheckResourceAttr(singularDatasourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(singularDatasourceName, "display_name", "displayName2"),
				resource.TestCheckResourceAttr(singularDatasourceName, "dns_label", "dnslabel"),
				resource.TestCheckResourceAttr(singularDatasourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttr(singularDatasourceName, "ipv6cidr_blocks.#", "1"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "id"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "state"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "time_created"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "vcn_domain_name"),
			),
		},
		// remove singular datasource from previous step so that it doesn't conflict with import tests
		{
			Config: config + compartmentIdVariableStr + VcnResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_core_vcn", "test_vcn", Optional, Update, vcnRepresentation),
		},
		// verify resource import
		{
			Config:            config,
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				"is_ipv6enabled",
			},
			ResourceName: resourceName,
		},
	})
}

func testAccCheckCoreVcnDestroy(s *terraform.State) error {
	noResourceFound := true
	client := testAccProvider.Meta().(*OracleClients).virtualNetworkClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oci_core_vcn" {
			noResourceFound = false
			request := oci_core.GetVcnRequest{}

			tmp := rs.Primary.ID
			request.VcnId = &tmp

			request.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "core")

			response, err := client.GetVcn(context.Background(), request)

			if err == nil {
				deletedLifecycleStates := map[string]bool{
					string(oci_core.VcnLifecycleStateTerminated): true,
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
	if !InSweeperExcludeList("CoreVcn") {
		resource.AddTestSweepers("CoreVcn", &resource.Sweeper{
			Name:         "CoreVcn",
			Dependencies: DependencyGraph["vcn"],
			F:            sweepCoreVcnResource,
		})
	}
}

func sweepCoreVcnResource(compartment string) error {
	virtualNetworkClient := GetTestClients(&schema.ResourceData{}).virtualNetworkClient()
	vcnIds, err := getVcnIds(compartment)
	if err != nil {
		return err
	}
	for _, vcnId := range vcnIds {
		if ok := SweeperDefaultResourceId[vcnId]; !ok {
			deleteVcnRequest := oci_core.DeleteVcnRequest{}

			deleteVcnRequest.VcnId = &vcnId

			deleteVcnRequest.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "core")
			_, error := virtualNetworkClient.DeleteVcn(context.Background(), deleteVcnRequest)
			if error != nil {
				fmt.Printf("Error deleting Vcn %s %s, It is possible that the resource is already deleted. Please verify manually \n", vcnId, error)
				continue
			}
			WaitTillCondition(testAccProvider, &vcnId, vcnSweepWaitCondition, time.Duration(3*time.Minute),
				vcnSweepResponseFetchOperation, "core", true)
		}
	}
	return nil
}

func getVcnIds(compartment string) ([]string, error) {
	ids := GetResourceIdsToSweep(compartment, "VcnId")
	if ids != nil {
		return ids, nil
	}
	var resourceIds []string
	compartmentId := compartment
	virtualNetworkClient := GetTestClients(&schema.ResourceData{}).virtualNetworkClient()

	listVcnsRequest := oci_core.ListVcnsRequest{}
	listVcnsRequest.CompartmentId = &compartmentId
	listVcnsRequest.LifecycleState = oci_core.VcnLifecycleStateAvailable
	listVcnsResponse, err := virtualNetworkClient.ListVcns(context.Background(), listVcnsRequest)

	if err != nil {
		return resourceIds, fmt.Errorf("Error getting Vcn list for compartment id : %s , %s \n", compartmentId, err)
	}
	for _, vcn := range listVcnsResponse.Items {
		id := *vcn.Id
		resourceIds = append(resourceIds, id)
		AddResourceIdToSweeperResourceIdMap(compartmentId, "VcnId", id)
		SweeperDefaultResourceId[*vcn.DefaultDhcpOptionsId] = true
		SweeperDefaultResourceId[*vcn.DefaultRouteTableId] = true
		SweeperDefaultResourceId[*vcn.DefaultSecurityListId] = true

	}
	return resourceIds, nil
}

func vcnSweepWaitCondition(response common.OCIOperationResponse) bool {
	// Only stop if the resource is available beyond 3 mins. As there could be an issue for the sweeper to delete the resource and manual intervention required.
	if vcnResponse, ok := response.Response.(oci_core.GetVcnResponse); ok {
		return vcnResponse.LifecycleState != oci_core.VcnLifecycleStateTerminated
	}
	return false
}

func vcnSweepResponseFetchOperation(client *OracleClients, resourceId *string, retryPolicy *common.RetryPolicy) error {
	_, err := client.virtualNetworkClient().GetVcn(context.Background(), oci_core.GetVcnRequest{
		VcnId: resourceId,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: retryPolicy,
		},
	})
	return err
}
