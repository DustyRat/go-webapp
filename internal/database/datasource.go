package db

import (
	"github.com/DustyRat/go-webapp/internal/config"
	"github.com/DustyRat/go-webapp/internal/database/mgo"
)

// Initialize creates a new Datasource object and populates it with tested connections to sql and mongo databases.
func Initialize(conf config.Config) (*mgo.Mongo, error) {
	mongo, err := mgo.Connect(conf.Mongo.Database, conf.Mongo.URL, conf.Collections)
	return mongo, err
}
