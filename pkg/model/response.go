package model

import (
	"net/http"
	"net/url"
	"strconv"

	"github.com/DustyRat/go-webapp/internal/options"

	"github.com/pkg/errors"
)

const (
	ForbiddenMsg = "The request is forbidden"
	// StaleUpdateMsg stale document error message
	StaleUpdateMsg  = "The document submitted is stale, please refresh and try again"
	UnauthorizedMsg = "The request is not authorized"
)

var (
	// ErrStaleUpdate stale document error
	ErrStaleUpdate = errors.New("stale update")
)

// UnauthorizedResponse ...
// swagger:model UnauthorizedResponse
type UnauthorizedResponse struct {
	// Message
	// Example: malformed token
	Message string `json:"message,omitempty"`
	// Error
	// Example: The request is not authorized
	Error string `json:"error,omitempty"`
}

// ForbiddenResponse ...
// swagger:model ForbiddenResponse
type ForbiddenResponse struct {
	// Message
	// Example: user does not have access to [ BusinessUnit: '8', BusinessLine: 'MEDICARE' ]
	Message string `json:"message,omitempty"`
	// Error
	// Example: The request is forbidden
	Error string `json:"error,omitempty"`
}

// ConflictResponse ...
// swagger:model ConflictResponse
type ConflictResponse struct {
	// Message
	// Example: The document submitted is stale, please refresh and try again
	Message string `json:"message,omitempty"`
	// Error
	// Example: stale update
	Error string `json:"error,omitempty"`
}

// OKResponse ...
// swagger:model OKResponse
type OKResponse struct {
	// ID
	// ReadOnly: true
	// Example: 000000000000000000000000
	// swagger:strfmt bsonobjectid
	ID interface{} `json:"id,omitempty"`
}

// CreatedResponse ...
// swagger:model CreatedResponse
type CreatedResponse struct {
	// ID
	// ReadOnly: true
	// Example: 000000000000000000000000
	// swagger:strfmt bsonobjectid
	ID interface{} `json:"id,omitempty"`
}

// UpdatedResponse ...
// swagger:model UpdatedResponse
type UpdatedResponse struct {
	// ID
	// ReadOnly: true
	// Example: 000000000000000000000000
	// swagger:strfmt bsonobjectid
	ID interface{} `json:"id,omitempty"`
}

// List paginated list of Documents
// swagger:model List
type List struct {
	Documents []Document `json:"documents"`
	Page      int        `json:"page"`
	Count     int        `json:"count"`
	Links     []Link     `json:"links"`
	Warnings  Errors     `json:"warnings,omitempty"`
}

// Link reference pagination links
type Link struct {
	// Relitive
	// Example: next
	Rel string `json:"rel"`
	// Link
	// Example: /birth?page=2
	Href string `json:"href"`
}

// BuildPagination build pagination
func BuildPagination(r *http.Request, opts options.Options, more bool) []Link {
	links := make([]Link, 0)
	if opts.Page() != 1 {
		links = append(links, buildLink("first", 1, opts.Limit(), *r.URL))
	}

	if opts.Page()-1 > 0 {
		links = append(links, buildLink("prev", opts.Page()-1, opts.Limit(), *r.URL))
	}

	if more {
		links = append(links, buildLink("next", opts.Page()+1, opts.Limit(), *r.URL))
	}
	return links
}

func buildLink(rel string, page, count int, uri url.URL) Link {
	q := uri.Query()
	q.Set("page", strconv.Itoa(page))
	q.Set("count", strconv.Itoa(count))
	uri.RawQuery = q.Encode()
	return Link{
		Rel:  rel,
		Href: uri.String(),
	}
}
