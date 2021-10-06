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
	instancePoolLoadBalancerAttachmentSingularDataSourceRepresentation = map[string]interface{}{
		"instance_pool_id":                          Representation{RepType: Required, Create: `${oci_core_instance_pool.test_instance_pool.id}`},
		"instance_pool_load_balancer_attachment_id": Representation{RepType: Required, Create: `${oci_core_instance_pool.test_instance_pool.load_balancers.0.id}`},
	}

	InstancePoolLoadBalancerAttachmentResourceConfig = OciImageIdsVariable +
		GenerateResourceFromRepresentationMap("oci_core_instance_configuration", "test_instance_configuration", Optional, Create, instanceConfigurationPoolRepresentation) +
		GenerateResourceFromRepresentationMap("oci_core_instance_pool", "test_instance_pool", Optional, Update, instancePoolRepresentation) +
		GenerateResourceFromRepresentationMap("oci_core_instance", "test_instance", Required, Create, instanceRepresentation) +
		GenerateResourceFromRepresentationMap("oci_core_network_security_group", "test_network_security_group", Required, Create, networkSecurityGroupRepresentation) +
		GenerateResourceFromRepresentationMap("oci_core_subnet", "test_subnet", Required, Create, subnetRepresentation) +
		GenerateResourceFromRepresentationMap("oci_core_vcn", "test_vcn", Required, Create, vcnRepresentation) +
		AvailabilityDomainConfig +
		DefinedTagsDependencies +
		GenerateResourceFromRepresentationMap("oci_load_balancer_backend_set", "test_backend_set", Required, Create, backendSetRepresentation) +
		GenerateResourceFromRepresentationMap("oci_load_balancer_certificate", "test_certificate", Required, Create, certificateRepresentation) +
		GenerateResourceFromRepresentationMap("oci_load_balancer_load_balancer", "test_load_balancer", Required, Create, loadBalancerRepresentation) +
		LoadBalancerSubnetDependencies
)

// issue-routing-tag: core/computeManagement
func TestCoreInstancePoolLoadBalancerAttachmentResource_basic(t *testing.T) {
	httpreplay.SetScenario("TestCoreInstancePoolLoadBalancerAttachmentResource_basic")
	defer httpreplay.SaveScenario()

	config := testProviderConfig()

	compartmentId := getEnvSettingWithBlankDefault("compartment_ocid")
	compartmentIdVariableStr := fmt.Sprintf("variable \"compartment_id\" { default = \"%s\" }\n", compartmentId)

	singularDatasourceName := "data.oci_core_instance_pool_load_balancer_attachment.test_instance_pool_load_balancer_attachment"

	SaveConfigContent("", "", "", t)

	ResourceTest(t, nil, []resource.TestStep{
		// verify singular datasource
		{
			Config: config +
				GenerateDataSourceFromRepresentationMap("oci_core_instance_pool_load_balancer_attachment", "test_instance_pool_load_balancer_attachment", Required, Create, instancePoolLoadBalancerAttachmentSingularDataSourceRepresentation) +
				compartmentIdVariableStr + InstancePoolLoadBalancerAttachmentResourceConfig,
			Check: ComposeAggregateTestCheckFuncWrapper(
				resource.TestCheckResourceAttrSet(singularDatasourceName, "instance_pool_id"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "instance_pool_load_balancer_attachment_id"),

				resource.TestCheckResourceAttrSet(singularDatasourceName, "backend_set_name"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "id"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "load_balancer_id"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "port"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "state"),
				resource.TestCheckResourceAttrSet(singularDatasourceName, "vnic_selection"),
			),
		},
	})
}
