// Copyright (c) 2016, 2018, 2022, Oracle and/or its affiliates.  All rights reserved.
// This software is dual-licensed to you under the Universal Permissive License (UPL) 1.0 as shown at https://oss.oracle.com/licenses/upl or Apache License 2.0 as shown at http://www.apache.org/licenses/LICENSE-2.0. You may choose either license.
// Code generated. DO NOT EDIT.

// VisionService API
//
// A description of the VisionService API.
//

package aivision

import (
	"encoding/json"
	"fmt"
	"github.com/oracle/oci-go-sdk/v57/common"
	"strings"
)

// ObjectListInlineInputLocation A list of object locations in Object Storage.
type ObjectListInlineInputLocation struct {

	// List of ObjectLocations.
	ObjectLocations []ObjectLocation `mandatory:"true" json:"objectLocations"`
}

func (m ObjectListInlineInputLocation) String() string {
	return common.PointerString(m)
}

// ValidateEnumValue returns an error when providing an unsupported enum value
// This function is being called during constructing API request process
// Not recommended for calling this function directly
func (m ObjectListInlineInputLocation) ValidateEnumValue() (bool, error) {
	errMessage := []string{}

	if len(errMessage) > 0 {
		return true, fmt.Errorf(strings.Join(errMessage, "\n"))
	}
	return false, nil
}

// MarshalJSON marshals to json representation
func (m ObjectListInlineInputLocation) MarshalJSON() (buff []byte, e error) {
	type MarshalTypeObjectListInlineInputLocation ObjectListInlineInputLocation
	s := struct {
		DiscriminatorParam string `json:"sourceType"`
		MarshalTypeObjectListInlineInputLocation
	}{
		"OBJECT_LIST_INLINE_INPUT_LOCATION",
		(MarshalTypeObjectListInlineInputLocation)(m),
	}

	return json.Marshal(&s)
}
