package mgo

import (
	"net/url"
	"strings"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GetSort parse map into a mongodb sort operation for authorizations
func GetSort(fields map[string]string) bson.M {
	sort := bson.M{}
	for key, direction := range fields {
		switch key {
		case "createdTs":
			key = "createdTs"
		case "updatedTs":
			key = "updatedTs"
		default:
			key = ""
		}

		if key != "" {
			switch direction {
			case "asc":
				sort[key] = 1
			case "desc":
				sort[key] = -1
			default:
				sort[key] = 1
			}
			break
		}
	}
	return sort
}

// ParseQuery parse query parameters into a mongodb filter for authorizations
func ParseQuery(query url.Values) (bson.M, []error, []error) {
	pipe := make([]bson.M, 0)
	errs := make([]error, 0)
	warnings := make([]error, 0)
	for key, values := range query {
		for i, value := range values {
			if value == "" {
				warnings = append(warnings, errors.Errorf("no value supplied for '%s='[%d]", key, i))
			}
		}

		switch key {
		case "page", "count", "sortBy", "sortOrder":
			if len(values) > 1 {
				warnings = append(warnings, errors.Errorf("param '%s' used %d times (limit 1);", key, len(values)))
			}
			continue // ignore
		case "id":
			in := make([]interface{}, 0)
			for _, value := range values {
				if value != "" {
					id, err := primitive.ObjectIDFromHex(value)
					if err != nil {
						errs = append(errs, errors.Errorf("invalid input: '%s=%s'; %s", key, value, err))
						continue
					}
					in = append(in, id)
				}
			}
			if len(in) > 0 {
				pipe = append(pipe, bson.M{"_id": bson.M{"$in": in}})
			}
		case "createdOn", "createdAfter", "createdBefore":
			or := make([]bson.M, 0)
			for _, value := range values {
				if value != "" {
					date, err := time.Parse("2006-01-02", value)
					if err != nil {
						errs = append(errs, errors.Errorf("invalid input: '%s=%s'; date must be in the format of 'YYYY-MM-DD' (%s=2006-01-02)", key, value, key))
						continue
					}
					if op := getDateRange(strings.TrimPrefix(key, "created"), "createdTs", date); op != nil {
						or = append(or, op)
					}
				}
			}
			if len(or) > 0 {
				pipe = append(pipe, bson.M{"$or": or})
			}
		case "updatedOn", "updatedAfter", "updatedBefore":
			or := make([]bson.M, 0)
			for _, value := range values {
				if value != "" {
					date, err := time.Parse("2006-01-02", value)
					if err != nil {
						errs = append(errs, errors.Errorf("invalid input: '%s=%s'; date must be in the format of 'YYYY-MM-DD' (%s=2006-01-02)", key, value, key))
						continue
					}
					if op := getDateRange(strings.TrimPrefix(key, "updated"), "updatedTs", date); op != nil {
						or = append(or, op)
					}
				}
			}
			if len(or) > 0 {
				pipe = append(pipe, bson.M{"$or": or})
			}
		case "":
			// log an error for any query parameters missing values
			for _, value := range values {
				errs = append(errs, errors.Errorf("no key supplied for '=%s'", value))
			}
		default:
			// keep track of any query parameters that don't match what we're looking for in our case statements (catch misspellings/invalid query params)
			for _, value := range values {
				errs = append(errs, errors.Errorf("'%s=%s' not processable: '%s' is not a valid query parameter;", key, value, key))
			}
		}
	}
	if len(pipe) > 0 {
		return bson.M{"$and": pipe}, errs, warnings
	}
	return bson.M{}, errs, warnings
}

// parseDateRange determine date range mongo filter
func parseDateRange(key string, date time.Time) bson.M {
	if strings.HasSuffix(key, "After") {
		return getDateRange("After", strings.TrimSuffix(key, "After"), date)
	} else if strings.HasSuffix(key, "Before") {
		return getDateRange("Before", strings.TrimSuffix(key, "Before"), date)
	} else if strings.HasSuffix(key, "On") {
		return getDateRange("On", strings.TrimSuffix(key, "On"), date)
	}
	return getDateRange("On", key, date)
}

// getDateRange get date range mongo filter
func getDateRange(op, field string, date time.Time) bson.M {
	switch op {
	case "After":
		return bson.M{field: bson.M{"$gte": date.Add(24 * time.Hour)}}
	case "Before":
		return bson.M{field: bson.M{"$lt": date}}
	case "On":
		return bson.M{field: bson.M{"$gte": date, "$lt": date.Add(24 * time.Hour)}}
	default:
		return bson.M{field: date}
	}
}
