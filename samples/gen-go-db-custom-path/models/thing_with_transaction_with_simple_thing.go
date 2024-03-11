// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ThingWithTransactionWithSimpleThing thing with transaction with simple thing
//
// swagger:model ThingWithTransactionWithSimpleThing
type ThingWithTransactionWithSimpleThing struct {

	// name
	Name string `json:"name,omitempty"`
}

// Validate validates this thing with transaction with simple thing
func (m *ThingWithTransactionWithSimpleThing) Validate(formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *ThingWithTransactionWithSimpleThing) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ThingWithTransactionWithSimpleThing) UnmarshalBinary(b []byte) error {
	var res ThingWithTransactionWithSimpleThing
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}