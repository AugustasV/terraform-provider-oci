// Copyright (c) 2016, 2018, 2019, Oracle and/or its affiliates. All rights reserved.
// Code generated. DO NOT EDIT.

// Database Service API
//
// The API for the Database Service.
//

package database

import (
	"encoding/json"
	"github.com/oracle/oci-go-sdk/common"
)

// CreateDbHomeBase Details for creating a Database Home.
// **Warning:** Oracle recommends that you avoid using any confidential information when you supply string values using the API.
type CreateDbHomeBase interface {

	// The user-provided name of the Database Home.
	GetDisplayName() *string
}

type createdbhomebase struct {
	JsonData    []byte
	DisplayName *string `mandatory:"false" json:"displayName"`
	Source      string  `json:"source"`
}

// UnmarshalJSON unmarshals json
func (m *createdbhomebase) UnmarshalJSON(data []byte) error {
	m.JsonData = data
	type Unmarshalercreatedbhomebase createdbhomebase
	s := struct {
		Model Unmarshalercreatedbhomebase
	}{}
	err := json.Unmarshal(data, &s.Model)
	if err != nil {
		return err
	}
	m.DisplayName = s.Model.DisplayName
	m.Source = s.Model.Source

	return err
}

// UnmarshalPolymorphicJSON unmarshals polymorphic json
func (m *createdbhomebase) UnmarshalPolymorphicJSON(data []byte) (interface{}, error) {

	if data == nil || string(data) == "null" {
		return nil, nil
	}

	var err error
	switch m.Source {
	case "DB_BACKUP":
		mm := CreateDbHomeWithDbSystemIdFromBackupDetails{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	case "NONE":
		mm := CreateDbHomeWithDbSystemIdDetails{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	case "VM_CLUSTER_NEW":
		mm := CreateDbHomeWithVmClusterIdDetails{}
		err = json.Unmarshal(data, &mm)
		return mm, err
	default:
		return *m, nil
	}
}

//GetDisplayName returns DisplayName
func (m createdbhomebase) GetDisplayName() *string {
	return m.DisplayName
}

func (m createdbhomebase) String() string {
	return common.PointerString(m)
}

// CreateDbHomeBaseSourceEnum Enum with underlying type: string
type CreateDbHomeBaseSourceEnum string

// Set of constants representing the allowable values for CreateDbHomeBaseSourceEnum
const (
	CreateDbHomeBaseSourceNone         CreateDbHomeBaseSourceEnum = "NONE"
	CreateDbHomeBaseSourceDbBackup     CreateDbHomeBaseSourceEnum = "DB_BACKUP"
	CreateDbHomeBaseSourceVmClusterNew CreateDbHomeBaseSourceEnum = "VM_CLUSTER_NEW"
)

var mappingCreateDbHomeBaseSource = map[string]CreateDbHomeBaseSourceEnum{
	"NONE":           CreateDbHomeBaseSourceNone,
	"DB_BACKUP":      CreateDbHomeBaseSourceDbBackup,
	"VM_CLUSTER_NEW": CreateDbHomeBaseSourceVmClusterNew,
}

// GetCreateDbHomeBaseSourceEnumValues Enumerates the set of values for CreateDbHomeBaseSourceEnum
func GetCreateDbHomeBaseSourceEnumValues() []CreateDbHomeBaseSourceEnum {
	values := make([]CreateDbHomeBaseSourceEnum, 0)
	for _, v := range mappingCreateDbHomeBaseSource {
		values = append(values, v)
	}
	return values
}
