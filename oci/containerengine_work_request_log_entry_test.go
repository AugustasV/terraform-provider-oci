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
	workRequestLogEntryDataSourceRepresentation = map[string]interface{}{
		"compartment_id":  Representation{RepType: Required, Create: `${var.compartment_id}`},
		"work_request_id": Representation{RepType: Required, Create: `${lookup(data.oci_containerengine_work_requests.test_work_requests.work_requests[0], "id")}`},
	}

	WorkRequestLogEntryResourceConfig = WorkRequestResourceConfig +
		GenerateDataSourceFromRepresentationMap("oci_containerengine_work_requests", "test_work_requests", Optional, Create, workRequestDataSourceRepresentation)
)

// issue-routing-tag: containerengine/default
func TestContainerengineWorkRequestLogEntryResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestContainerengineWorkRequestLogEntryResource_basic")
	defer httpreplay.SaveScenario()

	config := ProviderTestConfig()

	compartmentId := GetEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	datasourceName := "data.oci_containerengine_work_request_log_entries.test_work_request_log_entries"

	SaveConfigContent("", "", "", t)

	ResourceTest(t, nil, []resource.TestStep{
		// verify datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_containerengine_work_request_log_entries", "test_work_request_log_entries", Required, Create, workRequestLogEntryDataSourceRepresentation) +
				compartmentIdVariableStr + WorkRequestLogEntryResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(datasourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttrSet(datasourceName, "work_request_id"),

				resource.TestCheckResourceAttrSet(datasourceName, "work_request_log_entries.#"),
				resource.TestCheckResourceAttrSet(datasourceName, "work_request_log_entries.0.message"),
				resource.TestCheckResourceAttrSet(datasourceName, "work_request_log_entries.0.timestamp"),
			),
		},
	})
}
