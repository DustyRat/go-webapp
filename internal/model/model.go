package model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Model database model
type Model struct {
	ID    *primitive.ObjectID `json:"-" bson:"_id,omitempty"`
	Audit Audit               `json:"-" bson:",inline"`

	// TODO: Add fields
}
