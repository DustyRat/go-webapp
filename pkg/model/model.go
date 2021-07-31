package model

import (
	"encoding/json"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// swagger:model ErrorResponse
type ErrorResponse struct {
	Err    string `json:"error,omitempty"`
	Errors Errors `json:"errors,omitempty"`
}

// Errors holds multiple errors
type Errors []error

// MarshalJSON ...
func (e Errors) MarshalJSON() ([]byte, error) {
	errs := make([]string, 0)
	for _, err := range e {
		errs = append(errs, err.Error())
	}
	return json.Marshal(&errs)
}

// UnmarshalJSON ...
func (e *Errors) UnmarshalJSON(data []byte) error {
	strs := make([]string, 0)
	if err := json.Unmarshal(data, &strs); err != nil {
		return err
	}

	errs := make([]error, 0)
	for _, str := range strs {
		errs = append(errs, errors.New(str))
	}
	*e = Errors(errs)
	return nil
}

// Contains checks if any errors passed into the method matches any errors within Errors at any error level
func (e Errors) Contains(err ...error) bool {
	for i := range e {
		for j := range err {
			if e[i] == err[j] || errors.Cause(e[i]) == err[j] {
				return true
			}
		}
	}
	return false
}

// swagger:model Model
type Model struct {
	// ID
	// ReadOnly: true
	// Example: 000000000000000000000000
	// swagger:strfmt bsonobjectid
	ID *primitive.ObjectID `json:"id,omitempty" bson:"-"`

	Audit Audit `json:",inline" bson:"_"`

	// TODO: Add fields
}
