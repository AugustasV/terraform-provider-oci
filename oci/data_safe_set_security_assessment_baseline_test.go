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
	setSecurityAssessmentBaselineRepresentation = map[string]interface{}{
		"security_assessment_id": Representation{RepType: Required, Create: `${oci_data_safe_security_assessment.test_security_assessment1.id}`},
		"assessment_ids":         Representation{RepType: Optional, Create: []string{`${oci_data_safe_security_assessment.test_security_assessment2.id}`}},
	}

	unsetSecAssessmentBaselineRepresentation = map[string]interface{}{
		"security_assessment_id": Representation{RepType: Required, Create: `${oci_data_safe_security_assessment.test_security_assessment1.id}`},
	}

	SetSecurityAssessmentBaselineResourceDependencies = GenerateResourceFromRepresentationMap("oci_database_autonomous_database", "test_autonomous_database", Required, Create, autonomousDatabaseRepresentation) +
		GenerateResourceFromRepresentationMap("oci_data_safe_target_database", "test_target_database", Required, Create, targetDatabaseRepresentation) +
		GenerateResourceFromRepresentationMap("oci_data_safe_security_assessment", "test_security_assessment1", Required, Create, securityAssessmentRepresentation) +
		GenerateResourceFromRepresentationMap("oci_data_safe_security_assessment", "test_security_assessment2", Required, Create, securityAssessmentRepresentation)
)

func TestDataSafeSetSecurityAssessmentBaselineResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestDataSafeSetSecurityAssessmentBaselineResource_basic")
	defer httpreplay.SaveScenario()

	config := ProviderTestConfig()

	compartmentId := GetEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	resourceName := "oci_data_safe_set_security_assessment_baseline.test_set_security_assessment_baseline"

	var resId string
	// Save TF content to Create resource with optional properties. This has to be exactly the same as the config part in the "Create with optionals" step in the test.
	SaveConfigContent(config+compartmentIdVariableStr+SetSecurityAssessmentBaselineResourceDependencies+
		GenerateResourceFromRepresentationMap("oci_data_safe_set_security_assessment_baseline", "test_set_security_assessment_baseline", Optional, Create, setSecurityAssessmentBaselineRepresentation), "datasafe", "setSecurityAssessmentBaseline", t)

	ResourceTest(t, nil, []resource.TestStep{
		// verify Create
		{
			Config: config + compartmentIdVariableStr + SetSecurityAssessmentBaselineResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_data_safe_set_security_assessment_baseline", "test_set_security_assessment_baseline", Optional, Create, setSecurityAssessmentBaselineRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(resourceName, "security_assessment_id"),
			),
		},

		// delete before next Create
		{
			Config: config + compartmentIdVariableStr + SetSecurityAssessmentBaselineResourceDependencies,
		},
		// verify Create with optionals
		{
			Config: config + compartmentIdVariableStr + SetSecurityAssessmentBaselineResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_data_safe_set_security_assessment_baseline", "test_set_security_assessment_baseline", Optional, Create, setSecurityAssessmentBaselineRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "assessment_ids.#", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "security_assessment_id"),

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
		{
			Config: config + compartmentIdVariableStr + SetSecurityAssessmentBaselineResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_data_safe_unset_security_assessment_baseline", "test_unset_security_assessment_baseline", Required, Create, unsetSecAssessmentBaselineRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				func(s *terraform.State) (err error) {
					return err
				},
			),
		},
	})
}
