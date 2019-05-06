// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/swag"
)

// Deployment deployment
// swagger:model Deployment
type Deployment struct {

	// application
	Application string `json:"application,omitempty"`

	// date
	Date strfmt.DateTime `json:"date,omitempty"`

	// environment
	Environment string `json:"environment,omitempty"`

	// version
	Version string `json:"version,omitempty"`
}

// Validate validates this deployment
func (m *Deployment) Validate(formats strfmt.Registry) error {
	var res []error

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// MarshalBinary interface implementation
func (m *Deployment) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *Deployment) UnmarshalBinary(b []byte) error {
	var res Deployment
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}