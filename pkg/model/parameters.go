package model

import (
	"time"
)

// nolint
// swagger:parameters find_documents
type find struct {
	// in: query
	// Page Number
	// Example: 0
	Page int `json:"page"`
	// Page Count
	// Example: 10
	Count int `json:"count"`
	// Field to sort by (createdTs, updatedTs)
	// Example: updatedTs
	SortBy string `json:"sortBy"`
	// Sort Order (asc or desc)
	// Example: desc
	SortOrder string `json:"sortOrder"`

	// Created on Date
	// Example: 2006-01-02
	CreatedOn time.Time `json:"createdOn"`
	// Created After Date
	// Example: 2006-01-02
	CreatedAfter time.Time `json:"createdAfter"`
	// Created Before Date
	// Example: 2006-01-02
	CreatedBefore time.Time `json:"createdBefore"`

	// Updated on Date
	// Example: 2006-01-02
	UpdatedOn time.Time `json:"updatedOn"`
	// Updated After Date
	// Example: 2006-01-02
	UpdatedAfter time.Time `json:"updatedAfter"`
	// Updated Before Date
	// Example: 2006-01-02
	UpdatedBefore time.Time `json:"updatedBefore"`
}
