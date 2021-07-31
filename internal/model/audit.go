package model

import (
	"time"
)

// Audit database model
type Audit struct {
	CreatedBy *User      `json:"-" bson:"createdBy,omitempty"`
	CreatedTs *time.Time `json:"-" bson:"createdTs,omitempty"`
	UpdatedBy User       `json:"-" bson:"updatedBy"`
	UpdatedTs *time.Time `json:"-" bson:"updatedTs,omitempty"`
	Version   uint       `json:"-" bson:"version,omitempty"`
}

// User database model
type User struct {
	FirstName string `json:"-" bson:"firstName,omitempty"`
	LastName  string `json:"-" bson:"lastName,omitempty"`
	Username  string `json:"-" bson:"username,omitempty"`
}
