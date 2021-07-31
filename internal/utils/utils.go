package utils

import (
	"time"

	"github.com/google/go-cmp/cmp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ISODate helper function to convert a string into a *time.Time with "2006-01-02T15:04:05.000Z07:00" format
func ISODate(str string) *time.Time {
	t, err := time.Parse("2006-01-02T15:04:05.000Z07:00", str)
	if err != nil {
		return nil
	}
	return &t
}

// PrimitiveObjectID helper method to parse string into primitive.ObjectID
func PrimitiveObjectID(hex string) primitive.ObjectID {
	oid, _ := primitive.ObjectIDFromHex(hex)
	return oid
}

// PPrimitiveObjectID helper method to parse string into *primitive.ObjectID
func PPrimitiveObjectID(hex string) *primitive.ObjectID {
	id := PrimitiveObjectID(hex)
	return &id
}

// EquateErrors check to see if an error equals another error
// A modification on the provided cmpopts.EquateErrors()
// see: https://pkg.go.dev/github.com/google/go-cmp@v0.5.5/cmp/cmpopts#EquateErrors
func EquateErrors() cmp.Option {
	return cmp.FilterValues(areErrors, cmp.Comparer(compareErrors))
}

func areErrors(x, y interface{}) bool {
	_, ok1 := x.(error)
	_, ok2 := y.(error)
	return ok1 && ok2
}

func compareErrors(x, y interface{}) bool {
	xe := x.(error)
	ye := y.(error)
	return xe.Error() == ye.Error()
}
