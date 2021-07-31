package options

import (
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
)

// Options query parameters
type Options struct {
	sort  map[string]string
	page  string
	limit string
}

// New return new Options from query parts
func New(sort map[string]string, page, limit string) Options {
	return Options{
		sort:  sort,
		page:  page,
		limit: limit,
	}
}

type convert func(map[string]string) bson.M

// Sort get sort mongodb operation
func (o *Options) Sort(fn convert) bson.M {
	return fn(o.sort)
}

// Page get calculated page
func (o *Options) Page() int {
	if o.page == "" {
		return 1
	}

	page, err := strconv.Atoi(o.page)
	if err != nil || page < 1 {
		return 1
	}
	return page
}

// Limit get calculated limit
func (o *Options) Limit() int {
	if o.limit == "" {
		return 25
	}

	limit, err := strconv.Atoi(o.limit)
	if err != nil || limit < 1 || limit > 50 {
		return 25
	}
	return limit
}

// Skip calculate mongodb skip value
func (o *Options) Skip() int {
	page, limit := o.Page(), o.Limit()
	return (page - 1) * limit
}

// CreatePipeline create mongodb aggregate pipeline
func CreatePipeline(match bson.M, skip, limit int, sort bson.M) []bson.M {
	pipeline := make([]bson.M, 0)
	pipeline = append(pipeline, bson.M{"$match": match})
	if sort != nil && len(sort) > 0 {
		pipeline = append(pipeline, bson.M{"$sort": sort})
	}
	pipeline = append(pipeline, bson.M{"$skip": skip})
	pipeline = append(pipeline, bson.M{"$limit": limit})
	return pipeline
}
