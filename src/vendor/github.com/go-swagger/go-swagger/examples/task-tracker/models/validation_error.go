package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"

	"github.com/go-openapi/errors"
)

// ValidationError validation error
// swagger:model ValidationError
type ValidationError struct {
	Error

	// an optional field name to which this validation error applies
	Field string `json:"field,omitempty"`
}

// UnmarshalJSON unmarshals this object from a JSON structure
func (m *ValidationError) UnmarshalJSON(raw []byte) error {

	var aO0 Error
	if err := swag.ReadJSON(raw, &aO0); err != nil {
		return err
	}
	m.Error = aO0

	return nil
}

// MarshalJSON marshals this object to a JSON structure
func (m ValidationError) MarshalJSON() ([]byte, error) {
	var _parts [][]byte

	aO0, err := swag.WriteJSON(m.Error)
	if err != nil {
		return nil, err
	}
	_parts = append(_parts, aO0)

	return swag.ConcatJSON(_parts...), nil
}

// Validate validates this validation error
func (m *ValidationError) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.Error.Validate(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}