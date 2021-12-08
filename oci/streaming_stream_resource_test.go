// Copyright (c) 2017, 2021, Oracle and/or its affiliates. All rights reserved.
// Licensed under the Mozilla Public License v2.0

package oci

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	"github.com/terraform-providers/terraform-provider-oci/httpreplay"
)

var (
	StreampoolidResourceConfig = StreampoolidResourceDependencies +
		GenerateResourceFromRepresentationMap("oci_streaming_stream", "test_stream", Optional, Update, streampoolidRepresentation)

	streampoolidSingularDataSourceRepresentation = map[string]interface{}{
		"stream_id": Representation{RepType: Required, Create: `${oci_streaming_stream.test_stream.id}`},
	}

	streampoolidDataSourceRepresentation = RepresentationCopyWithNewProperties(RepresentationCopyWithRemovedProperties(streamDataSourceRepresentation, []string{"compartment_id"}), map[string]interface{}{"stream_pool_id": Representation{RepType: Required, Create: `${oci_streaming_stream_pool.test_stream_pool.id}`}})

	streampoolidRepresentation = RepresentationCopyWithNewProperties(RepresentationCopyWithRemovedProperties(streamRepresentation, []string{"compartment_id"}), map[string]interface{}{"stream_pool_id": Representation{RepType: Required, Create: `${oci_streaming_stream_pool.test_stream_pool.id}`}})

	StreampoolidResourceDependencies = DefinedTagsDependencies +
		GenerateResourceFromRepresentationMap("oci_streaming_stream_pool", "test_stream_pool", Required, Create, streamPoolRepresentation)
)

// issue-routing-tag: streaming/default
func TestStreamingStreamWithStreamPoolIdResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestStreamingStreamResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	compartmentIdU := getEnvSettingWithDefault("compartment_id_for_update", compartmentId)
	compartmentIdUVariableStr := fmt.Sprintf("variable \"compartment_id_for_update\" { default = \"%s\" }\n", compartmentIdU)

	resourceName := "oci_streaming_stream.test_stream"
	datasourceName := "data.oci_streaming_streams.test_streams"
	singularDatasourceName := "data.oci_streaming_stream.test_stream"

	var resId, resId2 string

	ResourceTest(t, testAccCheckStreamingStreamDestroy, []resource.TestStep{
		// verify Create
		{
			Config: config + compartmentIdVariableStr + StreampoolidResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_streaming_stream", "test_stream", Required, Create, streampoolidRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(resourceName, "name", "mynewstream"),
				resource.TestCheckResourceAttr(resourceName, "partitions", "1"),

				func(s *terraform.State) (err error) {
					resId, err = FromInstanceState(s, resourceName, "id")
					return err
				},
			),
		},
		// Verify that stream's stream_pool_id can be removed and compartment_id can be used
		{
			Config: config + compartmentIdVariableStr + compartmentIdUVariableStr + StreampoolidResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_streaming_stream", "test_stream", Optional, Create,
					RepresentationCopyWithNewProperties(streamRepresentation, map[string]interface{}{
						"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id_for_update}`},
					})),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentIdU),
				resource.TestCheckResourceAttr(resourceName, "name", "mynewstream"),
				resource.TestCheckResourceAttr(resourceName, "partitions", "1"),

				func(s *terraform.State) (err error) {
					resId, err = FromInstanceState(s, resourceName, "id")
					return err
				},
			),
		},

		// delete before next Create
		{
			Config: config + compartmentIdVariableStr + StreampoolidResourceDependencies,
		},
		// verify Create with optionals
		{
			Config: config + compartmentIdVariableStr + StreampoolidResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_streaming_stream", "test_stream", Optional, Create, streampoolidRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttrSet(resourceName, "messages_endpoint"),
				resource.TestCheckResourceAttr(resourceName, "name", "mynewstream"),
				resource.TestCheckResourceAttr(resourceName, "partitions", "1"),
				resource.TestCheckResourceAttr(resourceName, "retention_in_hours", "24"),
				resource.TestCheckResourceAttrSet(resourceName, "state"),
				resource.TestCheckResourceAttrSet(resourceName, "stream_pool_id"),
				resource.TestCheckResourceAttrSet(resourceName, "time_created"),

				func(s *terraform.State) (err error) {
					resId, err = FromInstanceState(s, resourceName, "id")
					return err
				},
			),
		},
		{
			Config: config + compartmentIdVariableStr + compartmentIdUVariableStr + StreampoolidResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_streaming_stream_pool", "test_stream_pool_new", Required, Create, RepresentationCopyWithNewProperties(streamPoolRepresentation, map[string]interface{}{
					"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id_for_update}`},
				})) +
				GenerateResourceFromRepresentationMap("oci_streaming_stream", "test_stream", Optional, Create, streampoolidRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
			),
		},
		// verify stream move to compartment_id_for_update when stream_pool_id is updated which is also present in compartment_id_for_update
		{
			Config: config + compartmentIdVariableStr + compartmentIdUVariableStr + StreampoolidResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_streaming_stream_pool", "test_stream_pool_new", Required, Create, RepresentationCopyWithNewProperties(streamPoolRepresentation, map[string]interface{}{
					"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id_for_update}`},
				})) +
				GenerateResourceFromRepresentationMap("oci_streaming_stream", "test_stream", Optional, Create, RepresentationCopyWithNewProperties(streampoolidRepresentation, map[string]interface{}{
					"stream_pool_id": Representation{RepType: Required, Create: `${oci_streaming_stream_pool.test_stream_pool_new.id}`},
				})),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentIdU),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttrSet(resourceName, "messages_endpoint"),
				resource.TestCheckResourceAttr(resourceName, "name", "mynewstream"),
				resource.TestCheckResourceAttr(resourceName, "partitions", "1"),
				resource.TestCheckResourceAttr(resourceName, "retention_in_hours", "24"),
				resource.TestCheckResourceAttrSet(resourceName, "state"),
				resource.TestCheckResourceAttrSet(resourceName, "stream_pool_id"),
				resource.TestCheckResourceAttrSet(resourceName, "time_created"),

				func(s *terraform.State) (err error) {
					resId, err = FromInstanceState(s, resourceName, "id")
					return err
				},
			),
		},

		// verify updates to updatable parameters and switch stream back to test_stream_pool in compartment_ocid
		{
			Config: config + compartmentIdVariableStr + compartmentIdUVariableStr + StreampoolidResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_streaming_stream_pool", "test_stream_pool_new", Required, Create, RepresentationCopyWithNewProperties(streamPoolRepresentation, map[string]interface{}{
					"compartment_id": Representation{RepType: Required, Create: `${var.compartment_id_for_update}`},
				})) +
				GenerateResourceFromRepresentationMap("oci_streaming_stream", "test_stream", Optional, Update, streampoolidRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttr(resourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(resourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(resourceName, "id"),
				resource.TestCheckResourceAttrSet(resourceName, "messages_endpoint"),
				resource.TestCheckResourceAttr(resourceName, "name", "mynewstream"),
				resource.TestCheckResourceAttr(resourceName, "partitions", "1"),
				resource.TestCheckResourceAttr(resourceName, "retention_in_hours", "24"),
				resource.TestCheckResourceAttrSet(resourceName, "state"),
				resource.TestCheckResourceAttrSet(resourceName, "stream_pool_id"),
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
				GenerateDataSourceFromRepresentationMap("oci_streaming_streams", "test_streams", Optional, Update, streampoolidDataSourceRepresentation) +
				compartmentIdVariableStr + StreampoolidResourceDependencies +
				GenerateResourceFromRepresentationMap("oci_streaming_stream", "test_stream", Optional, Update, streampoolidRepresentation),
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(datasourceName, "id"),
				resource.TestCheckResourceAttr(datasourceName, "name", "mynewstream"),
				resource.TestCheckResourceAttr(datasourceName, "state", "ACTIVE"),
				resource.TestCheckResourceAttrSet(datasourceName, "stream_pool_id"),

				resource.TestCheckResourceAttr(datasourceName, "streams.#", "1"),
				resource.TestCheckResourceAttr(datasourceName, "streams.0.compartment_id", compartmentId),
				resource.TestCheckResourceAttr(datasourceName, "streams.0.freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(datasourceName, "streams.0.id"),
				resource.TestCheckResourceAttrSet(datasourceName, "streams.0.messages_endpoint"),
				resource.TestCheckResourceAttr(datasourceName, "streams.0.name", "mynewstream"),
				resource.TestCheckResourceAttr(datasourceName, "streams.0.partitions", "1"),
				resource.TestCheckResourceAttrSet(datasourceName, "streams.0.state"),
				resource.TestCheckResourceAttrSet(datasourceName, "streams.0.stream_pool_id"),
				resource.TestCheckResourceAttrSet(datasourceName, "streams.0.time_created"),
			),
		},
		// verify singular datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_streaming_stream", "test_stream", Required, Create, streampoolidSingularDataSourceRepresentation) +
				compartmentIdVariableStr + StreampoolidResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(singularDatasourceName, "stream_id"),

				resource.TestCheckResourceAttr(singularDatasourceName, "compartment_id", compartmentId),
				resource.TestCheckResourceAttr(singularDatasourceName, "freeform_tags.%", "1"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "id"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "messages_endpoint"),
				resource.TestCheckResourceAttr(singularDatasourceName, "name", "mynewstream"),
				resource.TestCheckResourceAttr(singularDatasourceName, "partitions", "1"),
				resource.TestCheckResourceAttr(singularDatasourceName, "retention_in_hours", "24"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "state"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "time_created"),
			),
		},
		// remove singular datasource from previous step so that it doesn't conflict with import tests
		{
			Config: config + compartmentIdVariableStr + StreampoolidResourceConfig,
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
