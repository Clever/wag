// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// Branch branch
//
// swagger:model Branch
type Branch string

const (

	// BranchMaster captures enum value "master"
	BranchMaster Branch = "master"

	// BranchDEVBRANCH captures enum value "DEV_BRANCH"
	BranchDEVBRANCH Branch = "DEV_BRANCH"

	// BranchTest captures enum value "test"
	BranchTest Branch = "test"
)

// for schema
var branchEnum []interface{}

func init() {
	var res []Branch
	if err := json.Unmarshal([]byte(`["master","DEV_BRANCH","test"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		branchEnum = append(branchEnum, v)
	}
}

func (m Branch) validateBranchEnum(path, location string, value Branch) error {
	if err := validate.Enum(path, location, value, branchEnum); err != nil {
		return err
	}
	return nil
}

// Validate validates this branch
func (m Branch) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateBranchEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
