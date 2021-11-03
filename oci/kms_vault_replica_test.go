// Copyright (c) 2017, 2021, Oracle and/or its affiliates. All rights reserved.
// Licensed under the Mozilla Public License v2.0

package oci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	vaultReplicaDataSourceRepresentation = map[string]interface{}{
		"vault_id": Representation{RepType: Required, Create: `${data.oci_kms_vault.test_vault.id}`},
	}

	VaultReplicaResourceConfig = GenerateResourceFromRepresentationMap("oci_kms_vault", "test_vault", Required, Create, vaultRepresentation)
)

// issue-routing-tag: kms/default
func TestKmsVaultReplicaResource_basic(t *testing.T) {
	t.Skip("Skip this test because virtual private vault is needed")
	httpreplay.SetScenario("TestKmsVaultReplicaResource_basic")
	defer httpreplay.SaveScenario()

	config := ProviderTestConfig()

	compartmentId := GetEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	datasourceName := "data.oci_kms_vault_replicas.test_vault_replicas"

	SaveConfigContent("", "", "", t)

	ResourceTest(t, nil, []resource.TestStep{
		// verify datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_kms_vault_replicas", "test_vault_replicas", Required, Create, vaultReplicaDataSourceRepresentation) +
				compartmentIdVariableStr + VaultReplicaResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(datasourceName, "vault_id"),

				resource.TestCheckResourceAttrSet(datasourceName, "vault_replicas.#"),
				resource.TestCheckResourceAttrSet(datasourceName, "vault_replicas.0.crypto_endpoint"),
				resource.TestCheckResourceAttrSet(datasourceName, "vault_replicas.0.management_endpoint"),
				resource.TestCheckResourceAttrSet(datasourceName, "vault_replicas.0.region"),
				resource.TestCheckResourceAttrSet(datasourceName, "vault_replicas.0.status"),
			),
		},
	})
}
