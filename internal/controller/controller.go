package controller

import (
	"context"
	"net/url"
	"time"

	"github.com/DustyRat/go-webapp/internal/config"
	db "github.com/DustyRat/go-webapp/internal/database"
	"github.com/DustyRat/go-webapp/internal/database/mgo"
	"github.com/DustyRat/go-webapp/internal/model"
	"github.com/DustyRat/go-webapp/internal/options"
	dtoModel "github.com/DustyRat/go-webapp/pkg/model"

	collection "github.com/DustyRat/go-metrics/db/mgo"
	"github.com/DustyRat/go-metrics/metrics"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	moptions "go.mongodb.org/mongo-driver/mongo/options"
)

// Controller ...
type Controller struct {
	Mongo      *mgo.Mongo
	collection *collection.Collection
}

// New Create a new Controller
func New(cfg config.Config) (*Controller, error) {
	mongo, err := db.Initialize(cfg)
	if err != nil {
		return nil, err
	}

	return &Controller{
		Mongo:      mongo,
		collection: mongo.GetCollection(mgo.Collection),
	}, nil
}

// Ready ...
func (c *Controller) Ready() error {
	if err := c.Mongo.Ping(); err != nil {
		return err
	}
	return nil
}

func (c *Controller) GetVersion(ctx context.Context, hex string) (uint, error) {
	id, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		return 0, err
	}

	version, err := mgo.GetVersion(ctx, c.collection, id)
	if err != nil {
		return 0, err
	}
	return version, nil
}

// Insert ...
func (c *Controller) Insert(ctx context.Context, user model.User, dto dtoModel.Document) (*mongo.InsertOneResult, error) {
	start := time.Now()
	defer metrics.ObserveCaller("controller", start)

	document := model.TransformFromDTO(dto)

	/** Read Only **/
	document.ID = nil

	now := time.Now()
	document.Audit.CreatedTs = &now
	document.Audit.UpdatedTs = &now
	document.Audit.CreatedBy = &user
	document.Audit.UpdatedBy = user
	/** Read Only **/

	result, err := mgo.Insert(ctx, c.collection, document)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Find ...
func (c *Controller) Find(ctx context.Context, query url.Values, opts options.Options) ([]dtoModel.Document, int, bool, []error, []error, error) {
	start := time.Now()
	defer metrics.ObserveCaller("controller", start)

	match, errs, warnings := mgo.ParseQuery(query)
	if len(errs) > 0 {
		return []dtoModel.Document{}, 0, false, errs, warnings, nil
	}

	o := moptions.Find().SetSkip(int64(opts.Skip())).SetLimit(int64(opts.Limit() + 1))
	sort := opts.Sort(mgo.GetSort)
	if sort != nil && len(sort) > 0 {
		o.SetSort(sort)
	}

	details, err := mgo.Find(ctx, c.collection, match, o)
	if err != nil {
		return []dtoModel.Document{}, 0, false, nil, nil, err
	}

	dtos := make([]dtoModel.Document, 0)
	for _, document := range details {
		dtos = append(dtos, model.TransformToDTO(document))
	}

	more := (len(dtos) > opts.Limit())
	size := opts.Limit()
	if size > len(dtos) {
		size = len(dtos)
	}
	return dtos[:size], len(dtos[:size]), more, nil, warnings, nil
}

// Get ...
func (c *Controller) Get(ctx context.Context, hex string) (*dtoModel.Document, error) {
	start := time.Now()
	defer metrics.ObserveCaller("controller", start)
	id, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		return nil, err
	}

	document, err := mgo.Get(ctx, c.collection, id)
	if err != nil {
		return nil, err
	}
	dto := model.TransformToDTO(document)
	return &dto, nil
}

// Update ...
func (c *Controller) Update(ctx context.Context, user model.User, hex string, dto dtoModel.Document) (*mongo.UpdateResult, error) {
	start := time.Now()
	defer metrics.ObserveCaller("controller", start)
	id, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		return nil, err
	}

	document := model.TransformFromDTO(dto)

	/** Read Only **/
	document.Audit.UpdatedBy = user
	/** Read Only **/

	result, err := mgo.Update(ctx, c.collection, id, document)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Delete ...
func (c *Controller) Delete(ctx context.Context, hex string) (*mongo.DeleteResult, error) {
	start := time.Now()
	defer metrics.ObserveCaller("controller", start)
	id, err := primitive.ObjectIDFromHex(hex)
	if err != nil {
		return nil, err
	}

	result, err := mgo.Delete(ctx, c.collection, id)
	if err != nil {
		return nil, err
	}
	return result, nil
}
