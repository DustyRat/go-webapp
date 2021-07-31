package model

import (
	"time"
)

// swagger:model Audit
type Audit struct {
	// Created By
	// ReadOnly: true
	CreatedBy User `json:"createdBy" bson:"-"`

	// Created Timestamp
	// ReadOnly: true
	// Example: 2006-01-02T15:04:05.999Z
	CreatedTs *time.Time `json:"createdTs,omitempty" bson:"-"`

	// Updated By
	// ReadOnly: true
	UpdatedBy User `json:"updatedBy" bson:"-"`

	// Updated Timestamp
	// ReadOnly: true
	// Example: 2006-01-02T15:04:05.999Z
	UpdatedTs *time.Time `json:"updatedTs,omitempty" bson:"-"`

	// Version
	// ReadOnly: true
	// Example: 1
	Version uint `json:"version" bson:"-"`
}

// swagger:model User
type User struct {
	// First Name
	// ReadOnly: true
	// Example: John
	FirstName string `json:"firstName" bson:"-"`

	// Last Name
	// ReadOnly: true
	// Example: Doe
	LastName string `json:"lastName" bson:"-"`

	// Username
	// ReadOnly: true
	// Example: john.doe
	Username string `json:"username" bson:"-"`
}
