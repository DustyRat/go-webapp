package handler

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/DustyRat/go-webapp/internal/config"
	"github.com/DustyRat/go-webapp/internal/controller"
	"github.com/DustyRat/go-webapp/internal/middleware"
	"github.com/DustyRat/go-webapp/internal/utils"
	"github.com/DustyRat/go-webapp/pkg/model"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
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

func Test_insert(t *testing.T) {
	ctrl, _ := controller.New(conf)
	if err := ctrl.Ready(); err != nil {
		t.Skipf("handler.go insert() error = %v", err)
		return
	}
	defer ctrl.Mongo.Disconnect()

	handler := insert(ctrl)

	opts := cmp.Options{
		utils.EquateErrors(),
		cmpopts.IgnoreFields(model.CreatedResponse{}, "ID"),
		cmpopts.IgnoreFields(model.Model{}, "ID"),
		cmpopts.IgnoreFields(model.Audit{}, "CreatedTs", "UpdatedTs"),
	}

	type result struct {
		Response model.CreatedResponse
		Document model.Model
	}
	type request struct {
		vars  map[string]string
		query url.Values
		body  interface{}
		user  middleware.User
	}
	type response struct {
		StatusCode int
		Body       interface{}
	}
	tests := []struct {
		name     string
		request  request
		response response
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			requestBody, err := json.Marshal(test.request.body)
			if err != nil {
				t.Errorf("handler.go insert() error = %v", err)
				return
			}

			request, err := http.NewRequest(http.MethodPost, "", bytes.NewReader(requestBody))
			if err != nil {
				t.Errorf("handler.go insert() error = %v", err)
				return
			}

			request.URL = &url.URL{RawQuery: test.request.query.Encode()}
			request = mux.SetURLVars(request, test.request.vars)
			context.Set(request, "User", test.request.user)

			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, request)
			response := recorder.Result()

			if !cmp.Equal(test.response.StatusCode, response.StatusCode) {
				t.Errorf("handler.go insert() mismatch (-want +got):\n%s", cmp.Diff(test.response.StatusCode, response.StatusCode))
			}

			var responseBody interface{}
			switch response.StatusCode {
			case http.StatusCreated:
				result := result{}
				body := model.CreatedResponse{}
				if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
					t.Errorf("handler.go insert() error = %v", err)
					return
				}
				result.Response = body

				recorder := httptest.NewRecorder()
				request, err := http.NewRequest(http.MethodGet, "", nil)
				if err != nil {
					t.Errorf("handler.go insert() error = %v", err)
					return
				}

				if body.ID != nil {
					if id, ok := body.ID.(string); ok {
						test.request.vars["id"] = id
					} else if id, ok := body.ID.(primitive.ObjectID); ok {
						test.request.vars["id"] = id.Hex()
					}
				}

				request = mux.SetURLVars(request, test.request.vars)
				context.Set(request, "User", test.request.user)

				get(ctrl).ServeHTTP(recorder, request)
				response := recorder.Result()

				document := model.Model{}
				if err := json.NewDecoder(response.Body).Decode(&document); err != nil {
					t.Errorf("handler.go insert() error = %v", err)
					return
				}
				result.Document = document
				responseBody = result
			case http.StatusForbidden:
				body := model.UnauthorizedResponse{}
				if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
					t.Errorf("handler.go insert() error = %v", err)
					return
				}
				responseBody = body
			default:
				body := model.ErrorResponse{}
				if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
					t.Errorf("handler.go insert() error = %v", err)
					return
				}
				responseBody = body
			}

			if !cmp.Equal(test.response.Body, responseBody, opts) {
				t.Errorf("handler.go insert() mismatch (-want +got):\n%s", cmp.Diff(test.response.Body, responseBody, opts))
			}
		})
	}
}

func Test_find(t *testing.T) {
	ctrl, _ := controller.New(conf)
	if err := ctrl.Ready(); err != nil {
		t.Skipf("handler.go find() error = %v", err)
		return
	}
	defer ctrl.Mongo.Disconnect()

	handler := find(ctrl)

	opts := cmp.Options{
		utils.EquateErrors(),
	}

	type request struct {
		vars  map[string]string
		query url.Values
		body  interface{}
		user  middleware.User
	}
	type response struct {
		StatusCode int
		Body       interface{}
	}
	tests := []struct {
		name     string
		request  request
		response response
	}{
		// TODO: Add test cases.
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			requestBody, err := json.Marshal(test.request.body)
			if err != nil {
				t.Errorf("handler.go find() error = %v", err)
				return
			}

			request, err := http.NewRequest(http.MethodGet, "", bytes.NewReader(requestBody))
			if err != nil {
				t.Errorf("handler.go find() error = %v", err)
				return
			}

			request.URL = &url.URL{RawQuery: test.request.query.Encode()}
			request = mux.SetURLVars(request, test.request.vars)
			context.Set(request, "User", test.request.user)

			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, request)
			response := recorder.Result()

			if !cmp.Equal(test.response.StatusCode, response.StatusCode) {
				t.Errorf("handler.go find() mismatch (-want +got):\n%s", cmp.Diff(test.response.StatusCode, response.StatusCode))
			}

			var responseBody interface{}
			switch response.StatusCode {
			case http.StatusOK:
				body := model.List{}
				if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
					t.Errorf("handler.go find() error = %v", err)
					return
				}
				responseBody = body
			case http.StatusForbidden:
				body := model.UnauthorizedResponse{}
				if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
					t.Errorf("handler.go find() error = %v", err)
					return
				}
				responseBody = body
			default:
				body := model.ErrorResponse{}
				if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
					t.Errorf("handler.go find() error = %v", err)
					return
				}
				responseBody = body
			}

			if !cmp.Equal(test.response.Body, responseBody, opts) {
				t.Errorf("handler.go find() mismatch (-want +got):\n%s", cmp.Diff(test.response.Body, responseBody, opts))
			}
		})
	}
}

func Test_get(t *testing.T) {
	ctrl, _ := controller.New(conf)
	if err := ctrl.Ready(); err != nil {
		t.Skipf("handler.go get() error = %v", err)
		return
	}
	defer ctrl.Mongo.Disconnect()

	handler := get(ctrl)

	opts := cmp.Options{
		utils.EquateErrors(),
	}

	type request struct {
		vars  map[string]string
		query url.Values
		body  interface{}
		user  middleware.User
	}
	type response struct {
		StatusCode int
		Body       interface{}
	}
	tests := []struct {
		name     string
		request  request
		response response
	}{
		// TODO: Add test cases.
		{
			name: "Invalid Hex",
			request: request{
				user: middleware.User{
					GivenName:      "John",
					Sn:             "Doe",
					SAMAccountName: "john.doe",
				},
				vars: map[string]string{"id": "01"},
			},
			response: response{
				StatusCode: http.StatusBadRequest,
				Body: model.ErrorResponse{
					Err: primitive.ErrInvalidHex.Error(),
				},
			},
		},
		{
			name: "Invalid Byte",
			request: request{
				user: middleware.User{
					GivenName:      "John",
					Sn:             "Doe",
					SAMAccountName: "john.doe",
				},
				vars: map[string]string{"id": "ZXCVBNM<ASDFGHJKL:"},
			},
			response: response{
				StatusCode: http.StatusBadRequest,
				Body: model.ErrorResponse{
					Err: primitive.ErrInvalidHex.Error(),
				},
			},
		},
		{
			name: "Invalid Hex Length",
			request: request{
				user: middleware.User{
					GivenName:      "John",
					Sn:             "Doe",
					SAMAccountName: "john.doe",
				},
				vars: map[string]string{"id": "0"},
			},
			response: response{
				StatusCode: http.StatusBadRequest,
				Body: model.ErrorResponse{
					Err: primitive.ErrInvalidHex.Error(),
				},
			},
		},
		{
			name: "No Documents",
			request: request{
				user: middleware.User{
					GivenName:      "John",
					Sn:             "Doe",
					SAMAccountName: "john.doe",
				},
				vars: map[string]string{"id": "000000000000000000000000"},
			},
			response: response{
				StatusCode: http.StatusNotFound,
				Body: model.ErrorResponse{
					Err: mongo.ErrNoDocuments.Error(),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			requestBody, err := json.Marshal(test.request.body)
			if err != nil {
				t.Errorf("handler.go get() error = %v", err)
				return
			}

			request, err := http.NewRequest(http.MethodGet, "", bytes.NewReader(requestBody))
			if err != nil {
				t.Errorf("handler.go get() error = %v", err)
				return
			}

			request.URL = &url.URL{RawQuery: test.request.query.Encode()}
			request = mux.SetURLVars(request, test.request.vars)
			context.Set(request, "User", test.request.user)

			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, request)
			response := recorder.Result()

			if !cmp.Equal(test.response.StatusCode, response.StatusCode) {
				t.Errorf("handler.go get() mismatch (-want +got):\n%s", cmp.Diff(test.response.StatusCode, response.StatusCode))
			}

			var responseBody interface{}
			switch response.StatusCode {
			case http.StatusOK:
				body := model.Model{}
				if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
					t.Errorf("handler.go get() error = %v", err)
					return
				}
				responseBody = body
			case http.StatusForbidden:
				body := model.UnauthorizedResponse{}
				if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
					t.Errorf("handler.go get() error = %v", err)
					return
				}
				responseBody = body
			default:
				body := model.ErrorResponse{}
				if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
					t.Errorf("handler.go get() error = %v", err)
					return
				}
				responseBody = body
			}

			if !cmp.Equal(test.response.Body, responseBody, opts) {
				t.Errorf("handler.go get() mismatch (-want +got):\n%s", cmp.Diff(test.response.Body, responseBody, opts))
			}
		})
	}
}

func Test_update(t *testing.T) {
	ctrl, _ := controller.New(conf)
	if err := ctrl.Ready(); err != nil {
		t.Skipf("handler.go update() error = %v", err)
		return
	}
	defer ctrl.Mongo.Disconnect()
	handler := update(ctrl)

	opts := cmp.Options{
		utils.EquateErrors(),
		cmpopts.IgnoreFields(model.UpdatedResponse{}, "ID"),
		cmpopts.IgnoreFields(model.Audit{}, "UpdatedTs"),
	}

	type result struct {
		Response model.UpdatedResponse
		Document model.Model
	}
	type request struct {
		vars  map[string]string
		query url.Values
		body  interface{}
		user  middleware.User
	}
	type response struct {
		StatusCode int
		Body       interface{}
	}
	tests := []struct {
		name     string
		request  request
		response response
	}{
		// TODO: Add test cases.
		{
			name: "Invalid Hex",
			request: request{
				user: middleware.User{
					GivenName:      "John",
					Sn:             "Doe",
					SAMAccountName: "john.doe",
				},
				vars: map[string]string{"id": "01"},
				body: model.Model{},
			},
			response: response{
				StatusCode: http.StatusBadRequest,
				Body: model.ErrorResponse{
					Err: primitive.ErrInvalidHex.Error(),
				},
			},
		},
		{
			name: "Invalid Byte",
			request: request{
				user: middleware.User{
					GivenName:      "John",
					Sn:             "Doe",
					SAMAccountName: "john.doe",
				},
				vars: map[string]string{"id": "ZXCVBNM<ASDFGHJKL:"},
				body: model.Model{},
			},
			response: response{
				StatusCode: http.StatusBadRequest,
				Body: model.ErrorResponse{
					Err: primitive.ErrInvalidHex.Error(),
				},
			},
		},
		{
			name: "Invalid Hex Length",
			request: request{
				user: middleware.User{
					GivenName:      "John",
					Sn:             "Doe",
					SAMAccountName: "john.doe",
				},
				vars: map[string]string{"id": "0"},
				body: model.Model{},
			},
			response: response{
				StatusCode: http.StatusBadRequest,
				Body: model.ErrorResponse{
					Err: primitive.ErrInvalidHex.Error(),
				},
			},
		},
		{
			name: "No Documents",
			request: request{
				user: middleware.User{
					GivenName:      "John",
					Sn:             "Doe",
					SAMAccountName: "john.doe",
				},
				vars: map[string]string{"id": "000000000000000000000000"},
				body: model.Model{},
			},
			response: response{
				StatusCode: http.StatusNotFound,
				Body: model.ErrorResponse{
					Err: mongo.ErrNoDocuments.Error(),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			requestBody, err := json.Marshal(test.request.body)
			if err != nil {
				t.Errorf("handler.go update() error = %v", err)
				return
			}

			request, err := http.NewRequest(http.MethodPut, "", bytes.NewReader(requestBody))
			if err != nil {
				t.Errorf("handler.go update() error = %v", err)
				return
			}

			request.URL = &url.URL{RawQuery: test.request.query.Encode()}
			request = mux.SetURLVars(request, test.request.vars)
			context.Set(request, "User", test.request.user)

			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, request)
			response := recorder.Result()

			if !cmp.Equal(test.response.StatusCode, response.StatusCode) {
				t.Errorf("handler.go update() mismatch (-want +got):\n%s", cmp.Diff(test.response.StatusCode, response.StatusCode))
			}

			var responseBody interface{}
			switch response.StatusCode {
			case http.StatusOK:
				result := result{}
				body := model.UpdatedResponse{}
				if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
					t.Errorf("handler.go update() error = %v", err)
					return
				}
				result.Response = body

				recorder := httptest.NewRecorder()
				request, err := http.NewRequest(http.MethodGet, "", nil)
				if err != nil {
					t.Errorf("handler.go update() error = %v", err)
					return
				}

				if body.ID != nil {
					if id, ok := body.ID.(string); ok {
						test.request.vars["id"] = id
					} else if id, ok := body.ID.(primitive.ObjectID); ok {
						test.request.vars["id"] = id.Hex()
					}
				}

				request = mux.SetURLVars(request, test.request.vars)
				context.Set(request, "User", test.request.user)

				get(ctrl).ServeHTTP(recorder, request)
				response := recorder.Result()

				document := model.Model{}
				if err := json.NewDecoder(response.Body).Decode(&document); err != nil {
					t.Errorf("handler.go update() error = %v", err)
					return
				}
				result.Document = document
				responseBody = result
			case http.StatusConflict:
				body := model.ConflictResponse{}
				if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
					t.Errorf("handler.go update() error = %v", err)
					return
				}
				responseBody = body
			case http.StatusForbidden:
				body := model.UnauthorizedResponse{}
				if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
					t.Errorf("handler.go update() error = %v", err)
					return
				}
				responseBody = body
			default:
				body := model.ErrorResponse{}
				if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
					t.Errorf("handler.go update() error = %v", err)
					return
				}
				responseBody = body
			}

			if !cmp.Equal(test.response.Body, responseBody, opts) {
				t.Errorf("handler.go update() mismatch (-want +got):\n%s", cmp.Diff(test.response.Body, responseBody, opts))
			}
		})
	}
}

func Test_delete(t *testing.T) {
	ctrl, _ := controller.New(conf)
	if err := ctrl.Ready(); err != nil {
		t.Skipf("handler.go delete() error = %v", err)
		return
	}
	defer ctrl.Mongo.Disconnect()
	handler := delete(ctrl)

	opts := cmp.Options{
		utils.EquateErrors(),
	}

	type request struct {
		vars  map[string]string
		query url.Values
		body  interface{}
		user  middleware.User
	}
	type response struct {
		StatusCode int
		Body       interface{}
	}
	tests := []struct {
		name     string
		request  request
		response response
	}{
		{
			name: "Invalid Hex",
			request: request{
				user: middleware.User{
					GivenName:      "John",
					Sn:             "Doe",
					SAMAccountName: "john.doe",
				},
				vars: map[string]string{"id": "01"},
			},
			response: response{
				StatusCode: http.StatusBadRequest,
				Body: model.ErrorResponse{
					Err: primitive.ErrInvalidHex.Error(),
				},
			},
		},
		{
			name: "Invalid Byte",
			request: request{
				user: middleware.User{
					GivenName:      "John",
					Sn:             "Doe",
					SAMAccountName: "john.doe",
				},
				vars: map[string]string{"id": "ZXCVBNM<ASDFGHJKL:"},
			},
			response: response{
				StatusCode: http.StatusBadRequest,
				Body: model.ErrorResponse{
					Err: primitive.ErrInvalidHex.Error(),
				},
			},
		},
		{
			name: "Invalid Hex Length",
			request: request{
				user: middleware.User{
					GivenName:      "John",
					Sn:             "Doe",
					SAMAccountName: "john.doe",
				},
				vars: map[string]string{"id": "0"},
			},
			response: response{
				StatusCode: http.StatusBadRequest,
				Body: model.ErrorResponse{
					Err: primitive.ErrInvalidHex.Error(),
				},
			},
		},
		{
			name: "No Documents",
			request: request{
				user: middleware.User{
					GivenName:      "John",
					Sn:             "Doe",
					SAMAccountName: "john.doe",
				},
				vars: map[string]string{"id": "000000000000000000000000"},
			},
			response: response{
				StatusCode: http.StatusNotFound,
				Body: model.ErrorResponse{
					Err: mongo.ErrNoDocuments.Error(),
				},
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			requestBody, err := json.Marshal(test.request.body)
			if err != nil {
				t.Errorf("handler.go delete() error = %v", err)
				return
			}

			request, err := http.NewRequest(http.MethodPut, "", bytes.NewReader(requestBody))
			if err != nil {
				t.Errorf("handler.go delete() error = %v", err)
				return
			}

			request.URL = &url.URL{RawQuery: test.request.query.Encode()}
			request = mux.SetURLVars(request, test.request.vars)
			context.Set(request, "User", test.request.user)

			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, request)
			response := recorder.Result()

			if !cmp.Equal(test.response.StatusCode, response.StatusCode) {
				t.Errorf("handler.go delete() mismatch (-want +got):\n%s", cmp.Diff(test.response.StatusCode, response.StatusCode))
			}

			var responseBody interface{}
			switch response.StatusCode {
			case http.StatusNoContent:
				recorder := httptest.NewRecorder()
				request, err := http.NewRequest(http.MethodGet, "", nil)
				if err != nil {
					t.Errorf("handler.go delete() error = %v", err)
					return
				}

				request = mux.SetURLVars(request, test.request.vars)
				context.Set(request, "User", test.request.user)

				get(ctrl).ServeHTTP(recorder, request)
				response := recorder.Result()

				switch response.StatusCode {
				case http.StatusNotFound:
					responseBody = nil
				default:
					document := model.Model{}
					if err := json.NewDecoder(response.Body).Decode(&document); err != nil {
						t.Errorf("handler.go delete() error = %v", err)
						return
					}
					responseBody = document
				}
			case http.StatusForbidden:
				body := model.UnauthorizedResponse{}
				if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
					t.Errorf("handler.go delete() error = %v", err)
					return
				}
				responseBody = body
			default:
				body := model.ErrorResponse{}
				if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
					t.Errorf("handler.go delete() error = %v", err)
					return
				}
				responseBody = body
			}

			if !cmp.Equal(test.response.Body, responseBody, opts) {
				t.Errorf("handler.go delete() mismatch (-want +got):\n%s", cmp.Diff(test.response.Body, responseBody, opts))
			}
		})
	}
}
