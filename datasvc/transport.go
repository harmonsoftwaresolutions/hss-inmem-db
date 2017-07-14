// Package datasvc provides ...
package datasvc

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"

	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
)

var (
	// ErrBadRouting bad path
	ErrBadRouting = errors.New("inconsistent mapping between route and handler (programming error)")
)

// MakeHTTPHandler ...
func MakeHTTPHandler(s Service, logger log.Logger) http.Handler {
	r := mux.NewRouter()
	e := MakeServerEndpoints(s)
	options := []httptransport.ServerOption{
		httptransport.ServerErrorLogger(logger),
		httptransport.ServerErrorEncoder(encodeError),
	}

	// POST		/notes/		adds another note
	// GET   	/notes/:id	retrieves given note by id
	// PUT		/notes/:id      post updated note information about the note
	// PATCH	/notes/:id      partial updated note information
	// DELETE	/notes/:id      remove the given note

	r.Methods("POST").Path("/notes").Handler(httptransport.NewServer(
		e.PostNoteEndpoint,
		decodePostNoteRequest,
		encodeResponse,
		options...,
	))
	r.Methods("GET").Path("/notes/{id}").Handler(httptransport.NewServer(
		e.GetNoteEndpoint,
		decodeGetNoteRequest,
		encodeResponse,
		options...,
	))
	r.Methods("PUT").Path("/notes/{id}").Handler(httptransport.NewServer(
		e.PutNoteEndpoint,
		decodePutNoteRequest,
		encodeResponse,
		options...,
	))
	r.Methods("PATCH").Path("/notes/{id}").Handler(httptransport.NewServer(
		e.PatchNoteEndpoint,
		decodePatchNoteRequest,
		encodeResponse,
		options...,
	))
	r.Methods("DELETE").Path("/notes/{id}").Handler(httptransport.NewServer(
		e.DeleteNoteEndpoint,
		decodeDeleteNoteRequest,
		encodeResponse,
		options...,
	))
	return r
}

func decodePostNoteRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	var req postNoteRequest
	if e := json.NewDecoder(r.Body).Decode(&req.Note); e != nil {
		return nil, e
	}
	return req, nil
}

func decodeGetNoteRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return getNoteRequest{ID: id}, nil
}

func decodePutNoteRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	var note Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		return nil, err
	}
	return putNoteRequest{
		ID:   id,
		Note: note,
	}, nil
}

func decodePatchNoteRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	var note Note
	if err := json.NewDecoder(r.Body).Decode(&note); err != nil {
		return nil, err
	}
	return patchNoteRequest{
		ID:   id,
		Note: note,
	}, nil
}

func decodeDeleteNoteRequest(_ context.Context, r *http.Request) (request interface{}, err error) {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return nil, ErrBadRouting
	}
	return deleteNoteRequest{ID: id}, nil
}

func encodePostNoteRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("POST").Path("/notes/")
	req.Method, req.URL.Path = "POST", "/notes/"
	return encodeRequest(ctx, req, request)
}

func encodeGetNoteRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("GET").Path("/notes/{id}")
	r := request.(getNoteRequest)
	noteID := url.QueryEscape(r.ID)
	req.Method, req.URL.Path = "GET", "/notes/"+noteID
	return encodeRequest(ctx, req, request)
}

func encodePutNoteRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("PUT").Path("/notes/{id}")
	r := request.(putNoteRequest)
	noteID := url.QueryEscape(r.ID)
	req.Method, req.URL.Path = "PUT", "/notes/"+noteID
	return encodeRequest(ctx, req, request)
}

func encodePatchNoteRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("PATCH").Path("/notes/{id}")
	r := request.(patchNoteRequest)
	noteID := url.QueryEscape(r.ID)
	req.Method, req.URL.Path = "PATCH", "/notes/"+noteID
	return encodeRequest(ctx, req, request)
}

func encodeDeleteNoteRequest(ctx context.Context, req *http.Request, request interface{}) error {
	// r.Methods("DELETE").Path("/notes/{id}")
	r := request.(deleteNoteRequest)
	noteID := url.QueryEscape(r.ID)
	req.Method, req.URL.Path = "DELETE", "/notes/"+noteID
	return encodeRequest(ctx, req, request)
}

func decodePostNoteResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response postNoteResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeGetNoteResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response getNoteResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodePutNoteResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response putNoteResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodePatchNoteResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response patchNoteResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

func decodeDeleteNoteResponse(_ context.Context, resp *http.Response) (interface{}, error) {
	var response deleteNoteResponse
	err := json.NewDecoder(resp.Body).Decode(&response)
	return response, err
}

type errorer interface {
	error() error
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	if e, ok := response.(errorer); ok && e.error() != nil {
		encodeError(ctx, e.error(), w)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func encodeRequest(_ context.Context, req *http.Request, request interface{}) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(request)
	if err != nil {
		return err
	}
	req.Body = ioutil.NopCloser(&buf)
	return nil
}

func encodeError(_ context.Context, err error, w http.ResponseWriter) {
	if err != nil {
		panic("encodeError with nil error")
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(codeFrom(err))
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": err.Error(),
	})
}

func codeFrom(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrAlreadyExists, ErrInconsistentIDs:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
