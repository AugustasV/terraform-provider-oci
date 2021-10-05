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
	oci_core "github.com/oracle/oci-go-sdk/v49/core"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	InstanceConsoleConnectionRequiredOnlyResource = InstanceConsoleConnectionResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_core_instance_console_connection", "test_instance_console_connection", Required, Create, instanceConsoleConnectionRepresentation)

	instanceConsoleConnectionDataSourceRepresentation = map[string]interface{}{
		"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id}`},
		"instance_id":    Representation{RepType: Optional, Create: `${oci_core_instance.test_instance.id}`},
		"filter":         RepresentationGroup{Required, instanceConsoleConnectionDataSourceFilterRepresentation}}
	instanceConsoleConnectionDataSourceFilterRepresentation = map[string]interface{}{
		"name":   Representation{RepType: Required, Create: `id`},
		"values": Representation{RepType: Required, Create: []string{`${oci_core_instance_console_connection.test_instance_console_connection.id}`}},
	}

	instanceConsoleConnectionRepresentation = map[string]interface{}{
		"instance_id":   Representation{RepType: Required, Create: `${oci_core_instance.test_instance.id}`},
		"public_key":    Representation{RepType: Required, Create: `ssh-rsa KKKLK3NzaC1yc2EAAAADAQABAAABAQC+UC9MFNA55NIVtKPIBCNw7++ACXhD0hx+Zyj25JfHykjz/QU3Q5FAU3DxDbVXyubgXfb/GJnrKRY8O4QDdvnZZRvQFFEOaApThAmCAM5MuFUIHdFvlqP+0W+ZQnmtDhwVe2NCfcmOrMuaPEgOKO3DOW6I/qOOdO691Xe2S9NgT9HhN0ZfFtEODVgvYulgXuCCXsJs+NUqcHAOxxFUmwkbPvYi0P0e2DT8JKeiOOC8VKUEgvVx+GKmqasm+Y6zHFW7vv3g2GstE1aRs3mttHRoC/JPM86PRyIxeWXEMzyG5wHqUu4XZpDbnWNxi6ugxnAGiL3CrIFdCgRNgHz5qS1l MustWin`},
		"defined_tags":  Representation{RepType: Optional, Create: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "value")}`, Update: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "updatedValue")}`},
		"freeform_tags": Representation{RepType: Optional, Create: map[string]string{"Department": "Finance"}, Update: map[string]string{"Department": "Accounting"}},
	}

	InstanceConsoleConnectionResourceDependencies = GenerateResourceFromRepresentationMap("oci_core_subnet", "test_subnet", Required, Create, subnetRepresentation) +
		GenerateResourceFromRepresentationMap("oci_core_vcn", "test_vcn", Required, Create, vcnRepresentation) +
		OciImageIdsVariable +
		GenerateResourceFromRepresentationMap("oci_core_instance", "test_instance", Required, Create, instanceRepresentation) +
		AvailabilityDomainConfig +
		DefinedTagsDependencies
)

// issue-routing-tag: core/computeServices
func TestCoreInstanceConsoleConnectionResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestCoreInstanceConsoleConnectionResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	resourceName := "oci_core_instance_console_connection.test_instance_console_connection"
	datasourceName := "data.oci_core_instance_console_connections.test_instance_console_connections"

	var resId, resId2 string
	// Save TF content to Create resource with optional properties. This has to be exactly the same as the config part in the "Create with optionals" step in the test.
	SaveConfigContent(config+compartmentIdVariableStr+InstanceConsoleConnectionResourceDependencies+
		GenerateResourceFromRepresentationMap("oci_core_instance_console_connection", "test_instance_console_connection", Optional, Create, instanceConsoleConnectionRepresentation), "core", "instanceConsoleConnection", t)

	ResourceTest(t, testAccCheckCoreInstanceConsoleConnectionDestroy, []resource.TestStep{
		// verify Create
		{
			Config: config + compartmentIdVariableStr + InstanceConsoleConnectionResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_core_instance_console_connection", "test_instance_console_connection", Required, Create, instanceConsoleConnectionRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
				resource.TestCheckResourceAttr(resourceName, "public_key", "ssh-rsa KKKLK3NzaC1yc2EAAAADAQABAAABAQC+UC9MFNA55NIVtKPIBCNw7++ACXhD0hx+Zyj25JfHykjz/QU3Q5FAU3DxDbVXyubgXfb/GJnrKRY8O4QDdvnZZRvQFFEOaApThAmCAM5MuFUIHdFvlqP+0W+ZQnmtDhwVe2NCfcmOrMuaPEgOKO3DOW6I/qOOdO691Xe2S9NgT9HhN0ZfFtEODVgvYulgXuCCXsJs+NUqcHAOxxFUmwkbPvYi0P0e2DT8JKeiOOC8VKUEgvVx+GKmqasm+Y6zHFW7vv3g2GstE1aRs3mttHRoC/JPM86PRyIxeWXEMzyG5wHqUu4XZpDbnWNxi6ugxnAGiL3CrIFdCgRNgHz5qS1l MustWin"),

				func(s *terraform.State) (err error) {
					resId, err = FromInstanceState(s, resourceName, "id")
					return err
				},
			),
		},

		// delete before next Create
		{
			Config: config + compartmentIdVariableStr + InstanceConsoleConnectionResourceDependencies,
		},
		// verify Create with optionals
		{
			Config: config + compartmentIdVariableStr + InstanceConsoleConnectionResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_core_instance_console_connection", "test_instance_console_connection", Optional, Create, instanceConsoleConnectionRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
				resource.TestCheckResourceAttr(resourceName, "public_key", "ssh-rsa KKKLK3NzaC1yc2EAAAADAQABAAABAQC+UC9MFNA55NIVtKPIBCNw7++ACXhD0hx+Zyj25JfHykjz/QU3Q5FAU3DxDbVXyubgXfb/GJnrKRY8O4QDdvnZZRvQFFEOaApThAmCAM5MuFUIHdFvlqP+0W+ZQnmtDhwVe2NCfcmOrMuaPEgOKO3DOW6I/qOOdO691Xe2S9NgT9HhN0ZfFtEODVgvYulgXuCCXsJs+NUqcHAOxxFUmwkbPvYi0P0e2DT8JKeiOOC8VKUEgvVx+GKmqasm+Y6zHFW7vv3g2GstE1aRs3mttHRoC/JPM86PRyIxeWXEMzyG5wHqUu4XZpDbnWNxi6ugxnAGiL3CrIFdCgRNgHz5qS1l MustWin"),

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
			Config: config + compartmentIdVariableStr + InstanceConsoleConnectionResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_core_instance_console_connection", "test_instance_console_connection", Optional, Update, instanceConsoleConnectionRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "instance_id"),
				resource.TestCheckResourceAttr(resourceName, "public_key", "ssh-rsa KKKLK3NzaC1yc2EAAAADAQABAAABAQC+UC9MFNA55NIVtKPIBCNw7++ACXhD0hx+Zyj25JfHykjz/QU3Q5FAU3DxDbVXyubgXfb/GJnrKRY8O4QDdvnZZRvQFFEOaApThAmCAM5MuFUIHdFvlqP+0W+ZQnmtDhwVe2NCfcmOrMuaPEgOKO3DOW6I/qOOdO691Xe2S9NgT9HhN0ZfFtEODVgvYulgXuCCXsJs+NUqcHAOxxFUmwkbPvYi0P0e2DT8JKeiOOC8VKUEgvVx+GKmqasm+Y6zHFW7vv3g2GstE1aRs3mttHRoC/JPM86PRyIxeWXEMzyG5wHqUu4XZpDbnWNxi6ugxnAGiL3CrIFdCgRNgHz5qS1l MustWin"),

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
				GenerateDataSourceFromRepresentationMap("oci_core_instance_console_connections", "test_instance_console_connections", Optional, Update, instanceConsoleConnectionDataSourceRepresentation) +
				compartmentIdVariableStr + InstanceConsoleConnectionResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_core_instance_console_connection", "test_instance_console_connection", Optional, Update, instanceConsoleConnectionRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttrSet(datasourceName, "instance_id"),

				resource.TestCheckResourceAttr(datasourceName, "instance_console_connections.#", "1"),
				resource.TestCheckResourceAttrSet(datasourceName, "instance_console_connections.0.compartment_id"),
				resource.TestCheckResourceAttrSet(datasourceName, "instance_console_connections.0.connection_string"),
				resource.TestCheckResourceAttr(datasourceName, "instance_console_connections.0.defined_tags.%", "1"),
				resource.TestCheckResourceAttrSet(datasourceName, "instance_console_connections.0.fingerprint"),
				resource.TestCheckResourceAttr(datasourceName, "instance_console_connections.0.freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(datasourceName, "instance_console_connections.0.id"),
				resource.TestCheckResourceAttrSet(datasourceName, "instance_console_connections.0.instance_id"),
				resource.TestCheckResourceAttrSet(datasourceName, "instance_console_connections.0.service_host_key_fingerprint"),
				resource.TestCheckResourceAttrSet(datasourceName, "instance_console_connections.0.state"),
				resource.TestCheckResourceAttrSet(datasourceName, "instance_console_connections.0.vnc_connection_string"),
			),
		},
		// verify resource import
		{
			Config:            config,
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				"public_key",
			},
			ResourceName: resourceName,
		},
	})
}

func testAccCheckCoreInstanceConsoleConnectionDestroy(s *terraform.State) error {
	noResourceFound := true
	client := testAccProvider.Meta().(*OracleClients).computeClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oci_core_instance_console_connection" {
			noResourceFound = false
			request := oci_core.GetInstanceConsoleConnectionRequest{}

			tmp := rs.Primary.ID
			request.InstanceConsoleConnectionId = &tmp

			request.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "core")

			response, err := client.GetInstanceConsoleConnection(context.Background(), request)

			if err == nil {
				deletedLifecycleStates := map[string]bool{
					string(oci_core.InstanceConsoleConnectionLifecycleStateDeleted): true,
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
	if !InSweeperExcludeList("CoreInstanceConsoleConnection") {
		resource.AddTestSweepers("CoreInstanceConsoleConnection", &resource.Sweeper{
			Name:         "CoreInstanceConsoleConnection",
			Dependencies: DependencyGraph["instanceConsoleConnection"],
			F:            sweepCoreInstanceConsoleConnectionResource,
		})
	}
}

func sweepCoreInstanceConsoleConnectionResource(compartment string) error {
	computeClient := GetTestClients(&schema.ResourceData{}).computeClient()
	instanceConsoleConnectionIds, err := getInstanceConsoleConnectionIds(compartment)
	if err != nil {
		return err
	}
	for _, instanceConsoleConnectionId := range instanceConsoleConnectionIds {
		if ok := SweeperDefaultResourceId[instanceConsoleConnectionId]; !ok {
			deleteInstanceConsoleConnectionRequest := oci_core.DeleteInstanceConsoleConnectionRequest{}

			deleteInstanceConsoleConnectionRequest.InstanceConsoleConnectionId = &instanceConsoleConnectionId

			deleteInstanceConsoleConnectionRequest.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "core")
			_, error := computeClient.DeleteInstanceConsoleConnection(context.Background(), deleteInstanceConsoleConnectionRequest)
			if error != nil {
				fmt.Printf("Error deleting InstanceConsoleConnection %s %s, It is possible that the resource is already deleted. Please verify manually \n", instanceConsoleConnectionId, error)
				continue
			}
			WaitTillCondition(testAccProvider, &instanceConsoleConnectionId, instanceConsoleConnectionSweepWaitCondition, time.Duration(3*time.Minute),
				instanceConsoleConnectionSweepResponseFetchOperation, "core", true)
		}
	}
	return nil
}

func getInstanceConsoleConnectionIds(compartment string) ([]string, error) {
	ids := GetResourceIdsToSweep(compartment, "InstanceConsoleConnectionId")
	if ids != nil {
		return ids, nil
	}
	var resourceIds []string
	compartmentId := compartment
	computeClient := GetTestClients(&schema.ResourceData{}).computeClient()

	listInstanceConsoleConnectionsRequest := oci_core.ListInstanceConsoleConnectionsRequest{}
	listInstanceConsoleConnectionsRequest.CompartmentId = &compartmentId
	listInstanceConsoleConnectionsResponse, err := computeClient.ListInstanceConsoleConnections(context.Background(), listInstanceConsoleConnectionsRequest)

	if err != nil {
		return resourceIds, fmt.Errorf("Error getting InstanceConsoleConnection list for compartment id : %s , %s \n", compartmentId, err)
	}
	for _, instanceConsoleConnection := range listInstanceConsoleConnectionsResponse.Items {
		if instanceConsoleConnection.LifecycleState == oci_core.InstanceConsoleConnectionLifecycleStateDeleted {
			continue
		}
		id := *instanceConsoleConnection.Id
		resourceIds = append(resourceIds, id)
		AddResourceIdToSweeperResourceIdMap(compartmentId, "InstanceConsoleConnectionId", id)
	}
	return resourceIds, nil
}

func instanceConsoleConnectionSweepWaitCondition(response common.OCIOperationResponse) bool {
	// Only stop if the resource is available beyond 3 mins. As there could be an issue for the sweeper to delete the resource and manual intervention required.
	if instanceConsoleConnectionResponse, ok := response.Response.(oci_core.GetInstanceConsoleConnectionResponse); ok {
		return instanceConsoleConnectionResponse.LifecycleState != oci_core.InstanceConsoleConnectionLifecycleStateDeleted
	}
	return false
}

func instanceConsoleConnectionSweepResponseFetchOperation(client *OracleClients, resourceId *string, retryPolicy *common.RetryPolicy) error {
	_, err := client.computeClient().GetInstanceConsoleConnection(context.Background(), oci_core.GetInstanceConsoleConnectionRequest{
		InstanceConsoleConnectionId: resourceId,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: retryPolicy,
		},
	})
	return err
}
