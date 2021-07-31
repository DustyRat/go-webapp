package service

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/DustyRat/go-webapp/pkg/model"
)

// BuildInfo ...
type BuildInfo struct {
	Start     time.Time `json:"-"`
	Uptime    string    `json:"uptime,omitempty"`
	Version   string    `json:"version,omitempty"`
	BuildDate string    `json:"build_date,omitempty"`
	BuildHost string    `json:"build_host,omitempty"`
	GitURL    string    `json:"git_url,omitempty"`
	Branch    string    `json:"branch,omitempty"`
	Debug     bool      `json:"debug"`
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
	w.Write(body)
}
