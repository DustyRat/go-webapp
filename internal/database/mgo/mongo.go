package mgo

import (
	"context"
	"time"

	"github.com/DustyRat/go-metrics/db/mgo"
	"github.com/DustyRat/go-metrics/metrics"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config represents the basic configurations for the mongo.
type Config struct {
	Database string `json:"database"`
	URL      string `json:"url"`
}

// Mongo represents the session, collection, and database of a Mongo RTR object.
type Mongo struct {
	database    *mgo.Database
	collections map[string]*mgo.Collection
}

// Connect to mongodb instance
func Connect(dbname, rawurl string, collections map[string]string, opts ...*options.ClientOptions) (*Mongo, error) {
	mongo := Mongo{collections: make(map[string]*mgo.Collection)}
	database, err := mgo.Connect(dbname, rawurl, opts...)
	if err != nil {
		return nil, err
	}

	mongo.database = database
	for key, collection := range collections {
		mongo.collections[key] = mongo.database.Collection(collection)
	}
	return &mongo, nil
}

// Disconnect mongodb connection
func (db *Mongo) Disconnect() {
	db.database.Client().Disconnect(context.Background())
}

// GetCollection ...
func (db *Mongo) GetCollection(key string) *mgo.Collection {
	return db.collections[key]
}

// Ping ...
func (db *Mongo) Ping() error {
	return db.database.Ping()
}

// GetVersion ...
func GetVersion(ctx context.Context, collection *mgo.Collection, id primitive.ObjectID) (uint, error) {
	start := time.Now()
	defer metrics.ObserveCaller("mongo", start)

	var err error
	e := log.Debug().Str("method", "GetVersion")
	defer func(e *zerolog.Event, start time.Time) {
		e.Int64("resp_time", time.Since(start).Milliseconds()).Err(err).Send()
	}(e, start)

	document := struct {
		ID      *primitive.ObjectID `bson:"_id"`
		Version uint                `bson:"version"`
	}{}
	filter := bson.M{"_id": id}
	err = collection.FindOne(ctx, filter).Decode(&document)
	if err != nil {
		return 0, err
	}
	return document.Version, nil
}

// Count ...
func Count(ctx context.Context, collection *mgo.Collection, filter bson.M) (int64, error) {
	start := time.Now()
	defer metrics.ObserveCaller("mongo", start)

	var err error
	e := log.Debug().Str("method", "Count")
	defer func(e *zerolog.Event, start time.Time) {
		e.Int64("resp_time", time.Since(start).Milliseconds()).Err(err).Send()
	}(e, start)

	var count int64
	if len(filter) > 0 {
		count, err = collection.CountDocuments(ctx, filter)
	} else {
		count, err = collection.EstimatedDocumentCount(ctx)
	}
	if err != nil {
		return 0, err
	}
	return count, nil
}
