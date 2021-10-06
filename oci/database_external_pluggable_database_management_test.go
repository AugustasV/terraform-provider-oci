// Copyright (c) 2017, 2020, Oracle and/or its affiliates. All rights reserved.
// Licensed under the Mozilla Public License v2.0

package oci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	externalPluggableDatabaseManagementRepresentation = map[string]interface{}{
		"external_pluggable_database_id": Representation{RepType: Required, Create: `${oci_database_external_pluggable_database.test_external_pluggable_database.id}`},
		"external_database_connector_id": Representation{RepType: Required, Create: `${oci_database_external_database_connector.test_external_pluggable_database_connector.id}`},
		"enable_management":              Representation{RepType: Required, Create: `true`, Update: `false`},
	}
	externalPluggableDatabaseConnectorRepresentation = map[string]interface{}{
		"connection_credentials": RepresentationGroup{Required, externalDatabaseConnectorConnectionCredentialsRepresentation},
		"connection_string":      RepresentationGroup{Required, externalDatabaseConnectorConnectionStringRepresentation},
		"connector_agent_id":     Representation{RepType: Required, Create: `ocid1.managementagent.oc1.phx.amaaaaaajobtc3iaes4ijczgekzqigoji25xocsny7yundummydummydummy`},
		"display_name":           Representation{RepType: Required, Create: `myTestConn`},
		"external_database_id":   Representation{RepType: Required, Create: `${oci_database_external_pluggable_database.test_external_pluggable_database.id}`},
		"connector_type":         Representation{RepType: Optional, Create: `MACS`},
	}

	externalPluggable1DatabaseRepresentation = map[string]interface{}{
		"compartment_id":                 Representation{RepType: Required, Create: `${var.compartment_id}`},
		"display_name":                   Representation{RepType: Required, Create: `myTestExternalPdb`},
		"external_container_database_id": Representation{RepType: Required, Create: `${oci_database_external_container_database.test_external_container_database.id}`},
		"defined_tags":                   Representation{RepType: Optional, Create: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "value")}`, Update: `${map("${oci_identity_tag_namespace.tag-namespace1.name}.${oci_identity_tag.tag1.name}", "updatedValue")}`},
		"freeform_tags":                  Representation{RepType: Optional, Create: map[string]string{"Department": "Finance"}, Update: map[string]string{"Department": "Accounting"}},
	}

	ExternalPluggableDatabaseManagementResourceDependencies = GenerateResourceFromRepresentationMap("oci_database_external_container_database", "test_external_container_database", Required, Create, externalContainerDatabaseRepresentation) +
		GenerateResourceFromRepresentationMap("oci_database_external_pluggable_database", "test_external_pluggable_database", Required, Create, externalPluggable1DatabaseRepresentation) +
		GenerateResourceFromRepresentationMap("oci_database_external_database_connector", "test_external_pluggable_database_connector", Required, Create, externalPluggableDatabaseConnectorRepresentation) +
		GenerateResourceFromRepresentationMap("oci_database_external_database_connector", "test_external_database_connector", Required, Create, externalContainerDatabaseConnectorRepresentation)
)

// issue-routing-tag: database/default
func TestDatabaseExternalPluggableDatabaseManagementResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestDatabaseExternalPluggableDatabaseManagementResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	resourceName := "oci_database_external_pluggable_database_management.test_external_pluggable_database_management"
	resourcePDB := "oci_database_external_pluggable_database.test_external_pluggable_database"

	// Save TF content to Create resource with only required properties. This has to be exactly the same as the config part in the Create step in the test.
	SaveConfigContent(config+compartmentIdVariableStr+ExternalPluggableDatabaseManagementResourceDependencies+
		GenerateResourceFromRepresentationMap("oci_database_external_pluggable_database_management", "test_external_pluggable_database_management", Required, Create, externalPluggableDatabaseManagementRepresentation), "database", "externalPluggableDatabaseManagement", t)

	ResourceTest(t, nil, []resource.TestStep{
		// Enablement of parent CDB
		{
			Config: config + compartmentIdVariableStr + ExternalPluggableDatabaseManagementResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_database_external_container_database_management", "test_external_container_database_management", Required, Create, externalContainerDatabaseManagementRepresentation),
		},
		// Enablement of PDB
		{
			Config: config + compartmentIdVariableStr + ExternalPluggableDatabaseManagementResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_database_external_container_database_management", "test_external_container_database_management", Required, Create, externalContainerDatabaseManagementRepresentation) +
				GenerateResourceFromRepresentationMap("oci_database_external_pluggable_database_management", "test_external_pluggable_database_management", Required, Create, externalPluggableDatabaseManagementRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(resourceName, "external_pluggable_database_id"),
				resource.TestCheckResourceAttrSet(resourceName, "external_database_connector_id"),
			),
		},
		// Verify Enablement of PDB
		{
			Config: config + compartmentIdVariableStr + ExternalPluggableDatabaseManagementResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_database_external_container_database_management", "test_external_container_database_management", Required, Create, externalContainerDatabaseManagementRepresentation) +
				GenerateResourceFromRepresentationMap("oci_database_external_pluggable_database_management", "test_external_pluggable_database_management", Required, Create, externalPluggableDatabaseManagementRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourcePDB, "database_management_config.0.database_management_status", "ENABLED"),
			),
		},

		// delete before next Create
		{
			Config: config + compartmentIdVariableStr + ExternalPluggableDatabaseManagementResourceDependencies,
		},
		// Enablement of parent CDB
		{
			Config: config + compartmentIdVariableStr + ExternalPluggableDatabaseManagementResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_database_external_container_database_management", "test_external_container_database_management", Optional, Create, externalContainerDatabaseManagementRepresentation),
		},
		// Enablement of PDB
		{
			Config: config + compartmentIdVariableStr + ExternalPluggableDatabaseManagementResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_database_external_container_database_management", "test_external_container_database_management", Optional, Create, externalContainerDatabaseManagementRepresentation) +
				GenerateResourceFromRepresentationMap("oci_database_external_pluggable_database_management", "test_external_pluggable_database_management", Optional, Create, externalPluggableDatabaseManagementRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(resourceName, "external_pluggable_database_id"),
				resource.TestCheckResourceAttrSet(resourceName, "external_database_connector_id"),
			),
		},
		// Verify Enablement of PDB
		{
			Config: config + compartmentIdVariableStr + ExternalPluggableDatabaseManagementResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_database_external_container_database_management", "test_external_container_database_management", Optional, Create, externalContainerDatabaseManagementRepresentation) +
				GenerateResourceFromRepresentationMap("oci_database_external_pluggable_database_management", "test_external_pluggable_database_management", Optional, Create, externalPluggableDatabaseManagementRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourcePDB, "database_management_config.0.database_management_status", "ENABLED"),
			),
		},
		// Disablement of parent CDB
		{
			Config: config + compartmentIdVariableStr + ExternalPluggableDatabaseManagementResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_database_external_container_database_management", "test_external_container_database_management", Optional, Update, externalContainerDatabaseManagementRepresentation) +
				GenerateResourceFromRepresentationMap("oci_database_external_pluggable_database_management", "test_external_pluggable_database_management", Optional, Update, externalPluggableDatabaseManagementRepresentation),
		},
		// Disablement of PDB
		{
			Config: config + compartmentIdVariableStr + ExternalPluggableDatabaseManagementResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_database_external_container_database_management", "test_external_container_database_management", Optional, Update, externalContainerDatabaseManagementRepresentation) +
				GenerateResourceFromRepresentationMap("oci_database_external_pluggable_database_management", "test_external_pluggable_database_management", Optional, Update, externalPluggableDatabaseManagementRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(resourceName, "external_pluggable_database_id"),
				resource.TestCheckResourceAttrSet(resourceName, "external_database_connector_id"),
			),
		},
		// Verify Disablement of PDB
		{
			Config: config + compartmentIdVariableStr + ExternalPluggableDatabaseManagementResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_database_external_container_database_management", "test_external_container_database_management", Optional, Update, externalContainerDatabaseManagementRepresentation) +
				GenerateResourceFromRepresentationMap("oci_database_external_pluggable_database_management", "test_external_pluggable_database_management", Optional, Update, externalPluggableDatabaseManagementRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourcePDB, "database_management_config.0.database_management_status", "NOT_ENABLED"),
			),
		},
	})
}
