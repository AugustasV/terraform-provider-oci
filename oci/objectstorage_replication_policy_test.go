// Copyright (c) 2017, 2021, Oracle and/or its affiliates. All rights reserved.
// Licensed under the Mozilla Public License v2.0

package oci

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/oracle/oci-go-sdk/v49/common"
	oci_object_storage "github.com/oracle/oci-go-sdk/v49/objectstorage"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	replicationBucketName           = testBucketName + "_replication"
	ReplicationPolicyResourceConfig = ReplicationPolicyResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_objectstorage_replication_policy", "test_replication_policy", Optional, Update, replicationPolicyRepresentation)

	replicationPolicySingularDataSourceRepresentation = map[string]interface{}{
		"bucket":         Representation{RepType: Required, Create: `${oci_objectstorage_bucket.test_bucket.name}`},
		"namespace":      Representation{RepType: Required, Create: `${oci_objectstorage_bucket.test_bucket.namespace}`},
		"replication_id": Representation{RepType: Required, Create: `${data.oci_objectstorage_replication_policies.test_replication_policies.replication_policies.0.id}`},
	}

	replicationPolicyDataSourceRepresentation = map[string]interface{}{
		"bucket":    Representation{RepType: Required, Create: `${oci_objectstorage_bucket.test_bucket.name}`},
		"namespace": Representation{RepType: Required, Create: `${oci_objectstorage_bucket.test_bucket.namespace}`},
	}

	replicationPolicyRepresentation = map[string]interface{}{
		"bucket":                  Representation{RepType: Required, Create: `${oci_objectstorage_bucket.test_bucket.name}`},
		"destination_bucket_name": Representation{RepType: Required, Create: `${oci_objectstorage_bucket.test_bucket_replication.name}`},
		"destination_region_name": Representation{RepType: Required, Create: `${var.region}`},
		"name":                    Representation{RepType: Required, Create: `mypolicy`},
		"namespace":               Representation{RepType: Required, Create: `${oci_objectstorage_bucket.test_bucket.namespace}`},
	}

	ReplicationPolicyResourceDependencies = GenerateDataSourceFromRepresentationMap("oci_identity_regions", "test_regions", Required, Create, regionDataSourceRepresentation) +
		GenerateResourceFromRepresentationMap("oci_objectstorage_bucket", "test_bucket", Required, Create, bucketRepresentation) +
		GenerateResourceFromRepresentationMap("oci_objectstorage_bucket", "test_bucket_replication", Required, Create,
			RepresentationCopyWithNewProperties(bucketRepresentation, map[string]interface{}{
				"name": Representation{RepType: Required, Create: replicationBucketName},
			})) + GenerateDataSourceFromRepresentationMap("oci_objectstorage_namespace", "test_namespace", Required, Create, namespaceSingularDataSourceRepresentation)
)

// issue-routing-tag: object_storage/default
func TestObjectStorageReplicationPolicyResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestObjectStorageReplicationPolicyResource_basic")
	defer httpreplay.SaveScenario()

	config := ProviderTestConfig()

	compartmentId := GetEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	resourceName := "oci_objectstorage_replication_policy.test_replication_policy"
	datasourceName := "data.oci_objectstorage_replication_policies.test_replication_policies"
	singularDatasourceName := "data.oci_objectstorage_replication_policy.test_replication_policy"

	var resId string
	// Save TF content to Create resource with only required properties. This has to be exactly the same as the config part in the Create step in the test.
	SaveConfigContent(config+compartmentIdVariableStr+ReplicationPolicyResourceDependencies+
		GenerateResourceFromRepresentationMap("oci_objectstorage_replication_policy", "test_replication_policy", Required, Create, replicationPolicyRepresentation), "objectstorage", "replicationPolicy", t)

	ResourceTest(t, testAccCheckObjectStorageReplicationPolicyDestroy, []resource.TestStep{
		// verify Create
		{
			Config: config + compartmentIdVariableStr + ReplicationPolicyResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_objectstorage_replication_policy", "test_replication_policy", Required, Create, replicationPolicyRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "bucket", testBucketName),
				resource.TestCheckResourceAttrSet(resourceName, "destination_bucket_name"),
				resource.TestCheckResourceAttrSet(resourceName, "destination_region_name"),
				resource.TestCheckResourceAttr(resourceName, "name", "mypolicy"),
				resource.TestCheckResourceAttrSet(resourceName, "namespace"),

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

		// verify datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_objectstorage_replication_policies", "test_replication_policies", Optional, Update, replicationPolicyDataSourceRepresentation) +
				compartmentIdVariableStr + ReplicationPolicyResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_objectstorage_replication_policy", "test_replication_policy", Optional, Update, replicationPolicyRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(datasourceName, "bucket", testBucketName),

				resource.TestCheckResourceAttr(datasourceName, "replication_policies.#", "1"),
				resource.TestCheckResourceAttrSet(datasourceName, "replication_policies.0.destination_bucket_name"),
				resource.TestCheckResourceAttrSet(datasourceName, "replication_policies.0.destination_region_name"),
				resource.TestCheckResourceAttrSet(datasourceName, "replication_policies.0.id"),
				resource.TestCheckResourceAttr(datasourceName, "replication_policies.0.name", "mypolicy"),
				resource.TestCheckResourceAttrSet(datasourceName, "replication_policies.0.status"),
				resource.TestCheckResourceAttrSet(datasourceName, "replication_policies.0.status_message"),
				resource.TestCheckResourceAttrSet(datasourceName, "replication_policies.0.time_created"),
			),
		},

		// verify singular datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_objectstorage_replication_policies", "test_replication_policies", Optional, Update, replicationPolicyDataSourceRepresentation) +
				GenerateDataSourceFromRepresentationMap("oci_objectstorage_replication_policy", "test_replication_policy", Required, Create, replicationPolicySingularDataSourceRepresentation) +
				compartmentIdVariableStr + ReplicationPolicyResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(singularDatasourceName, "bucket", testBucketName),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "replication_id"),

				resource.TestCheckResourceAttrSet(singularDatasourceName, "id"),
				resource.TestCheckResourceAttr(singularDatasourceName, "name", "mypolicy"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "status"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "status_message"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "time_created"),
			),
		},
		// remove singular datasource from previous step so that it doesn't conflict with import tests
		{
			Config: config + compartmentIdVariableStr + ReplicationPolicyResourceConfig,
		},
		// verify resource import
		{
			Config:            config,
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				"time_last_sync",
			},
			ResourceName: resourceName,
		},
	})
}

func testAccCheckObjectStorageReplicationPolicyDestroy(s *terraform.State) error {
	noResourceFound := true
	client := TestAccProvider.Meta().(*OracleClients).objectStorageClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oci_objectstorage_replication_policy" {
			noResourceFound = false
			request := oci_object_storage.GetReplicationPolicyRequest{}

			bucket, namespace, replicationId, err := parseReplicationPolicyCompositeId(rs.Primary.ID)
			if err == nil {
				request.BucketName = &bucket
				request.NamespaceName = &namespace
				request.ReplicationId = &replicationId
			} else {
				log.Printf("[WARN] Get() unable to parse current ID: %s", rs.Primary.ID)
			}

			request.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "object_storage")

			_, err = client.GetReplicationPolicy(context.Background(), request)

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
	if !InSweeperExcludeList("ObjectStorageReplicationPolicy") {
		resource.AddTestSweepers("ObjectStorageReplicationPolicy", &resource.Sweeper{
			Name:         "ObjectStorageReplicationPolicy",
			Dependencies: DependencyGraph["replicationPolicy"],
			F:            sweepObjectStorageReplicationPolicyResource,
		})
	}
}

func sweepObjectStorageReplicationPolicyResource(compartment string) error {
	objectStorageClient := GetTestClients(&schema.ResourceData{}).objectStorageClient()
	replicationPolicyIds, err := getReplicationPolicyIds(compartment)
	if err != nil {
		return err
	}
	for _, replicationPolicyId := range replicationPolicyIds {
		if ok := SweeperDefaultResourceId[replicationPolicyId]; !ok {
			deleteReplicationPolicyRequest := oci_object_storage.DeleteReplicationPolicyRequest{}

			deleteReplicationPolicyRequest.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "object_storage")
			_, error := objectStorageClient.DeleteReplicationPolicy(context.Background(), deleteReplicationPolicyRequest)
			if error != nil {
				fmt.Printf("Error deleting ReplicationPolicy %s %s, It is possible that the resource is already deleted. Please verify manually \n", replicationPolicyId, error)
				continue
			}
		}
	}
	return nil
}

func getReplicationPolicyIds(compartment string) ([]string, error) {
	ids := GetResourceIdsToSweep(compartment, "ReplicationPolicyId")
	if ids != nil {
		return ids, nil
	}
	var resourceIds []string
	compartmentId := compartment
	objectStorageClient := GetTestClients(&schema.ResourceData{}).objectStorageClient()

	listReplicationPoliciesRequest := oci_object_storage.ListReplicationPoliciesRequest{}

	buckets, error := getBucketIds(compartment)
	if error != nil {
		return resourceIds, fmt.Errorf("Error getting bucket required for ReplicationPolicy resource requests \n")
	}
	for _, bucket := range buckets {
		listReplicationPoliciesRequest.BucketName = &bucket

		namespaces, error := getNamespaces(compartment)
		if error != nil {
			return resourceIds, fmt.Errorf("Error getting namespace required for ReplicationPolicy resource requests \n")
		}
		for _, namespace := range namespaces {
			listReplicationPoliciesRequest.NamespaceName = &namespace

			listReplicationPoliciesResponse, err := objectStorageClient.ListReplicationPolicies(context.Background(), listReplicationPoliciesRequest)

			if err != nil {
				return resourceIds, fmt.Errorf("Error getting ReplicationPolicy list for compartment id : %s , %s \n", compartmentId, err)
			}
			for _, replicationPolicy := range listReplicationPoliciesResponse.Items {
				id := *replicationPolicy.Id
				resourceIds = append(resourceIds, id)
				AddResourceIdToSweeperResourceIdMap(compartmentId, "ReplicationPolicyId", id)
			}

		}
	}
	return resourceIds, nil
}
