package controller

import (
	"context"
	"net/url"
	"os"
	"reflect"
	"testing"

	"github.com/DustyRat/go-webapp/internal/config"
	"github.com/DustyRat/go-webapp/internal/database/mgo"
	"github.com/DustyRat/go-webapp/internal/model"
	"github.com/DustyRat/go-webapp/internal/options"
	"github.com/DustyRat/go-webapp/internal/utils"
	dto "github.com/DustyRat/go-webapp/pkg/model"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/rs/zerolog"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func TestNew(t *testing.T) {
	type args struct {
		cfg config.Config
	}
	tests := []struct {
		name    string
		args    args
		want    *Controller
		wantErr bool
	}{
		// TODO: Add test cases..
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.args.cfg)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestController_Ready(t *testing.T) {
	mongodb, err := mgo.Connect(conf.Mongo.Database, conf.Mongo.URL, conf.Collections)
	if err != nil {
		t.Errorf("controller.go Ready() error = %v", err)
		return
	}
	defer mongodb.Disconnect()

	if err := mongodb.Ping(); err != nil {
		t.Skipf("controller.go Ready() error = %v", err)
		return
	}

	ctrl := &Controller{
		Mongo:      mongodb,
		collection: mongodb.GetCollection(mgo.Collection),
	}

	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases..
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ctrl.Ready(); (err != nil) != tt.wantErr {
				t.Errorf("controller.go Ready() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestController_Insert(t *testing.T) {
	mongodb, err := mgo.Connect(conf.Mongo.Database, conf.Mongo.URL, conf.Collections)
	if err != nil {
		t.Errorf("controller.go Insert() error = %v", err)
		return
	}
	defer mongodb.Disconnect()

	if err := mongodb.Ping(); err != nil {
		t.Skipf("controller.go Insert() error = %v", err)
		return
	}

	ctrl := &Controller{
		Mongo:      mongodb,
		collection: mongodb.GetCollection(mgo.Collection),
	}

	opts := cmp.Options{
		utils.EquateErrors(),
		cmpopts.IgnoreFields(dto.Model{}, "ID"),
		cmpopts.IgnoreFields(dto.Audit{}, "CreatedTs", "UpdatedTs"),
	}
	type args struct {
		user model.User
		dto  dto.Model
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
			result, err := ctrl.Insert(context.Background(), test.args.user, test.args.dto)
			got := want{
				Result: result,
				Err:    err,
			}
			if !cmp.Equal(test.want, got, opts) {
				t.Errorf("controller.go Insert() mismatch (-want +got):\n%s", cmp.Diff(test.want, got, opts))
			}
		})
	}
}

func TestController_Find(t *testing.T) {
	mongodb, err := mgo.Connect(conf.Mongo.Database, conf.Mongo.URL, conf.Collections)
	if err != nil {
		t.Errorf("controller.go Find() error = %v", err)
		return
	}
	defer mongodb.Disconnect()

	if err := mongodb.Ping(); err != nil {
		t.Skipf("controller.go Find() error = %v", err)
		return
	}

	ctrl := &Controller{
		Mongo:      mongodb,
		collection: mongodb.GetCollection(mgo.Collection),
	}

	opts := cmp.Options{
		utils.EquateErrors(),
	}
	type args struct {
		query url.Values
		opts  options.Options
	}
	type want struct {
		Documents []dto.Model
		Count     int
		More      bool
		Errs      []error
		Warnings  []error
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
			documents, count, more, errs, warnings, err := ctrl.Find(context.Background(), test.args.query, test.args.opts)
			got := want{
				Documents: documents,
				Count:     count,
				More:      more,
				Errs:      errs,
				Warnings:  warnings,
				Err:       err,
			}

			if !cmp.Equal(test.want, got, opts) {
				t.Errorf("controller.go Find() mismatch (-want +got):\n%s", cmp.Diff(test.want, got, opts))
			}
		})
	}
}

func TestController_Get(t *testing.T) {
	mongodb, err := mgo.Connect(conf.Mongo.Database, conf.Mongo.URL, conf.Collections)
	if err != nil {
		t.Errorf("controller.go Get() error = %v", err)
		return
	}
	defer mongodb.Disconnect()

	if err := mongodb.Ping(); err != nil {
		t.Skipf("controller.go Get() error = %v", err)
		return
	}

	ctrl := &Controller{
		Mongo:      mongodb,
		collection: mongodb.GetCollection(mgo.Collection),
	}

	opts := cmp.Options{
		utils.EquateErrors(),
	}

	type args struct {
		hex string
	}
	type want struct {
		Document *dto.Model
		Err      error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		// TODO: Add test cases.
		{
			name: "Invalid Hex",
			args: args{
				hex: "01",
			},
			want: want{
				Document: nil,
				Err:      primitive.ErrInvalidHex,
			},
		},
		{
			name: "Invalid Byte",
			args: args{
				hex: "ZXCVBNM<ASDFGHJKL:",
			},
			want: want{
				Document: nil,
				Err:      primitive.ErrInvalidHex,
			},
		},
		{
			name: "Invalid Hex Length",
			args: args{
				hex: "0",
			},
			want: want{
				Document: nil,
				Err:      primitive.ErrInvalidHex,
			},
		},
		{
			name: "No Documents",
			args: args{
				hex: "000000000000000000000000",
			},
			want: want{
				Document: nil,
				Err:      mongo.ErrNoDocuments,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			document, err := ctrl.Get(context.Background(), test.args.hex)
			got := want{
				Document: document,
				Err:      err,
			}
			if !cmp.Equal(test.want, got, opts) {
				t.Errorf("controller.go Get() mismatch (-want +got):\n%s", cmp.Diff(test.want, got, opts))
			}
		})
	}
}

func TestController_Update(t *testing.T) {
	mongodb, err := mgo.Connect(conf.Mongo.Database, conf.Mongo.URL, conf.Collections)
	if err != nil {
		t.Errorf("controller.go Update() error = %v", err)
		return
	}
	defer mongodb.Disconnect()

	if err := mongodb.Ping(); err != nil {
		t.Skipf("controller.go Update() error = %v", err)
		return
	}

	ctrl := &Controller{
		Mongo:      mongodb,
		collection: mongodb.GetCollection(mgo.Collection),
	}

	opts := cmp.Options{
		utils.EquateErrors(),
	}

	type args struct {
		user model.User
		hex  string
		dto  dto.Model
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
			name: "Invalid Hex",
			args: args{
				user: model.User{},
				hex:  "01",
				dto:  dto.Model{},
			},
			want: want{
				Result: nil,
				Err:    primitive.ErrInvalidHex,
			},
		},
		{
			name: "Invalid Byte",
			args: args{
				user: model.User{},
				hex:  "ZXCVBNM<ASDFGHJKL:",
				dto:  dto.Model{},
			},
			want: want{
				Result: nil,
				Err:    primitive.ErrInvalidHex,
			},
		},
		{
			name: "Invalid Hex Length",
			args: args{
				user: model.User{},
				hex:  "0",
				dto:  dto.Model{},
			},
			want: want{
				Result: nil,
				Err:    primitive.ErrInvalidHex,
			},
		},
		{
			name: "No Documents",
			args: args{
				user: model.User{
					Username:  "john.doe",
					LastName:  "Doe",
					FirstName: "John",
				},
				hex: "000000000000000000000000",
				dto: dto.Model{},
			},
			want: want{
				Result: &mongo.UpdateResult{},
				Err:    nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := ctrl.Update(context.Background(), test.args.user, test.args.hex, test.args.dto)
			got := want{
				Result: result,
				Err:    err,
			}
			if !cmp.Equal(test.want, got, opts) {
				t.Errorf("controller.go Update() mismatch (-want +got):\n%s", cmp.Diff(test.want, got, opts))
			}
		})
	}
}

func TestController_Delete(t *testing.T) {
	mongodb, err := mgo.Connect(conf.Mongo.Database, conf.Mongo.URL, conf.Collections)
	if err != nil {
		t.Errorf("controller.go Delete() error = %v", err)
		return
	}
	defer mongodb.Disconnect()

	if err := mongodb.Ping(); err != nil {
		t.Skipf("controller.go Delete() error = %v", err)
		return
	}

	ctrl := &Controller{
		Mongo:      mongodb,
		collection: mongodb.GetCollection(mgo.Collection),
	}

	opts := cmp.Options{
		utils.EquateErrors(),
	}

	type args struct {
		hex string
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
			name: "Invalid Hex",
			args: args{
				hex: "01",
			},
			want: want{
				Result: nil,
				Err:    primitive.ErrInvalidHex,
			},
		},
		{
			name: "Invalid Byte",
			args: args{
				hex: "ZXCVBNM<ASDFGHJKL:",
			},
			want: want{
				Result: nil,
				Err:    primitive.ErrInvalidHex,
			},
		},
		{
			name: "Invalid Hex Length",
			args: args{
				hex: "0",
			},
			want: want{
				Result: nil,
				Err:    primitive.ErrInvalidHex,
			},
		},
		{
			name: "No Documents",
			args: args{
				hex: "000000000000000000000000",
			},
			want: want{
				Result: &mongo.DeleteResult{},
				Err:    nil,
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := ctrl.Delete(context.Background(), test.args.hex)
			got := want{
				Result: result,
				Err:    err,
			}
			if !cmp.Equal(test.want, got, opts) {
				t.Errorf("controller.go Delete() mismatch (-want +got):\n%s", cmp.Diff(test.want, got, opts))
			}
		})
	}
}
