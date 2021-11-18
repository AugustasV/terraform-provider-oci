// Copyright (c) 2017, 2021, Oracle and/or its affiliates. All rights reserved.
// Licensed under the Mozilla Public License v2.0

package oci

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	oci_apigateway "github.com/oracle/oci-go-sdk/v52/apigateway"
	"github.com/oracle/oci-go-sdk/v52/common"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	ApiGatewayCertificateRequiredOnlyResource = ApiGatewayCertificateResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_apigateway_certificate", "test_certificate", Required, Create, apiGatewaycertificateRepresentation)

	CertificateResourceConfig = ApiGatewayCertificateResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_apigateway_certificate", "test_certificate", Optional, Update, apiGatewaycertificateRepresentation)

	apiGatewaycertificateSingularDataSourceRepresentation = map[string]interface{}{
		"certificate_id": Representation{RepType: Required, Create: `${oci_apigateway_certificate.test_certificate.id}`},
	}

	apiGatewaycertificateDataSourceRepresentation = map[string]interface{}{
		"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id}`},
		"display_name":   Representation{RepType: Optional, Create: `displayName`, Update: `displayName2`},
		"state":          Representation{RepType: Optional, Create: `ACTIVE`},
		"filter":         RepresentationGroup{Required, apiGatewaycertificateDataSourceFilterRepresentation}}
	apiGatewaycertificateDataSourceFilterRepresentation = map[string]interface{}{
		"name":   Representation{RepType: Required, Create: `id`},
		"values": Representation{RepType: Required, Create: []string{`${oci_apigateway_certificate.test_certificate.id}`}},
	}

	apiGatewaycertificateRepresentation = map[string]interface{}{
		"certificate":               Representation{RepType: Required, Create: "${var.api_certificate_value}"},
		"compartment_id":            Representation{RepType: Required, Create: `${var.compartment_id}`},
		"private_key":               Representation{RepType: Required, Create: "${var.api_private_key_value}"},
		"defined_tags":              Representation{RepType: Optional, Create: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "value")}`, Update: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "updatedValue")}`},
		"display_name":              Representation{RepType: Optional, Create: `displayName`, Update: `displayName2`},
		"freeform_tags":             Representation{RepType: Optional, Create: map[string]string{"Department": "Finance"}, Update: map[string]string{"Department": "Accounting"}},
		"intermediate_certificates": Representation{RepType: Optional, Create: "${var.api_intermediate_certificate_value}"},
	}

	apiCertificate            = getEnvSettingWithBlankDefault("api_certificate")
	apiCertificateVariableStr = fmt.Sprintf("variable \"api_certificate_value\" { default = \"%s\" }\n", apiCertificate)

	apiPrivateKey            = getEnvSettingWithBlankDefault("api_private_key")
	apiPrivateKeyVariableStr = fmt.Sprintf("variable \"api_private_key_value\" { default = \"%s\" }\n", apiPrivateKey)

	apiIntermediateCertificate            = getEnvSettingWithBlankDefault("api_intermediate_certificate")
	apiIntermediateCertificateVariableStr = fmt.Sprintf("variable \"api_intermediate_certificate_value\" { default = \"%s\" }\n", apiIntermediateCertificate)

	ApiGatewayCertificateResourceDependencies = DefinedTagsDependencies + apiCertificateVariableStr + apiPrivateKeyVariableStr + apiIntermediateCertificateVariableStr
)

// issue-routing-tag: apigateway/default
func TestApigatewayCertificateResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestApigatewayCertificateResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	compartmentIdU := getEnvSettingWithDefault("compartment_id_for_update", compartmentId)
	compartmentIdUVariableStr := fmt.Sprintf("variable \"compartment_id_for_update\" { default = \"%s\" }\n", compartmentIdU)

	resourceName := "oci_apigateway_certificate.test_certificate"
	datasourceName := "data.oci_apigateway_certificates.test_certificates"
	singularDatasourceName := "data.oci_apigateway_certificate.test_certificate"

	var resId, resId2 string
	// Save TF content to Create resource with optional properties. This has to be exactly the same as the config part in the "Create with optionals" step in the test.
	SaveConfigContent(config+compartmentIdVariableStr+CertificateResourceDependencies+
		GenerateResourceFromRepresentationMap("oci_apigateway_certificate", "test_certificate", Optional, Create, certificateRepresentation), "apigateway", "certificate", t)

	ResourceTest(t, testAccCheckApigatewayCertificateDestroy, []resource.TestStep{
		// verify Create
		{
			Config: config + compartmentIdVariableStr + ApiGatewayCertificateResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_apigateway_certificate", "test_certificate", Required, Create, apiGatewaycertificateRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestMatchResourceAttr(resourceName, "certificate", regexp.MustCompile("-----BEGIN CERT.*")),
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestMatchResourceAttr(resourceName, "private_key", regexp.MustCompile("-----BEGIN RSA.*")),

				func(s *terraform.State) (err error) {
					resId, err = FromInstanceState(s, resourceName, "id")
					return err
				},
			),
		},

		// delete before next Create
		{
			Config: config + compartmentIdVariableStr + ApiGatewayCertificateResourceDependencies,
		},
		// verify Create with optionals
		{
			Config: config + compartmentIdVariableStr + ApiGatewayCertificateResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_apigateway_certificate", "test_certificate", Optional, Create, apiGatewaycertificateRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(resourceName, "certificate"),
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "display_name", "displayName"),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttrSet(resourceName, "intermediate_certificates"),
				resource.TestCheckResourceAttrSet(resourceName, "private_key"),
				resource.TestCheckResourceAttrSet(resourceName, "subject_names.0"),
				resource.TestCheckResourceAttrSet(resourceName, "time_created"),
				resource.TestCheckResourceAttrSet(resourceName, "time_not_valid_after"),

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
			Config: config + compartmentIdVariableStr + compartmentIdUVariableStr + ApiGatewayCertificateResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_apigateway_certificate", "test_certificate", Optional, Create,
					RepresentationCopyWithNewProperties(apiGatewaycertificateRepresentation, map[string]interface{}{
						"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id_for_update}`},
					})),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(resourceName, "certificate"),
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentIdU),
				resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "display_name", "displayName"),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttrSet(resourceName, "intermediate_certificates"),
				resource.TestCheckResourceAttrSet(resourceName, "private_key"),
				resource.TestCheckResourceAttrSet(resourceName, "subject_names.0"),
				resource.TestCheckResourceAttrSet(resourceName, "time_created"),
				resource.TestCheckResourceAttrSet(resourceName, "time_not_valid_after"),

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
			Config: config + compartmentIdVariableStr + ApiGatewayCertificateResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_apigateway_certificate", "test_certificate", Optional, Update, apiGatewaycertificateRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(resourceName, "certificate"),
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "display_name", "displayName2"),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttrSet(resourceName, "intermediate_certificates"),
				resource.TestCheckResourceAttrSet(resourceName, "private_key"),
				resource.TestCheckResourceAttrSet(resourceName, "subject_names.0"),
				resource.TestCheckResourceAttrSet(resourceName, "time_created"),
				resource.TestCheckResourceAttrSet(resourceName, "time_not_valid_after"),

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
				GenerateDataSourceFromRepresentationMap("oci_apigateway_certificates", "test_certificates", Optional, Update, apiGatewaycertificateDataSourceRepresentation) +
				compartmentIdVariableStr + ApiGatewayCertificateResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_apigateway_certificate", "test_certificate", Optional, Update, apiGatewaycertificateRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(datasourceName, "display_name", "displayName2"),
				resource.TestCheckResourceAttr(datasourceName, "state", "ACTIVE"),

				resource.TestCheckResourceAttr(datasourceName, "certificate_collection.#", "1"),
				resource.TestCheckResourceAttr(datasourceName, "certificate_collection.0.items.#", "1"),
			),
		},
		// verify singular datasource
		{
			Config: config +
				GenerateResourceFromRepresentationMap("oci_apigateway_certificate", "test_certificate", Optional, Update, apiGatewaycertificateRepresentation) +
				GenerateDataSourceFromRepresentationMap("oci_apigateway_certificate", "test_certificate", Required, Create, apiGatewaycertificateSingularDataSourceRepresentation) +
				compartmentIdVariableStr + ApiGatewayCertificateResourceDependencies,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(singularDatasourceName, "certificate_id"),
				resource.TestCheckResourceAttrSet(resourceName, "certificate"),
				resource.TestCheckResourceAttr(singularDatasourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(singularDatasourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(singularDatasourceName, "display_name", "displayName2"),
				resource.TestCheckResourceAttr(singularDatasourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "id"),
				resource.TestCheckResourceAttrSet(resourceName, "intermediate_certificates"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "state"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "subject_names.0"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "time_created"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "time_not_valid_after"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "time_updated"),
			),
		},
		// remove singular datasource from previous step so that it doesn't conflict with import tests
		{
			Config: config + compartmentIdVariableStr + CertificateResourceConfig,
		},
		// verify resource import
		{
			Config:            config,
			ImportState:       true,
			ImportStateVerify: true,
			ImportStateVerifyIgnore: []string{
				"private_key",
			},
			ResourceName: resourceName,
		},
	})
}

func testAccCheckApigatewayCertificateDestroy(s *terraform.State) error {
	noResourceFound := true
	client := testAccProvider.Meta().(*OracleClients).apiGatewayClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oci_apigateway_certificate" {
			noResourceFound = false
			request := oci_apigateway.GetCertificateRequest{}

			tmp := rs.Primary.ID
			request.CertificateId = &tmp

			request.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "apigateway")

			response, err := client.GetCertificate(context.Background(), request)

			if err == nil {
				deletedLifecycleStates := map[string]bool{
					string(oci_apigateway.CertificateLifecycleStateDeleted): true,
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
	if !InSweeperExcludeList("ApigatewayCertificate") {
		resource.AddTestSweepers("ApigatewayCertificate", &resource.Sweeper{
			Name:         "ApigatewayCertificate",
			Dependencies: DependencyGraph["certificate"],
			F:            sweepApigatewayCertificateResource,
		})
	}
}

func sweepApigatewayCertificateResource(compartment string) error {
	apiGatewayClient := GetTestClients(&schema.ResourceData{}).apiGatewayClient()
	certificateIds, err := getApiGatewayCertificateIds(compartment)
	if err != nil {
		return err
	}
	for _, certificateId := range certificateIds {
		if ok := SweeperDefaultResourceId[certificateId]; !ok {
			deleteCertificateRequest := oci_apigateway.DeleteCertificateRequest{}

			deleteCertificateRequest.CertificateId = &certificateId

			deleteCertificateRequest.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "apigateway")
			_, error := apiGatewayClient.DeleteCertificate(context.Background(), deleteCertificateRequest)
			if error != nil {
				fmt.Printf("Error deleting Certificate %s %s, It is possible that the resource is already deleted. Please verify manually \n", certificateId, error)
				continue
			}
			WaitTillCondition(testAccProvider, &certificateId, apiGatewayCertificateSweepWaitCondition, time.Duration(3*time.Minute),
				apiGatewayCertificateSweepResponseFetchOperation, "apigateway", true)
		}
	}
	return nil
}

func getApiGatewayCertificateIds(compartment string) ([]string, error) {
	ids := GetResourceIdsToSweep(compartment, "CertificateId")
	if ids != nil {
		return ids, nil
	}
	var resourceIds []string
	compartmentId := compartment
	apiGatewayClient := GetTestClients(&schema.ResourceData{}).apiGatewayClient()

	listCertificatesRequest := oci_apigateway.ListCertificatesRequest{}
	listCertificatesRequest.CompartmentId = &compartmentId
	listCertificatesRequest.LifecycleState = oci_apigateway.CertificateLifecycleStateActive
	listCertificatesResponse, err := apiGatewayClient.ListCertificates(context.Background(), listCertificatesRequest)

	if err != nil {
		return resourceIds, fmt.Errorf("Error getting Certificate list for compartment id : %s , %s \n", compartmentId, err)
	}
	for _, certificate := range listCertificatesResponse.Items {
		id := *certificate.Id
		resourceIds = append(resourceIds, id)
		AddResourceIdToSweeperResourceIdMap(compartmentId, "CertificateId", id)
	}
	return resourceIds, nil
}

func apiGatewayCertificateSweepWaitCondition(response common.OCIOperationResponse) bool {
	// Only stop if the resource is available beyond 3 mins. As there could be an issue for the sweeper to delete the resource and manual intervention required.
	if certificateResponse, ok := response.Response.(oci_apigateway.GetCertificateResponse); ok {
		return certificateResponse.LifecycleState != oci_apigateway.CertificateLifecycleStateDeleted
	}
	return false
}

func apiGatewayCertificateSweepResponseFetchOperation(client *OracleClients, resourceId *string, retryPolicy *common.RetryPolicy) error {
	_, err := client.apiGatewayClient().GetCertificate(context.Background(), oci_apigateway.GetCertificateRequest{
		CertificateId: resourceId,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: retryPolicy,
		},
	})
	return err
}
