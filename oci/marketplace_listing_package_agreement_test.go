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
	listingPackageAgreementManagementRepresentation = map[string]interface{}{
		"agreement_id":    Representation{RepType: Required, Create: `${data.oci_marketplace_listing_package_agreements.test_listing_package_agreements.agreements.0.id}`},
		"listing_id":      Representation{RepType: Required, Create: `${data.oci_marketplace_listing.test_listing.id}`},
		"package_version": Representation{RepType: Required, Create: `${data.oci_marketplace_listing.test_listing.default_package_version}`},
		"compartment_id":  Representation{RepType: Optional, Create: `${var.compartment_id}`},
	}

	listingPackageAgreementDataSourceRepresentation = map[string]interface{}{
		"listing_id":      Representation{RepType: Required, Create: `${data.oci_marketplace_listing.test_listing.id}`},
		"package_version": Representation{RepType: Required, Create: `${data.oci_marketplace_listing.test_listing.default_package_version}`},
		"compartment_id":  Representation{RepType: Optional, Create: `${var.compartment_id}`},
	}

	ListingPackageAgreementResourceConfig = GenerateDataSourceFromRepresentationMap("oci_marketplace_listing", "test_listing", Required, Create, listingSingularDataSourceRepresentation) +
		GenerateDataSourceFromRepresentationMap("oci_marketplace_listings", "test_listings", Required, Create, listingDataSourceRepresentation)
)

// issue-routing-tag: marketplace/default
func TestMarketplaceListingPackageAgreementResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestMarketplaceListingPackageAgreementResource_basic")
	defer httpreplay.SaveScenario()

	config := ProviderTestConfig()

	compartmentId := GetEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	datasourceName := "data.oci_marketplace_listing_package_agreements.test_listing_package_agreements"
	resourceName := "oci_marketplace_listing_package_agreement.test_listing_package_agreement"

	SaveConfigContent("", "", "", t)

	ResourceTest(t, nil, []resource.TestStep{
		// verify resource
		{
			Config: config +
				GenerateResourceFromRepresentationMap("oci_marketplace_listing_package_agreement", "test_listing_package_agreement", Required, Create, listingPackageAgreementManagementRepresentation) +
				GenerateDataSourceFromRepresentationMap("oci_marketplace_listing_package_agreements", "test_listing_package_agreements", Required, Create, listingPackageAgreementDataSourceRepresentation) +
				compartmentIdVariableStr + ListingPackageAgreementResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(resourceName, "agreement_id"),
				resource.TestCheckResourceAttrSet(resourceName, "listing_id"),
				resource.TestCheckResourceAttrSet(resourceName, "package_version"),

				resource.TestCheckResourceAttrSet(resourceName, "content_url"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttrSet(resourceName, "prompt"),
				resource.TestCheckResourceAttrSet(resourceName, "signature"),
			),
		},
		// verify datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_marketplace_listing_package_agreements", "test_listing_package_agreements", Required, Create, listingPackageAgreementDataSourceRepresentation) +
				compartmentIdVariableStr + ListingPackageAgreementResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(datasourceName, "listing_id"),

				resource.TestCheckResourceAttrSet(datasourceName, "agreements.#"),
				resource.TestCheckResourceAttrSet(datasourceName, "agreements.0.author"),
				resource.TestCheckResourceAttrSet(datasourceName, "agreements.0.content_url"),
				resource.TestCheckResourceAttrSet(datasourceName, "agreements.0.id"),
				resource.TestCheckResourceAttrSet(datasourceName, "agreements.0.prompt"),
			),
		},
	})
}
