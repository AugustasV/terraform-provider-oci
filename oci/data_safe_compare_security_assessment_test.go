// Copyright (c) 2017, 2021, Oracle and/or its affiliates. All rights reserved.
// Licensed under the Mozilla Public License v2.0

package oci

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	compareSecurityAssessmentRepresentation = map[string]interface{}{
		"comparison_security_assessment_id": Representation{RepType: Required, Create: `${oci_data_safe_security_assessment.test_security_assessment1.id}`},
		"security_assessment_id":            Representation{RepType: Required, Create: `${oci_data_safe_security_assessment.test_security_assessment2.id}`},
	}

	CompareSecurityAssessmentResourceDependencies = GenerateResourceFromRepresentationMap("oci_database_autonomous_database", "test_autonomous_database", Required, Create, autonomousDatabaseRepresentation) +
		GenerateResourceFromRepresentationMap("oci_data_safe_target_database", "test_target_database", Required, Create, targetDatabaseRepresentation) +
		GenerateResourceFromRepresentationMap("oci_data_safe_security_assessment", "test_security_assessment1", Required, Create, securityAssessmentRepresentation) +
		GenerateResourceFromRepresentationMap("oci_data_safe_security_assessment", "test_security_assessment2", Required, Create, securityAssessmentRepresentation)
)

func TestDataSafeCompareSecurityAssessmentResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestDataSafeCompareSecurityAssessmentResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	resourceName := "oci_data_safe_compare_security_assessment.test_compare_security_assessment"

	var resId string
	// Save TF content to Create resource with only required properties. This has to be exactly the same as the config part in the Create step in the test.
	SaveConfigContent(config+compartmentIdVariableStr+CompareSecurityAssessmentResourceDependencies+
		GenerateResourceFromRepresentationMap("oci_data_safe_compare_security_assessment", "test_compare_security_assessment", Required, Create, compareSecurityAssessmentRepresentation), "datasafe", "compareSecurityAssessment", t)

	ResourceTest(t, nil, []resource.TestStep{
		// verify Create
		{
			Config: config + compartmentIdVariableStr + CompareSecurityAssessmentResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_data_safe_compare_security_assessment", "test_compare_security_assessment", Required, Create, compareSecurityAssessmentRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(resourceName, "comparison_security_assessment_id"),
				resource.TestCheckResourceAttrSet(resourceName, "security_assessment_id"),
				resource.TestCheckResourceAttr(resourceName, "summary.#", "0"),

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
	})
}
