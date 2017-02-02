package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	strfmt "github.com/go-openapi/strfmt"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/validate"
)

// RankingCell ranking cell
// swagger:model RankingCell
type RankingCell struct {

	// score
	// Required: true
	Score *int64 `json:"score"`

	// The format is hh:mm:ss.
	// Required: true
	Time *string `json:"time"`

	// wrong answer
	// Required: true
	WrongAnswer *int64 `json:"wrongAnswer"`
}

// Validate validates this ranking cell
func (m *RankingCell) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateScore(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateTime(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if err := m.validateWrongAnswer(formats); err != nil {
		// prop
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *RankingCell) validateScore(formats strfmt.Registry) error {

	if err := validate.Required("score", "body", m.Score); err != nil {
		return err
	}

	return nil
}

func (m *RankingCell) validateTime(formats strfmt.Registry) error {

	if err := validate.Required("time", "body", m.Time); err != nil {
		return err
	}

	return nil
}

func (m *RankingCell) validateWrongAnswer(formats strfmt.Registry) error {

	if err := validate.Required("wrongAnswer", "body", m.WrongAnswer); err != nil {
		return err
	}

	return nil
}
