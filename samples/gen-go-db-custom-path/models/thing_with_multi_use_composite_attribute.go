// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/go-openapi/validate"
)

// ThingWithMultiUseCompositeAttribute thing with multi use composite attribute
//
// swagger:model ThingWithMultiUseCompositeAttribute
type ThingWithMultiUseCompositeAttribute struct {

	// four
	// Required: true
	Four *string `json:"four"`

	// one
	// Required: true
	One *string `json:"one"`

	// three
	// Required: true
	Three *string `json:"three"`

	// two
	// Required: true
	Two *string `json:"two"`
}

// Validate validates this thing with multi use composite attribute
func (m *ThingWithMultiUseCompositeAttribute) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateFour(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateOne(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateThree(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateTwo(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ThingWithMultiUseCompositeAttribute) validateFour(formats strfmt.Registry) error {

	if err := validate.Required("four", "body", m.Four); err != nil {
		return err
	}

	return nil
}

func (m *ThingWithMultiUseCompositeAttribute) validateOne(formats strfmt.Registry) error {

	if err := validate.Required("one", "body", m.One); err != nil {
		return err
	}

	return nil
}

func (m *ThingWithMultiUseCompositeAttribute) validateThree(formats strfmt.Registry) error {

	if err := validate.Required("three", "body", m.Three); err != nil {
		return err
	}

	return nil
}

func (m *ThingWithMultiUseCompositeAttribute) validateTwo(formats strfmt.Registry) error {

	if err := validate.Required("two", "body", m.Two); err != nil {
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ThingWithMultiUseCompositeAttribute) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ThingWithMultiUseCompositeAttribute) UnmarshalBinary(b []byte) error {
	var res ThingWithMultiUseCompositeAttribute
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
