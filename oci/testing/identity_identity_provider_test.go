// Copyright (c) 2017, 2021, Oracle and/or its affiliates. All rights reserved.
// Licensed under the Mozilla Public License v2.0

package oci

import (
	"context"
	"fmt"
	"strconv"
	"testing"

	"io/ioutil"
	"log"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/oracle/oci-go-sdk/v49/common"
	oci_identity "github.com/oracle/oci-go-sdk/v49/identity"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

const (
	IdentityProviderPropertyVariables = `
variable "identity_provider_metadata" { default = "" }
variable "identity_provider_metadata_file" { default = "{{.metadata_file}}" }
`
)

var (
	IdentityProviderRequiredOnlyResource = IdentityProviderResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_identity_identity_provider", "test_identity_provider", Required, Create, identityProviderRepresentation)

	identityProviderDataSourceRepresentation = map[string]interface{}{
		"compartment_id": Representation{RepType: Required, Create: `${var.tenancy_ocid}`},
		"protocol":       Representation{RepType: Required, Create: `SAML2`},
		"name":           Representation{RepType: Optional, Create: `test-idp-saml2-adfs`},
		"state":          Representation{RepType: Optional, Create: `ACTIVE`},
		"filter":         RepresentationGroup{Required, identityProviderDataSourceFilterRepresentation}}
	identityProviderDataSourceFilterRepresentation = map[string]interface{}{
		"name":   Representation{RepType: Required, Create: `id`},
		"values": Representation{RepType: Required, Create: []string{`${oci_identity_identity_provider.test_identity_provider.id}`}},
	}

	identityProviderRepresentation = map[string]interface{}{
		"compartment_id":      Representation{RepType: Required, Create: `${var.tenancy_ocid}`},
		"description":         Representation{RepType: Required, Create: `description`, Update: `description2`},
		"metadata":            Representation{RepType: Required, Create: `${file("${var.identity_provider_metadata_file}")}`},
		"metadata_url":        Representation{RepType: Required, Create: `metadataUrl`, Update: `metadataUrl2`},
		"name":                Representation{RepType: Required, Create: `test-idp-saml2-adfs`},
		"product_type":        Representation{RepType: Required, Create: `ADFS`},
		"protocol":            Representation{RepType: Required, Create: `SAML2`},
		"defined_tags":        Representation{RepType: Optional, Create: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "value")}`, Update: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "updatedValue")}`},
		"freeform_attributes": Representation{RepType: Optional, Create: map[string]string{"clientId": "app_sf3kdjf3"}},
		"freeform_tags":       Representation{RepType: Optional, Create: map[string]string{"Department": "Finance"}, Update: map[string]string{"Department": "Accounting"}},
	}

	IdentityProviderResourceDependencies = IdentityProviderPropertyVariables +
		DefinedTagsDependencies
)

// issue-routing-tag: identity/default
func TestIdentityIdentityProviderResource_basic(t *testing.T) {
	metadataFile := GetEnvSettingWithBlankDefault("identity_provider_metadata_file")
	if metadataFile == "" {
		t.Skip("Skipping generated test for now as it has a dependency on federation metadata file")
	}

	httpreplay.SetScenario("TestIdentityIdentityProviderResource_basic")
	defer httpreplay.SaveScenario()

	config := ProviderTestConfig()

	compartmentId := GetEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)
	tenancyId := GetEnvSettingWithBlankDefault("tenancy_ocid")

	resourceName := "oci_identity_identity_provider.test_identity_provider"
	datasourceName := "data.oci_identity_identity_providers.test_identity_providers"

	var resId, resId2 string
	// Save TF content to Create resource with optional properties. This has to be exactly the same as the config part in the "Create with optionals" step in the test.
	SaveConfigContent(config+compartmentIdVariableStr+IdentityProviderResourceDependencies+
		GenerateResourceFromRepresentationMap("oci_identity_identity_provider", "test_identity_provider", Optional, Create, identityProviderRepresentation), "identity", "identityProvider", t)

	metadataContents, err := ioutil.ReadFile(metadataFile)
	if err != nil {
		log.Panic("Unable to read the file ", metadataFile)
	}
	metadata := string(metadataContents)

	_, tokenFn := TokenizeWithHttpReplay("identity_provider")
	IdentityProviderResourceDependencies = tokenFn(IdentityProviderResourceDependencies, map[string]string{"metadata_file": metadataFile})

	ResourceTest(t, testAccCheckIdentityIdentityProviderDestroy, []resource.TestStep{
		// verify Create
		{
			Config: config + compartmentIdVariableStr + IdentityProviderResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_identity_identity_provider", "test_identity_provider", Required, Create, identityProviderRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", tenancyId),
				resource.TestCheckResourceAttr(resourceName, "description", "description"),
				resource.TestCheckResourceAttr(resourceName, "metadata", metadata),
				resource.TestCheckResourceAttr(resourceName, "metadata_url", "metadataUrl"),
				resource.TestCheckResourceAttr(resourceName, "name", "test-idp-saml2-adfs"),
				resource.TestCheckResourceAttr(resourceName, "product_type", "ADFS"),
				resource.TestCheckResourceAttr(resourceName, "protocol", "SAML2"),

				func(s *terraform.State) (err error) {
					resId, err = FromInstanceState(s, resourceName, "id")
					return err
				},
			),
		},

		// delete before next Create
		{
			Config: config + compartmentIdVariableStr + IdentityProviderResourceDependencies,
		},

		// verify Create with optionals
		{
			Config: config + compartmentIdVariableStr + IdentityProviderResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_identity_identity_provider", "test_identity_provider", Optional, Create, identityProviderRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", tenancyId),
				resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "description", "description"),
				resource.TestCheckResourceAttr(resourceName, "freeform_attributes.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttr(resourceName, "metadata", metadata),
				resource.TestCheckResourceAttr(resourceName, "metadata_url", "metadataUrl"),
				resource.TestCheckResourceAttr(resourceName, "name", "test-idp-saml2-adfs"),
				resource.TestCheckResourceAttr(resourceName, "product_type", "ADFS"),
				resource.TestCheckResourceAttr(resourceName, "protocol", "SAML2"),
				resource.TestCheckResourceAttrSet(resourceName, "redirect_url"),
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
			Config: config + compartmentIdVariableStr + IdentityProviderResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_identity_identity_provider", "test_identity_provider", Optional, Update, identityProviderRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", tenancyId),
				resource.TestCheckResourceAttr(resourceName, "defined_tags.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "description", "description2"),
				resource.TestCheckResourceAttr(resourceName, "freeform_attributes.%", "1"),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttr(resourceName, "metadata", metadata),
				resource.TestCheckResourceAttr(resourceName, "metadata_url", "metadataUrl2"),
				resource.TestCheckResourceAttr(resourceName, "name", "test-idp-saml2-adfs"),
				resource.TestCheckResourceAttr(resourceName, "product_type", "ADFS"),
				resource.TestCheckResourceAttr(resourceName, "protocol", "SAML2"),
				resource.TestCheckResourceAttrSet(resourceName, "redirect_url"),
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
				GenerateDataSourceFromRepresentationMap("oci_identity_identity_providers", "test_identity_providers", Optional, Update, identityProviderDataSourceRepresentation) +
				compartmentIdVariableStr + IdentityProviderResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_identity_identity_provider", "test_identity_provider", Optional, Update, identityProviderRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(datasourceName, "compartment_id", tenancyId),
				resource.TestCheckResourceAttr(datasourceName, "name", "test-idp-saml2-adfs"),
				resource.TestCheckResourceAttr(datasourceName, "protocol", "SAML2"),
				resource.TestCheckResourceAttr(datasourceName, "state", "ACTIVE"),

				resource.TestCheckResourceAttr(datasourceName, "identity_providers.#", "1"),
				resource.TestCheckResourceAttr(datasourceName, "identity_providers.0.compartment_id", tenancyId),
				resource.TestCheckResourceAttr(datasourceName, "identity_providers.0.defined_tags.%", "1"),
				resource.TestCheckResourceAttr(datasourceName, "identity_providers.0.description", "description2"),
				resource.TestCheckResourceAttr(datasourceName, "identity_providers.0.freeform_attributes.%", "1"),
				resource.TestCheckResourceAttr(datasourceName, "identity_providers.0.freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(datasourceName, "identity_providers.0.id"),
				resource.TestCheckResourceAttrSet(datasourceName, "identity_providers.0.metadata"),
				resource.TestCheckResourceAttr(datasourceName, "identity_providers.0.metadata_url", "metadataUrl2"),
				resource.TestCheckResourceAttr(datasourceName, "identity_providers.0.name", "test-idp-saml2-adfs"),
				resource.TestCheckResourceAttr(datasourceName, "identity_providers.0.product_type", "ADFS"),
				resource.TestCheckResourceAttr(datasourceName, "identity_providers.0.protocol", "SAML2"),
				resource.TestCheckResourceAttrSet(datasourceName, "identity_providers.0.redirect_url"),
				resource.TestCheckResourceAttrSet(datasourceName, "identity_providers.0.state"),
				resource.TestCheckResourceAttrSet(datasourceName, "identity_providers.0.time_created"),
			),
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

func testAccCheckIdentityIdentityProviderDestroy(s *terraform.State) error {
	noResourceFound := true
	client := TestAccProvider.Meta().(*OracleClients).identityClient()
	for _, rs := range s.RootModule().Resources {
		if rs.Type == "oci_identity_identity_provider" {
			noResourceFound = false
			request := oci_identity.GetIdentityProviderRequest{}

			tmp := rs.Primary.ID
			request.IdentityProviderId = &tmp

			request.RequestMetadata.RetryPolicy = GetRetryPolicy(true, "identity")

			response, err := client.GetIdentityProvider(context.Background(), request)

			if err == nil {
				deletedLifecycleStates := map[string]bool{
					string(oci_identity.IdentityProviderLifecycleStateDeleted): true,
				}
				if _, ok := deletedLifecycleStates[string(response.GetLifecycleState())]; !ok {
					//resource lifecycle state is not in expected deleted lifecycle states.
					return fmt.Errorf("resource lifecycle state: %s is not in expected deleted lifecycle states", response.GetLifecycleState())
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
