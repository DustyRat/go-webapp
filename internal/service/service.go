package service

import (
	"encoding/json"
	"net/http"
	"text/template"
	"time"

	"github.com/dustyrat/go-webapp/internal/controller"
	"github.com/dustyrat/go-webapp/pkg/model"

	router "github.com/dustyrat/go-metrics/router/mux"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog/log"
)

// BuildInfo ...
type BuildInfo struct {
	Start     time.Time `json:"-"`
	Uptime    string    `json:"uptime,omitempty"`
	BuildDate string    `json:"build_date,omitempty"`
	BuildHost string    `json:"build_host,omitempty"`
	GitURL    string    `json:"git_url,omitempty"`
	Branch    string    `json:"branch,omitempty"`
	SHA       string    `json:"sha,omitempty"`
	Version   string    `json:"version,omitempty"`
	Debug     bool      `json:"debug"`
}

// AddHandlers add system endpoints used by k8s and prometheus gather readiness, health, metrics, ect...
func AddHandlers(r *router.Router, buildinfo *BuildInfo, ctrl *controller.Controller, debug bool) {
	r.Handle("/info", info(buildinfo)).Methods(http.MethodGet, http.MethodHead)
	r.Handle("/ready", ready(ctrl)).Methods(http.MethodGet, http.MethodHead)
	r.Handle("/health", health()).Methods(http.MethodGet, http.MethodHead)
	r.Handle("/metrics", promhttp.Handler())

	// fs := http.FileServer(http.Dir("./swagger/"))
	// r.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", fs))
	// http.Handle("/swagger/", r)

	if debug {
		log.Warn().Msg("pprof enabled")
		r.PathPrefix("/debug/pprof").Handler(http.DefaultServeMux)
		go func() {
			log.Error().Err(http.ListenAndServe("localhost:6060", nil)).Send()
		}()
	}
}

func info(buildinfo *BuildInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buildinfo.Uptime = time.Since(buildinfo.Start).String()
		RespondWithJSON(w, http.StatusOK, buildinfo)
	}
}

func ready(ctrl *controller.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		err := ctrl.Ready()
		if err != nil {
			Respond(w, http.StatusServiceUnavailable, []byte(http.StatusText(http.StatusServiceUnavailable)))
		} else {
			Respond(w, http.StatusOK, []byte(http.StatusText(http.StatusOK)))
		}
	}
}

func health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		Respond(w, http.StatusOK, []byte(http.StatusText(http.StatusOK)))
	}
}

// RespondWithError ...
func RespondWithError(w http.ResponseWriter, code int, err error) {
	RespondWithJSON(w, code, model.ErrorResponse{Err: err.Error()})
}

// RespondWithErrors ...
func RespondWithErrors(w http.ResponseWriter, code int, errs []error) {
	RespondWithJSON(w, code, model.ErrorResponse{Errors: errs})
}

// RespondWithJSON ...
func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	body, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	Respond(w, code, body)
}

// Respond ...
func Respond(w http.ResponseWriter, code int, body []byte) {
	w.WriteHeader(code)
	if _, err := w.Write(body); err != nil {
		log.Error().Err(err).Send()
	}
}

func Build(files []string) (*template.Template, error) {
	templates, err := template.ParseFiles(files...)
	if err != nil {
		return nil, err
	}
	return templates, nil
}

func Render(w http.ResponseWriter, templates *template.Template) {
	if err := templates.Execute(w, nil); err != nil {
		log.Error().Err(err).Send()
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
