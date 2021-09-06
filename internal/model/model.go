package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Document database model
type Document struct {
	ID    *primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Audit `json:"-" bson:",inline"`

	// TODO: Add fields
}

// func (d Document) MarshalJSON() ([]byte, error) {
// 	panic("do not marshal internal struct")
// }

// func (d *Document) UnmarshalJSON([]byte) error  {
// 	panic("do not unmarshal to internal struct")
// }
