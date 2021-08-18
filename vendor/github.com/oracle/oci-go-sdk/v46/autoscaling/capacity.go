// Copyright (c) 2016, 2018, 2021, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.
// Code generated. DO NOT EDIT.

// Autoscaling API
//
// APIs for dynamically scaling Compute resources to meet application requirements. For more information about
// autoscaling, see Autoscaling (https://docs.cloud.oracle.com/Content/Compute/Tasks/autoscalinginstancepools.htm). For information about the
// Compute service, see Overview of the Compute Service (https://docs.cloud.oracle.com/Content/Compute/Concepts/computeoverview.htm).
// **Note:** Autoscaling is not available in US Government Cloud tenancies. For more information, see
// Oracle Cloud Infrastructure US Government Cloud (https://docs.cloud.oracle.com/Content/General/Concepts/govoverview.htm).
//

package autoscaling

import (
	"github.com/oracle/oci-go-sdk/v46/common"
)

// Capacity Capacity limits for the instance pool.
type Capacity struct {

	// For a threshold-based autoscaling policy, this value is the maximum number of instances the instance pool is allowed
	// to increase to (scale out).
	// For a schedule-based autoscaling policy, this value is not used.
	Max *int `mandatory:"false" json:"max"`

	// For a threshold-based autoscaling policy, this value is the minimum number of instances the instance pool is allowed
	// to decrease to (scale in).
	// For a schedule-based autoscaling policy, this value is not used.
	Min *int `mandatory:"false" json:"min"`

	// For a threshold-based autoscaling policy, this value is the initial number of instances to launch in the instance pool
	// immediately after autoscaling is enabled. After autoscaling retrieves performance metrics, the number of
	// instances is automatically adjusted from this initial number to a number that is based on the limits that
	// you set.
	// For a schedule-based autoscaling policy, this value is the target pool size to scale to when executing the schedule
	// that's defined in the autoscaling policy.
	Initial *int `mandatory:"false" json:"initial"`
}

func (m Capacity) String() string {
	return common.PointerString(m)
}
