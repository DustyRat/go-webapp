package handler

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/DustyRat/go-webapp/internal/controller"
	"github.com/DustyRat/go-webapp/internal/middleware"
	"github.com/DustyRat/go-webapp/internal/model"
	"github.com/DustyRat/go-webapp/internal/options"
	"github.com/DustyRat/go-webapp/internal/service"
	dto "github.com/DustyRat/go-webapp/pkg/model"

	router "github.com/DustyRat/go-metrics/router/mux"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
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

func rbac(w http.ResponseWriter, r *http.Request) (middleware.User, error) {
	// TODO: RBAC implementation
	return middleware.User{SAMAccountName: "ANONYMOUS"}, nil
}

// AddHandlers adds handlers.
func AddHandlers(r *router.Router, buildinfo *service.BuildInfo, ctrl *controller.Controller, debug bool) {
	r.HandleWithMetrics("/", middleware.Logger(middleware.RBAC(rbac, insert(ctrl)))).Methods(http.MethodPost)
	r.HandleWithMetrics("/", middleware.Logger(middleware.RBAC(rbac, find(ctrl)))).Methods(http.MethodGet)
	r.HandleWithMetrics("/{id}", middleware.Logger(middleware.RBAC(rbac, get(ctrl)))).Methods(http.MethodGet)
	r.HandleWithMetrics("/{id}", middleware.Logger(middleware.RBAC(rbac, update(ctrl)))).Methods(http.MethodPut)
	r.HandleWithMetrics("/{id}", middleware.Logger(middleware.RBAC(rbac, delete(ctrl)))).Methods(http.MethodDelete)

	// create basic endpoints used for dev ops and prod support
	r.Handle("/info", info(buildinfo)).Methods(http.MethodGet, http.MethodHead)
	r.Handle("/ready", ready(ctrl)).Methods(http.MethodGet, http.MethodHead)
	r.Handle("/health", health()).Methods(http.MethodGet, http.MethodHead)
	r.Handle("/metrics", promhttp.Handler())

	// fs := http.FileServer(http.Dir("./swagger/"))
	// r.PathPrefix("/swagger/").Handler(http.StripPrefix("/swagger/", fs))
	// http.Handle("/swagger/", r)

	// if debug {
	// 	log.Warn().Msg("pprof enabled")
	// 	r.PathPrefix("/debug/pprof").Handler(http.DefaultServeMux)
	// 	go func() {
	// 		log.Error().Err(http.ListenAndServe("localhost:6060", nil)).Send()
	// 	}()
	// }
}

func ready(ctrl *controller.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		err := ctrl.Ready()
		if err != nil {
			service.Respond(w, http.StatusServiceUnavailable, []byte(http.StatusText(http.StatusServiceUnavailable)))
		} else {
			service.Respond(w, http.StatusOK, []byte(http.StatusText(http.StatusOK)))
		}
	}
}

func health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		service.Respond(w, http.StatusOK, []byte(http.StatusText(http.StatusOK)))
	}
}

func info(buildinfo *service.BuildInfo) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		buildinfo.Uptime = time.Now().Sub(buildinfo.Start).String()
		service.RespondWithJSON(w, http.StatusOK, buildinfo)
	}
}

// swagger:route POST / Model insert
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Security:
//       Bearer:
//
//     Parameters:
//     + name: request
//       in: body
//       type: Model
//
//     Responses:
//       default: ErrorResponse
//       201: CreatedResponse Created
//       400: ErrorResponse Bad Request
//       401: UnauthorizedResponse Unauthorized
//       403: ForbiddenResponse Forbidden
//       422: ErrorResponse Unprocessable Entity
//       500: ErrorResponse Internal Server Error
//       501: ErrorResponse Not Implemented
func insert(ctrl *controller.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		usr := middleware.GetUser(r)

		vars := mux.Vars(r)
		hex := vars["id"]

		start := time.Now()
		var err error
		e := log.Info()
		defer func(e *zerolog.Event, start time.Time) {
			if err != nil {
				e = log.Error().Stack().Err(err)
			}
			e.Str("handler", "insert").Str("id", hex).AnErr("context", ctx.Err()).Int64("resp_time", time.Now().Sub(start).Milliseconds()).Send()
		}(e, start)

		var document dto.Model
		err = json.NewDecoder(r.Body).Decode(&document)
		if err != nil {
			service.RespondWithError(w, http.StatusBadRequest, err)
			return
		}

		result, err := ctrl.Insert(ctx, model.User{FirstName: usr.GivenName, LastName: usr.Sn, Username: usr.SAMAccountName}, document)
		if ctx.Err() != nil { // Check for canceled/timedout requests before setting response
			return
		}

		if err != nil {
			service.RespondWithError(w, http.StatusUnprocessableEntity, err)
			return
		}

		if result.InsertedID != nil {
			if id, ok := result.InsertedID.(primitive.ObjectID); ok {
				w.Header().Add("Location", fmt.Sprintf("/%s", id.Hex()))
			}
		}
		service.RespondWithJSON(w, http.StatusCreated, dto.CreatedResponse{
			ID: result.InsertedID,
		})
	}
}

// swagger:route GET / Model find
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Security:
//       Bearer:
//
//     Responses:
//       default: ErrorResponse
//       200: List OK
//       400: ErrorResponse Bad Request
//       401: UnauthorizedResponse Unauthorized
//       500: ErrorResponse Internal Server Error
//       501: ErrorResponse Not Implemented
func find(ctrl *controller.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// user := middleware.GetUser(r)

		query := r.URL.Query()

		page := query.Get("page")
		limit := query.Get("count")
		sortBy := query.Get("sortBy")
		sortOrder := query.Get("sortOrder")

		opts := options.New(map[string]string{sortBy: sortOrder}, page, limit)
		start := time.Now()
		var err error
		e := log.Info()
		defer func(e *zerolog.Event, start time.Time) {
			if err != nil {
				e = log.Error().Stack().Err(err)
			}
			e.Str("handler", "find").AnErr("context", ctx.Err()).Str("query", r.URL.RawQuery).Int64("resp_time", time.Now().Sub(start).Milliseconds()).Send()
		}(e, start)

		documents, count, more, errs, warnings, err := ctrl.Find(ctx, query, opts)
		if ctx.Err() != nil { // Check for canceled/timedout requests before setting response
			return
		}

		if err != nil {
			service.RespondWithError(w, http.StatusNotImplemented, err)
			return
		} else if len(errs) > 0 {
			service.RespondWithErrors(w, http.StatusBadRequest, errs)
			return
		}

		service.RespondWithJSON(w, http.StatusOK, dto.List{
			Documents: documents,
			Page:      opts.Page(),
			Count:     count,
			Links:     dto.BuildPagination(r, opts, more),
			Warnings:  warnings,
		})
	}
}

// swagger:route GET /{id} Model get
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Security:
//       Bearer:
//
//     Parameters:
//     + in: path
//       name: id
//       description: MongoDB Object ID.
//       type: string
//       required: true
//
//     Responses:
//       default: ErrorResponse
//       200: Model OK
//       400: ErrorResponse Bad Request
//       401: UnauthorizedResponse Unauthorized
//       403: ForbiddenResponse Forbidden
//       404: ErrorResponse Not Found
//       500: ErrorResponse Internal Server Error
//       501: ErrorResponse Not Implemented
func get(ctrl *controller.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		// user := middleware.GetUser(r)

		vars := mux.Vars(r)
		id := vars["id"]

		var err error
		start := time.Now()
		e := log.Info()
		defer func(e *zerolog.Event, start time.Time) {
			if err != nil {
				e = log.Error().Stack().Err(err)
			}
			e.Str("handler", "get").Str("id", id).AnErr("context", ctx.Err()).Int64("resp_time", time.Now().Sub(start).Milliseconds()).Send()
		}(e, start)

		document, err := ctrl.Get(ctx, id)
		if ctx.Err() != nil { // Check for canceled/timedout requests before setting response
			return
		}

		if err != nil {
			// Note: hex.ErrLength and hex.InvalidByteError are no longer returned from primitive.ObjectIDFromHex
			// https://github.com/mongodb/mongo-go-driver/commit/3fd62610449ee969dc7069ee21f2d94172aef148
			switch err {
			case hex.ErrLength, primitive.ErrInvalidHex:
				service.RespondWithError(w, http.StatusBadRequest, err)
			case mongo.ErrNoDocuments:
				service.RespondWithError(w, http.StatusNotFound, err)
			default:
				switch t := err.(type) {
				case hex.InvalidByteError:
					service.RespondWithError(w, http.StatusBadRequest, t)
				default:
					service.RespondWithError(w, http.StatusNotImplemented, err)
				}
			}
			return
		}
		service.RespondWithJSON(w, http.StatusOK, document)
	}
}

// swagger:route PUT /{id} Model update
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Security:
//       Bearer:
//
//     Parameters:
//     + in: path
//       name: id
//       description: MongoDB Object ID.
//       type: string
//       required: true
//     + name: request
//       in: body
//       type: Model
//
//     Responses:
//       default: ErrorResponse
//       200: UpdatedResponse OK
//       400: ErrorResponse Bad Request
//       401: UnauthorizedResponse Unauthorized
//       403: ForbiddenResponse Forbidden
//       404: ErrorResponse Not Found
//       409: ConflictResponse Conflict
//       500: ErrorResponse Internal Server Error
//       501: ErrorResponse Not Implemented
func update(ctrl *controller.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		usr := middleware.GetUser(r)

		vars := mux.Vars(r)
		id := vars["id"]

		start := time.Now()
		var err error
		e := log.Info()
		defer func(e *zerolog.Event, start time.Time) {
			if err != nil {
				e = log.Error().Stack().Err(err)
			}
			e.Str("handler", "update").Str("id", id).AnErr("context", ctx.Err()).Int64("resp_time", time.Now().Sub(start).Milliseconds()).Send()
		}(e, start)

		version, err := ctrl.GetVersion(ctx, id)
		if ctx.Err() != nil { // Check for canceled/timedout requests before setting response
			return
		}

		if err != nil {
			// Note: hex.ErrLength and hex.InvalidByteError are no longer returned from primitive.ObjectIDFromHex
			// https://github.com/mongodb/mongo-go-driver/commit/3fd62610449ee969dc7069ee21f2d94172aef148
			switch err {
			case hex.ErrLength, primitive.ErrInvalidHex:
				service.RespondWithError(w, http.StatusBadRequest, err)
			case mongo.ErrNoDocuments:
				service.RespondWithError(w, http.StatusNotFound, err)
			default:
				switch t := err.(type) {
				case hex.InvalidByteError:
					service.RespondWithError(w, http.StatusBadRequest, t)
				default:
					service.RespondWithError(w, http.StatusNotImplemented, err)
				}
			}
			return
		}

		var document dto.Model
		err = json.NewDecoder(r.Body).Decode(&document)
		if err != nil {
			service.RespondWithError(w, http.StatusBadRequest, err)
			return
		}

		if document.Audit.Version != version {
			err = dto.ErrStaleUpdate
			service.RespondWithJSON(w, http.StatusConflict, dto.ConflictResponse{
				Message: dto.StaleUpdateMsg,
				Error:   err.Error(),
			})
			return
		}

		_, err = ctrl.Update(ctx, model.User{FirstName: usr.GivenName, LastName: usr.Sn, Username: usr.SAMAccountName}, id, document)
		if ctx.Err() != nil { // Check for canceled/timedout requests before setting response
			return
		}

		if err != nil {
			switch err {
			case mongo.ErrNoDocuments: // This actually never gets returned from a mongo update. Keeping it here anyway...
				service.RespondWithError(w, http.StatusNotFound, err)
			default:
				switch t := err.(type) {
				case mongo.WriteException:
					service.RespondWithError(w, http.StatusBadRequest, t)
				default:
					service.RespondWithError(w, http.StatusUnprocessableEntity, err)
				}
			}
			return
		}

		w.Header().Add("Location", fmt.Sprintf("/%s", id))
		service.RespondWithJSON(w, http.StatusOK, dto.UpdatedResponse{
			ID: id,
		})
	}
}

// swagger:route DELETE /{id} Model delete
//
//     Consumes:
//     - application/json
//
//     Produces:
//     - application/json
//
//     Schemes: http, https
//
//     Security:
//       Bearer:
//
//     Parameters:
//     + in: path
//       name: id
//       description: MongoDB Object ID.
//       type: string
//       required: true
//
//     Responses:
//       default: ErrorResponse
//       204: null No Content
//       401: UnauthorizedResponse Unauthorized
//       403: ForbiddenResponse Forbidden
//       404: ErrorResponse Not Found
//       409: ConflictResponse Conflict
//       500: ErrorResponse Internal Server Error
//       501: ErrorResponse Not Implemented
func delete(ctrl *controller.Controller) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		vars := mux.Vars(r)
		id := vars["id"]

		start := time.Now()
		var err error
		e := log.Info()
		defer func(e *zerolog.Event, start time.Time) {
			if err != nil {
				e = log.Error().Stack().Err(err)
			}
			e.Str("handler", "delete").Str("id", id).AnErr("context", ctx.Err()).Int64("resp_time", time.Now().Sub(start).Milliseconds()).Send()
		}(e, start)

		result, err := ctrl.Delete(ctx, id)
		if ctx.Err() != nil { // Check for canceled/timedout requests before setting response
			return
		}

		if err != nil {
			// Note: hex.ErrLength and hex.InvalidByteError are no longer returned from primitive.ObjectIDFromHex
			// https://github.com/mongodb/mongo-go-driver/commit/3fd62610449ee969dc7069ee21f2d94172aef148
			switch err {
			case hex.ErrLength, primitive.ErrInvalidHex:
				service.RespondWithError(w, http.StatusBadRequest, err)
			case mongo.ErrNoDocuments: // This actually never gets returned from a mongo update. Keeping it here anyway...
				service.RespondWithError(w, http.StatusNotFound, err)
			default:
				switch t := err.(type) {
				case hex.InvalidByteError, mongo.WriteException:
					service.RespondWithError(w, http.StatusBadRequest, t)
				default:
					service.RespondWithError(w, http.StatusUnprocessableEntity, err)
				}
			}
			return
		}

		if result.DeletedCount > 0 {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		service.RespondWithError(w, http.StatusNotFound, mongo.ErrNoDocuments)
	}
}
