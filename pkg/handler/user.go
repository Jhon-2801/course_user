package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Jhon-2801/course_user/internal/user"
	"github.com/Jhon-2801/lib-response/response"
	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func NewUserHTTPServer(ctx context.Context, endpoints user.Endpoints) http.Handler {
	r := mux.NewRouter()

	opts := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(encodeError),
	}
	// Handler Create
	r.Handle("/users", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Create),
		decodeCreateUser, encodeResponse,
		opts...,
	)).Methods("POST")

	// Handler Get
	r.Handle("/users/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Get),
		decodeGetUser, encodeResponse,
		opts...,
	)).Methods("GET")

	// Handler GetAll
	r.Handle("/users", httptransport.NewServer(
		endpoint.Endpoint(endpoints.GetAll),
		decodeGetAllUser, encodeResponse,
		opts...,
	)).Methods("GET")

	// Handler Update
	r.Handle("/users/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Update),
		decodeUpdateUser, encodeResponse,
		opts...,
	)).Methods("PATCH")

	// Handler Delete
	r.Handle("/users/{id}", httptransport.NewServer(
		endpoint.Endpoint(endpoints.Delete),
		decodeDeleteUser, encodeResponse,
		opts...,
	)).Methods("DELETE")

	return r
}

// Create
func decodeCreateUser(_ context.Context, r *http.Request) (interface{}, error) {
	var req user.CreateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid a request format: %v", err.Error()))
	}

	return req, nil
}

// Get
func decodeGetUser(_ context.Context, r *http.Request) (interface{}, error) {
	p := mux.Vars(r)
	req := user.GetReq{
		ID: p["id"],
	}
	return req, nil
}

// GetAll
func decodeGetAllUser(_ context.Context, r *http.Request) (interface{}, error) {

	v := r.URL.Query()

	limit, _ := strconv.Atoi(v.Get("limit"))
	page, _ := strconv.Atoi(v.Get("page"))

	req := user.GetAllReq{
		FirstName: v.Get("first_name"),
		LastName:  v.Get("last_name"),
		Limit:     limit,
		Page:      page,
	}
	return req, nil
}

// Update
func decodeUpdateUser(_ context.Context, r *http.Request) (interface{}, error) {
	var req user.UpdateReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, response.BadRequest(fmt.Sprintf("invalid a request format: %v", err.Error()))
	}
	path := mux.Vars(r)
	req.ID = path["id"]

	return req, nil
}

// Delete

func decodeDeleteUser(_ context.Context, r *http.Request) (interface{}, error) {
	p := mux.Vars(r)
	req := user.DeleteReq{
		ID: p["id"],
	}
	return req, nil
}

// Response
func encodeResponse(ctx context.Context, w http.ResponseWriter, resp interface{}) error {
	r := resp.(response.Response)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(r.StatusCode())
	return json.NewEncoder(w).Encode(r)
}

// Error
func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; chatset=utf-8")
	resp := err.(response.Response)
	w.WriteHeader(resp.StatusCode())
	_ = json.NewEncoder(w).Encode(resp)
}
