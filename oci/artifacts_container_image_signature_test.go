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
	oci_artifacts "github.com/oracle/oci-go-sdk/v55/artifacts"
	"github.com/oracle/oci-go-sdk/v55/common"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	message             = getEnvSettingWithBlankDefault("container_image_signing_signature")
	signingAlgorithm    = "SHA_224_RSA_PKCS_PSS"
	signingAlgorithmStr = fmt.Sprintf("variable \"signingAlgorithm\" { default = \"%s\" }\n", signingAlgorithm)
	description         = "Image built by TC"
	descriptionStr      = fmt.Sprintf("variable \"description\" { default = \"%s\" }\n", description)
	metadata            = "{\\\"buildNumber\\\":\\\"123\\\"}"
	metadataStr         = fmt.Sprintf("variable \"metadata\" { default = \"%s\" }\n", metadata)

	containerImageResourceConfig = GenerateDataSourceFromRepresentationMap("oci_artifacts_container_image", "test_container_image", Required, Create, containerImageSingularDataSourceRepresentation)

	containerImageSignatureKmsSignResourceDependencies = SignResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_kms_sign", "test_container_image_signature_kms_sign", Required, Create, containerImageSignatureKmsSignRepresentation)

	containerImageSignatureKmsSignRepresentation = map[string]interface{}{
		"crypto_endpoint":   Representation{RepType: Required, Create: `${data.oci_kms_vault.test_vault.crypto_endpoint}`},
		"key_id":            Representation{RepType: Required, Create: `${lookup(data.oci_kms_keys.test_keys_dependency_RSA.keys[0], "id")}`},
		"message":           Representation{RepType: Required, Create: message},
		"signing_algorithm": Representation{RepType: Required, Create: signingAlgorithm},
		"message_type":      Representation{RepType: Optional, Create: `RAW`},
	}

	containerImageSignatureRepresentation = map[string]interface{}{
		"compartment_id":     Representation{RepType: Required, Create: `${data.oci_artifacts_container_image.test_container_image.compartment_id}`},
		"image_id":           Representation{RepType: Required, Create: `${data.oci_artifacts_container_image.test_container_image.image_id}`},
		"kms_key_id":         Representation{RepType: Required, Create: `${lookup(data.oci_kms_keys.test_keys_dependency_RSA.keys[0], "id")}`},
		"kms_key_version_id": Representation{RepType: Required, Create: `${oci_kms_sign.test_container_image_signature_kms_sign.key_version_id}`},
		"message":            Representation{RepType: Required, Create: message},
		"signature":          Representation{RepType: Required, Create: `${oci_kms_sign.test_container_image_signature_kms_sign.signature}`},
		"signing_algorithm":  Representation{RepType: Required, Create: signingAlgorithm},
	}

	containerImageSignatureSingularDataSourceRepresentation = map[string]interface{}{
		"image_signature_id": Representation{RepType: Required, Create: `${oci_artifacts_container_image_signature.test_container_image_signature.id}`},
	}

	containerImageSignatureDataSourceRepresentation = map[string]interface{}{
		"compartment_id":            Representation{RepType: Required, Create: `${data.oci_artifacts_container_image.test_container_image.compartment_id}`},
		"compartment_id_in_subtree": Representation{RepType: Optional, Create: `false`},
		"display_name":              Representation{RepType: Optional, Create: `${oci_artifacts_container_image_signature.test_container_image_signature.display_name}`},
		"image_digest":              Representation{RepType: Optional, Create: `${data.oci_artifacts_container_image.test_container_image.digest}`},
		"image_id":                  Representation{RepType: Optional, Create: `${data.oci_artifacts_container_image.test_container_image.image_id}`},
		"kms_key_id":                Representation{RepType: Optional, Create: `${lookup(data.oci_kms_keys.test_keys_dependency_RSA.keys[0], "id")}`},
		"kms_key_version_id":        Representation{RepType: Optional, Create: `${oci_kms_sign.test_container_image_signature_kms_sign.key_version_id}`},
		"repository_id":             Representation{RepType: Optional, Create: `${data.oci_artifacts_container_image.test_container_image.repository_id}`},
		"repository_name":           Representation{RepType: Optional, Create: `${data.oci_artifacts_container_image.test_container_image.repository_name}`},
		"signing_algorithm":         Representation{RepType: Optional, Create: signingAlgorithm},
		"filter":                    RepresentationGroup{Required, containerImageSignatureDataSourceFilterRepresentation},
	}

	containerImageSignatureDataSourceFilterRepresentation = map[string]interface{}{
		"name":   Representation{RepType: Required, Create: `id`},
		"values": Representation{RepType: Required, Create: []string{`${oci_artifacts_container_image_signature.test_container_image_signature.id}`}},
	}
)

// issue-routing-tag: artifacts/default
func TestArtifactsContainerImageSignatureResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestArtifactsContainerImageSignatureResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	resourceName := "oci_artifacts_container_image_signature.test_container_image_signature"
	datasourceName := "data.oci_artifacts_container_image_signatures.test_container_image_signatures"
	singularDatasourceName := "data.oci_artifacts_container_image_signature.test_container_image_signature"

	var resId string

	ResourceTest(t, testAccCheckArtifactsContainerImageSignatureDestroy, []resource.TestStep{
		// verify Create
		{
			Config: config + descriptionStr + metadataStr + signingAlgorithmStr +
				containerImageResourceConfig +
				containerImageSignatureKmsSignResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_artifacts_container_image_signature", "test_container_image_signature", Required, Create, containerImageSignatureRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttrSet(resourceName, "created_by"),
				resource.TestCheckResourceAttrSet(resourceName, "display_name"),
				resource.TestCheckResourceAttrSet(resourceName, "image_id"),
				resource.TestCheckResourceAttrSet(resourceName, "kms_key_id"),
				resource.TestCheckResourceAttrSet(resourceName, "kms_key_version_id"),
				resource.TestCheckResourceAttrSet(resourceName, "message"),
				resource.TestCheckResourceAttrSet(resourceName, "signature"),
				resource.TestCheckResourceAttrSet(resourceName, "time_created"),
				resource.TestCheckResourceAttr(resourceName, "signing_algorithm", signingAlgorithm),

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

		// verify datasource
		{
			Config: config + descriptionStr + metadataStr + signingAlgorithmStr +
				containerImageResourceConfig +
				containerImageSignatureKmsSignResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_artifacts_container_image_signature", "test_container_image_signature", Optional, Update, containerImageSignatureRepresentation) +
				GenerateDataSourceFromRepresentationMap("oci_artifacts_container_image_signatures", "test_container_image_signatures", Optional, Update, containerImageSignatureDataSourceRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(datasourceName, "compartment_id_in_subtree", "false"),
				resource.TestCheckResourceAttrSet(datasourceName, "display_name"),
				resource.TestCheckResourceAttrSet(datasourceName, "image_digest"),
				resource.TestCheckResourceAttrSet(datasourceName, "image_id"),
				resource.TestCheckResourceAttrSet(datasourceName, "kms_key_id"),
				resource.TestCheckResourceAttrSet(datasourceName, "kms_key_version_id"),
				resource.TestCheckResourceAttrSet(datasourceName, "repository_id"),
				resource.TestCheckResourceAttrSet(datasourceName, "repository_name"),
				resource.TestCheckResourceAttr(datasourceName, "signing_algorithm", signingAlgorithm),

				resource.TestCheckResourceAttr(datasourceName, "container_image_signature_collection.#", "1"),
				resource.TestCheckResourceAttr(datasourceName, "container_image_signature_collection.0.items.#", "1"),
			),
		},

		// verify singular datasource
		{
			Config: config + descriptionStr + metadataStr + signingAlgorithmStr +
				containerImageResourceConfig +
				containerImageSignatureKmsSignResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_artifacts_container_image_signature", "test_container_image_signature", Optional, Update, containerImageSignatureRepresentation) +
				GenerateDataSourceFromRepresentationMap("oci_artifacts_container_image_signature", "test_container_image_signature", Required, Create, containerImageSignatureSingularDataSourceRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(singularDatasourceName, "image_signature_id"),
				resource.TestCheckResourceAttr(singularDatasourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "created_by"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "display_name"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "id"),
				//resource.TestCheckResourceAttr(singularDatasourceName, "message", encodedMessage),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "signature"),
				resource.TestCheckResourceAttr(singularDatasourceName, "signing_algorithm", signingAlgorithm),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "time_created"),
			),
		},

		// remove singular datasource from previous step so that it doesn't conflict with import tests
		{
			Config: config + descriptionStr + metadataStr + signingAlgorithmStr +
				containerImageResourceConfig +
				containerImageSignatureKmsSignResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_artifacts_container_image_signature", "test_container_image_signature", Optional, Update, containerImageSignatureRepresentation),
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

func testAccCheckArtifactsContainerImageSignatureDestroy(s *terraform.State) error {
	noResourceFound := true
	client := testAccProvider.Meta().(*OracleClients).artifactsClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oci_artifacts_container_image_signature" {
			noResourceFound = false
			request := oci_artifacts.GetContainerImageSignatureRequest{}

			if value, ok := rs.Primary.Attributes["id"]; ok {
				request.ImageSignatureId = &value
			}

			request.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "artifacts")

			_, err := client.GetContainerImageSignature(context.Background(), request)

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
	if !InSweeperExcludeList("ArtifactsContainerImageSignature") {
		resource.AddTestSweepers("ArtifactsContainerImageSignature", &resource.Sweeper{
			Name:         "ArtifactsContainerImageSignature",
			Dependencies: DependencyGraph["containerImageSignature"],
			F:            sweepArtifactsContainerImageSignatureResource,
		})
	}
}

func sweepArtifactsContainerImageSignatureResource(compartment string) error {
	artifactsClient := GetTestClients(&schema.ResourceData{}).artifactsClient()
	containerImageSignatureIds, err := getContainerImageSignatureIds(compartment)
	if err != nil {
		return err
	}
	for _, containerImageSignatureId := range containerImageSignatureIds {
		if ok := SweeperDefaultResourceId[containerImageSignatureId]; !ok {
			deleteContainerImageSignatureRequest := oci_artifacts.DeleteContainerImageSignatureRequest{}

			deleteContainerImageSignatureRequest.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "artifacts")
			_, error := artifactsClient.DeleteContainerImageSignature(context.Background(), deleteContainerImageSignatureRequest)
			if error != nil {
				fmt.Printf("Error deleting ContainerImageSignature %s %s, It is possible that the resource is already deleted. Please verify manually \n", containerImageSignatureId, error)
				continue
			}
		}
	}
	return nil
}

func getContainerImageSignatureIds(compartment string) ([]string, error) {
	ids := GetResourceIdsToSweep(compartment, "ContainerImageSignatureId")
	if ids != nil {
		return ids, nil
	}
	var resourceIds []string
	compartmentId := compartment
	artifactsClient := GetTestClients(&schema.ResourceData{}).artifactsClient()

	listContainerImageSignaturesRequest := oci_artifacts.ListContainerImageSignaturesRequest{}
	listContainerImageSignaturesRequest.CompartmentId = &compartmentId
	listContainerImageSignaturesResponse, err := artifactsClient.ListContainerImageSignatures(context.Background(), listContainerImageSignaturesRequest)

	if err != nil {
		return resourceIds, fmt.Errorf("Error getting ContainerImageSignature list for compartment id : %s , %s \n", compartmentId, err)
	}
	for _, containerImageSignature := range listContainerImageSignaturesResponse.Items {
		id := *containerImageSignature.Id
		resourceIds = append(resourceIds, id)
		AddResourceIdToSweeperResourceIdMap(compartmentId, "ContainerImageSignatureId", id)
	}
	return resourceIds, nil
}
