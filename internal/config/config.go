package config

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"

	"github.com/rs/zerolog/log"
)

// Config app configuration
type Config struct {
	Name  string `json:"name"`
	Port  int    `json:"port"`
	Debug bool   `json:"debug"`

	LogLevel string `json:"log-level"`

	Mongo       Mongo             `json:"mongo"`
	Collections map[string]string `json:"collections"`
}

type Mongo struct {
	Database string `json:"database"`
	URL      string `json:"url"`
}

// GetConfig mappings for app
func GetConfig() (Config, error) {
	conf := Config{}
	if err := read("config/config.json", &conf); err != nil {
		return conf, err
	}
	return override(conf), nil
}

// override datasource url values with envorinment variable values if found
func override(conf Config) Config {
	if val, ok := os.LookupEnv("MONGO_URL"); ok {
		conf.Mongo.URL = val
	}
	return conf
}

func read(path string, conf interface{}) error {
	if reflect.ValueOf(conf).Kind() != reflect.Ptr {
		return errors.New("interface is not a pointer")
	}

	file, err := os.Open(path)
	if err != nil {
		log.Error().Stack().Caller().Err(err).Msgf("unable to open file: %s", path)
		return err
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(conf)
	if err != nil {
		log.Error().Stack().Caller().Err(err).Msgf("unable to parse file: %s", path)
		return err
	}
	return nil
}
