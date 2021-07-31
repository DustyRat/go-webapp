package mgo

import (
	"context"
	"os"
	"reflect"
	"testing"

	"github.com/DustyRat/go-metrics/db/mgo"
	"github.com/DustyRat/go-webapp/internal/config"
	"github.com/DustyRat/go-webapp/internal/utils"
	"github.com/google/go-cmp/cmp"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	conf = config.Config{
		Mongo: config.Mongo{
			Database: "Example",
			URL:      "mongodb://localhost:27017",
		},
		Collections: map[string]string{
			"Model": "Model",
		},
	}
)

func init() {
	if uri, ok := os.LookupEnv("MONGO_URL"); ok {
		conf.Mongo.URL = uri
	}
	zerolog.SetGlobalLevel(zerolog.NoLevel)
}

func TestConnect(t *testing.T) {
	type args struct {
		dbname      string
		rawurl      string
		collections map[string]string
		opts        []*options.ClientOptions
	}
	tests := []struct {
		name    string
		args    args
		want    *Mongo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Connect(test.args.dbname, test.args.rawurl, test.args.collections, test.args.opts...)
			if (err != nil) != test.wantErr {
				t.Errorf("mongo.go Connect() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("mongo.go Connect() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestMongo_Disconnect(t *testing.T) {
	tests := []struct {
		name string
		db   *Mongo
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.db.Disconnect()
		})
	}
}

func TestMongo_GetCollection(t *testing.T) {
	type args struct {
		key string
	}
	tests := []struct {
		name string
		db   *Mongo
		args args
		want *mgo.Collection
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.db.GetCollection(test.args.key); !reflect.DeepEqual(got, test.want) {
				t.Errorf("mongo.go Mongo.GetCollection() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestMongo_Ping(t *testing.T) {
	tests := []struct {
		name    string
		db      *Mongo
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if err := test.db.Ping(); (err != nil) != test.wantErr {
				t.Errorf("mongo.go Mongo.Ping() error = %v, wantErr %v", err, test.wantErr)
			}
		})
	}
}

func TestGetVersion(t *testing.T) {
	mongodb, err := Connect(conf.Mongo.Database, conf.Mongo.URL, conf.Collections)
	if err != nil {
		t.Errorf("mongo.go GetVersion() error = %v", err)
		return
	}
	defer mongodb.Disconnect()

	if err := mongodb.Ping(); err != nil {
		t.Skipf("mongo.go GetVersion() error = %v", err)
		return
	}
	collection := mongodb.GetCollection(Collection)

	opts := cmp.Options{
		utils.EquateErrors(),
	}

	type args struct {
		collection *mgo.Collection
		id         primitive.ObjectID
	}
	type want struct {
		Version uint
		Err     error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			version, err := GetVersion(context.Background(), collection, test.args.id)
			got := want{
				Version: version,
				Err:     err,
			}
			if !cmp.Equal(test.want, got, opts) {
				t.Errorf("mongo.go GetVersion() mismatch (-want +got):\n%s", cmp.Diff(test.want, got, opts))
			}
		})
	}
}

func TestCount(t *testing.T) {
	mongodb, err := Connect(conf.Mongo.Database, conf.Mongo.URL, conf.Collections)
	if err != nil {
		t.Errorf("mongo.go Count() error = %v", err)
		return
	}
	defer mongodb.Disconnect()

	if err := mongodb.Ping(); err != nil {
		t.Skipf("mongo.go Count() error = %v", err)
		return
	}
	collection := mongodb.GetCollection(Collection)

	opts := cmp.Options{
		utils.EquateErrors(),
	}

	type args struct {
		collection *mgo.Collection
		filter     bson.M
	}
	type want struct {
		Count int64
		Err   error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			count, err := Count(context.Background(), collection, test.args.filter)
			got := want{
				Count: count,
				Err:   err,
			}
			if !cmp.Equal(test.want, got, opts) {
				t.Errorf("mongo.go Count() mismatch (-want +got):\n%s", cmp.Diff(test.want, got, opts))
			}
		})
	}
}
