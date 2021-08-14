package mgo

import (
	"context"
	"time"

	"github.com/DustyRat/go-webapp/internal/model"

	"github.com/DustyRat/go-metrics/db/mgo"
	"github.com/DustyRat/go-metrics/metrics"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Collection Key
const Collection = "Document"

// Insert ...
func Insert(ctx context.Context, collection *mgo.Collection, document model.Document) (*mongo.InsertOneResult, error) {
	start := time.Now()
	defer metrics.ObserveCaller("mongo", start)

	var err error
	e := log.Debug().Str("method", "Insert")
	defer func(e *zerolog.Event, start time.Time) {
		e.Int64("resp_time", time.Since(start).Milliseconds()).Err(err).Send()
	}(e, start)

	/** Read Only **/
	document.ID = nil
	document.Audit.CreatedTs = &start
	document.Audit.UpdatedTs = &start
	document.Audit.Version = 1
	/** Read Only **/

	opts := options.InsertOne()
	result, err := collection.InsertOne(ctx, document, opts)
	if err != nil {
		return result, err
	}
	return result, nil
}

// Find ...
func Find(ctx context.Context, collection *mgo.Collection, filter bson.M, opts ...*options.FindOptions) ([]model.Document, error) {
	start := time.Now()
	defer metrics.ObserveCaller("mongo", start)

	var err error
	e := log.Debug().Str("method", "Find")
	defer func(e *zerolog.Event, start time.Time) {
		e.Int64("resp_time", time.Since(start).Milliseconds()).Err(err).Send()
	}(e, start)

	documents := make([]model.Document, 0)
	cur, err := collection.Find(ctx, filter, opts...)
	if err != nil {
		return documents, err
	}
	defer cur.Close(ctx)

	for cur.Next(ctx) {
		var document model.Document
		err := cur.Decode(&document)
		if err != nil {
			return documents, err
		}
		documents = append(documents, document)
	}
	return documents, err
}

// Get ...
func Get(ctx context.Context, collection *mgo.Collection, id primitive.ObjectID) (model.Document, error) {
	start := time.Now()
	defer metrics.ObserveCaller("mongo", start)

	var err error
	e := log.Debug().Str("method", "Get")
	defer func(e *zerolog.Event, start time.Time) {
		e.Int64("resp_time", time.Since(start).Milliseconds()).Err(err).Send()
	}(e, start)

	filter := bson.M{"_id": id}

	document := model.Document{}
	err = collection.FindOne(ctx, filter).Decode(&document)
	if err != nil {
		return document, err
	}
	return document, nil
}

// Update ...
func Update(ctx context.Context, collection *mgo.Collection, id primitive.ObjectID, document model.Document) (*mongo.UpdateResult, error) {
	start := time.Now()
	defer metrics.ObserveCaller("mongo", start)

	var err error
	e := log.Debug().Str("method", "Update")
	defer func(e *zerolog.Event, start time.Time) {
		e.Int64("resp_time", time.Since(start).Milliseconds()).Err(err).Send()
	}(e, start)

	filter := bson.M{
		"_id":     id,
		"version": document.Audit.Version, // prevent concurrent updates
	}

	/** Read Only **/
	document.ID = nil
	document.Audit.CreatedTs = nil
	document.Audit.CreatedBy = nil
	document.Audit.UpdatedTs = nil
	document.Audit.Version = 0
	/** Read Only **/

	update := bson.M{
		"$currentDate": bson.M{"updatedTs": bson.M{"$type": "date"}}, // set modified timestamp
		"$inc":         bson.M{"version": uint(1)},
		"$set":         &document,
	}
	opts := options.Update()
	result, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return result, err
	}
	return result, nil
}

// Delete ...
func Delete(ctx context.Context, collection *mgo.Collection, id primitive.ObjectID) (*mongo.DeleteResult, error) {
	start := time.Now()
	defer metrics.ObserveCaller("mongo", start)

	var err error
	e := log.Debug().Str("method", "Delete")
	defer func(e *zerolog.Event, start time.Time) {
		e.Int64("resp_time", time.Since(start).Milliseconds()).Err(err).Send()
	}(e, start)

	filter := bson.M{"_id": id}

	opts := options.Delete()
	result, err := collection.DeleteOne(ctx, filter, opts)
	if err != nil {
		return result, err
	}
	return result, nil
}
