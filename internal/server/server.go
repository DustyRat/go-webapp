package server

import (
	"fmt"
	"net/http"

	"github.com/dustyrat/go-webapp/internal/config"
	"github.com/dustyrat/go-webapp/internal/controller"
	"github.com/dustyrat/go-webapp/internal/rbac"
	"github.com/dustyrat/go-webapp/internal/service"
	"github.com/dustyrat/go-webapp/internal/service/file"
	"github.com/dustyrat/go-webapp/internal/service/handler"

	router "github.com/dustyrat/go-metrics/router/mux"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// Run configures and creates a new http.Server to be used for the application to listen on
func Run(info *service.BuildInfo) error {
	conf, err := config.GetConfig()
	if err != nil {
		return err
	}

	info.Debug = conf.Debug
	level, err := zerolog.ParseLevel(conf.LogLevel)
	if err != nil {
		level = zerolog.ErrorLevel
		log.Warn().Err(err).Msgf("unable to parse log level, logging level is set to %s", level.String())
	}
	zerolog.SetGlobalLevel(level)
	log.Logger = log.With().Str("application", conf.Name).Logger()

	ctrl, err := controller.New(conf)
	if err != nil {
		return errors.Wrap(err, "unable to create controller")
	}

	mux := mux.NewRouter()
	r := router.New(mux)
	service.AddHandlers(r, info, ctrl, conf.Debug)
	rbac := rbac.RBAC{} // TODO: RBAC initialization
	handler.AddHandlers(r, ctrl, rbac.Middleware)

	file.AddHandlers(r)

	srv := http.Server{
		Addr:    fmt.Sprintf(":%d", conf.Port),
		Handler: cors.Default().Handler(r),
	}

	log.Info().Msgf("Server running %v", srv.Addr)
	return srv.ListenAndServe()
}
