package model

import (
	"time"
)

// swagger:parameters find_details
type find struct {
	// in: query
	// Page Number
	// Example: 0
	Page int `json:"page"`
	// Page Count
	// Example: 10
	Count int `json:"count"`
	// Authorization field to sort by (createdTs, updatedTs)
	// Example: updatedTs
	SortBy string `json:"sortBy"`
	// Sort Order (asc or desc)
	// Example: desc
	SortOrder string `json:"sortOrder"`

	// Mother's Member ID
	// Example: 000000000000000000000000
	// swagger:strfmt bsonobjectid
	Mother []string `json:"mother"`

	// Member ID
	// Example: 000000000000000000000000
	// swagger:strfmt bsonobjectid
	Member []string `json:"member"`

	// Delivery or Sick Baby Authorization ID
	// Example: 000000000000000000000000
	// swagger:strfmt bsonobjectid
	Authorization []string `json:"authorization"`

	// Delivery or Sick Baby Authorization Cloud/Classic Number
	// Example: [449M-Y6MM, IP123456789]
	AuthorizationNbr []string `json:"authorizationNbr"`

	// Delivery Authorization ID
	// Example: 000000000000000000000000
	// swagger:strfmt bsonobjectid
	Delivery []string `json:"delivery"`

	// Delivery Authorization Cloud/Classic Number
	// Example: [449M-Y6MM, IP123456789]
	DeliveryNbr []string `json:"deliveryNbr"`

	// Sick Baby Authorization ID
	// Example: 000000000000000000000000
	// swagger:strfmt bsonobjectid
	SickBaby []string `json:"sickBaby"`

	// Sick Baby Authorization Cloud/Classic Number
	// Example: [449M-Y6MM, IP123456789]
	SickBabyNbr []string `json:"sickBabyNbr"`

	// Business Unit Code
	// Example: 8
	BusinessUnitCd []int `json:"businessUnitCd"`
	// Business Line
	// Example: Medicaid
	BusinessLine []string `json:"businessLine"`

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
