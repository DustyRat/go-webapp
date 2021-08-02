package mgo

import (
	"context"
	"os"
	"testing"

	"github.com/DustyRat/go-webapp/internal/model"
	"github.com/DustyRat/go-webapp/internal/utils"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func init() {
	if uri, ok := os.LookupEnv("MONGO_URL"); ok {
		conf.Mongo.URL = uri
	}
	zerolog.SetGlobalLevel(zerolog.NoLevel)
}

func TestInsert(t *testing.T) {
	mongodb, err := Connect(conf.Mongo.Database, conf.Mongo.URL, conf.Collections)
	if err != nil {
		t.Errorf("collection.go Insert() error = %v", err)
		return
	}
	defer mongodb.Disconnect()

	if err := mongodb.Ping(); err != nil {
		t.Skipf("collection.go Insert() error = %v", err)
		return
	}
	collection := mongodb.GetCollection(Collection)

	opts := cmp.Options{
		utils.EquateErrors(),
		cmpopts.IgnoreFields(mongo.InsertOneResult{}, "InsertedID"),
	}

	type args struct {
		document model.Document
	}
	type want struct {
		Result *mongo.InsertOneResult
		Err    error
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
			result, err := Insert(context.Background(), collection, test.args.document)
			got := want{
				Result: result,
				Err:    err,
			}
			if !cmp.Equal(test.want, got, opts) {
				t.Errorf("collection.go Insert() mismatch (-want +got):\n%s", cmp.Diff(test.want, got, opts))
			}
		})
	}
}

func TestFind(t *testing.T) {
	mongodb, err := Connect(conf.Mongo.Database, conf.Mongo.URL, conf.Collections)
	if err != nil {
		t.Errorf("collection.go Find() error = %v", err)
		return
	}
	defer mongodb.Disconnect()

	if err := mongodb.Ping(); err != nil {
		t.Skipf("collection.go Find() error = %v", err)
		return
	}
	collection := mongodb.GetCollection(Collection)

	opts := cmp.Options{
		utils.EquateErrors(),
	}

	type args struct {
		filter bson.M
		opts   []*options.FindOptions
	}
	type want struct {
		Documents []model.Document
		Err       error
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
			documents, err := Find(context.Background(), collection, test.args.filter, test.args.opts...)
			got := want{
				Documents: documents,
				Err:       err,
			}
			if !cmp.Equal(test.want, got, opts) {
				t.Errorf("collection.go Find() mismatch (-want +got):\n%s", cmp.Diff(test.want, got, opts))
			}
		})
	}
}

func TestGet(t *testing.T) {
	mongodb, err := Connect(conf.Mongo.Database, conf.Mongo.URL, conf.Collections)
	if err != nil {
		t.Errorf("collection.go Get() error = %v", err)
		return
	}
	defer mongodb.Disconnect()

	if err := mongodb.Ping(); err != nil {
		t.Skipf("collection.go Get() error = %v", err)
		return
	}
	collection := mongodb.GetCollection(Collection)

	opts := cmp.Options{
		utils.EquateErrors(),
	}

	type args struct {
		id primitive.ObjectID
	}
	type want struct {
		Document model.Document
		Err      error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "ErrNoDocuments",
			args: args{
				id: utils.PrimitiveObjectID("000000000000000000000000"),
			},
			want: want{
				Document: model.Document{},
				Err:      mongo.ErrNoDocuments,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			document, err := Get(context.Background(), collection, test.args.id)
			got := want{
				Document: document,
				Err:      err,
			}
			if !cmp.Equal(test.want, got, opts) {
				t.Errorf("collection.go Get() mismatch (-want +got):\n%s", cmp.Diff(test.want, got, opts))
			}
		})
	}
}

func TestUpdate(t *testing.T) {
	mongodb, err := Connect(conf.Mongo.Database, conf.Mongo.URL, conf.Collections)
	if err != nil {
		t.Errorf("collection.go Update() error = %v", err)
		return
	}
	defer mongodb.Disconnect()

	if err := mongodb.Ping(); err != nil {
		t.Skipf("collection.go Update() error = %v", err)
		return
	}
	collection := mongodb.GetCollection(Collection)

	opts := cmp.Options{
		utils.EquateErrors(),
	}

	type args struct {
		id       primitive.ObjectID
		document model.Document
	}
	type want struct {
		Result *mongo.UpdateResult
		Err    error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		// TODO: Add test cases.
		{
			name: "No Documents",
			args: args{
				id:       utils.PrimitiveObjectID("000000000000000000000000"),
				document: model.Document{},
			},
			want: want{
				Result: &mongo.UpdateResult{
					MatchedCount:  0,
					ModifiedCount: 0,
					UpsertedCount: 0,
					UpsertedID:    nil,
				},
				Err: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Update(context.Background(), collection, test.args.id, test.args.document)
			got := want{
				Result: result,
				Err:    err,
			}
			if !cmp.Equal(test.want, got, opts) {
				t.Errorf("collection.go Update() mismatch (-want +got):\n%s", cmp.Diff(test.want, got, opts))
			}
		})
	}
}

func TestDelete(t *testing.T) {
	mongodb, err := Connect(conf.Mongo.Database, conf.Mongo.URL, conf.Collections)
	if err != nil {
		t.Errorf("collection.go Delete() error = %v", err)
		return
	}
	defer mongodb.Disconnect()

	if err := mongodb.Ping(); err != nil {
		t.Skipf("collection.go Delete() error = %v", err)
		return
	}
	collection := mongodb.GetCollection(Collection)

	opts := cmp.Options{
		utils.EquateErrors(),
	}

	type args struct {
		id       primitive.ObjectID
		document model.Document
	}
	type want struct {
		Result *mongo.DeleteResult
		Err    error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		// TODO: Add test cases.
		{
			name: "No Documents",
			args: args{
				id:       utils.PrimitiveObjectID("000000000000000000000000"),
				document: model.Document{},
			},
			want: want{
				Result: &mongo.DeleteResult{
					DeletedCount: 0,
				},
				Err: nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Delete(context.Background(), collection, test.args.id)
			got := want{
				Result: result,
				Err:    err,
			}
			if !cmp.Equal(test.want, got, opts) {
				t.Errorf("collection.go Delete() mismatch (-want +got):\n%s", cmp.Diff(test.want, got, opts))
			}
		})
	}
}
