// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// Pet pet
//
// swagger:model Pet
type Pet struct {
	Animal

	// name
	Name string `json:"name,omitempty"`
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (m *Pet) UnmarshalJSON(raw []byte) error {
	// AO0
	var aO0 Animal
	if err := swag.ReadJSON(raw, &aO0); err != nil {
		return err
	}
	m.Animal = aO0

	// AO1
	var dataAO1 struct {
		Name string `json:"name,omitempty"`
	}
	if err := swag.ReadJSON(raw, &dataAO1); err != nil {
		return err
	}

	m.Name = dataAO1.Name

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (m Pet) MarshalJSON() ([]byte, error) {
	_parts := make([][]byte, 0, 2)

	aO0, err := swag.WriteJSON(m.Animal)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO0)
	var dataAO1 struct {
		Name string `json:"name,omitempty"`
	}

	dataAO1.Name = m.Name

	jsonDataAO1, errAO1 := swag.WriteJSON(dataAO1)
	if errAO1 != nil {
		return nil, errAO1
	}
	_parts = append(_parts, jsonDataAO1)
	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this pet
func (m *Pet) Validate(formats strfmt.Registry) error {
	var res []error

	// validation for a type composition with Animal
	if err := m.Animal.Validate(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *Pet) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Pet) UnmarshalBinary(b []byte) error {
	var res Pet
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}