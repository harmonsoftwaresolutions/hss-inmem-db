// Package datasvc provides ...
package datasvc

import (
	"context"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
)

// Endpoints ...
type Endpoints struct {
	PostNoteEndpoint   endpoint.Endpoint
	GetNoteEndpoint    endpoint.Endpoint
	PutNoteEndpoint    endpoint.Endpoint
	PatchNoteEndpoint  endpoint.Endpoint
	DeleteNoteEndpoint endpoint.Endpoint
}

// MakeServerEndpoints ...
func MakeServerEndpoints(s Service) Endpoints {
	return Endpoints{
		PostNoteEndpoint:   MakePostNoteEndpoint(s),
		GetNoteEndpoint:    MakeGetNoteEndpoint(s),
		PutNoteEndpoint:    MakePutNoteEndpoint(s),
		PatchNoteEndpoint:  MakePatchNoteEndpoint(s),
		DeleteNoteEndpoint: MakeDeleteNoteEndpoint(s),
	}
}

// MakeClientEndpoints ...
func MakeClientEndpoints(instance string) (Endpoints, error) {
	if !strings.HasPrefix(instance, "http") {
		instance = "http://" + instance
	}
	tgt, err := url.Parse(instance)
	if err != nil {
		return Endpoints{}, err
	}
	tgt.Path = ""

	options := []httptransport.ClientOption{}

	return Endpoints{
		PostNoteEndpoint:   httptransport.NewClient("POST", tgt, encodePostNoteRequest, decodePostNoteResponse, options...).Endpoint(),
		GetNoteEndpoint:    httptransport.NewClient("GET", tgt, encodeGetNoteRequest, decodeGetNoteResponse, options...).Endpoint(),
		PutNoteEndpoint:    httptransport.NewClient("PUT", tgt, encodePutNoteRequest, decodePutNoteResponse, options...).Endpoint(),
		PatchNoteEndpoint:  httptransport.NewClient("PATCH", tgt, encodePatchNoteRequest, decodePatchNoteResponse, options...).Endpoint(),
		DeleteNoteEndpoint: httptransport.NewClient("DELETE", tgt, encodeDeleteNoteRequest, decodeDeleteNoteResponse, options...).Endpoint(),
	}, nil
}

// PostNote implements Service. Primarily useful in a client.
func (e Endpoints) PostNote(ctx context.Context, n Note) error {
	request := postNoteRequest{Note: n}
	response, err := e.PostNoteEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(postNoteResponse)
	return resp.Err
}

// GetNote ...
func (e Endpoints) GetNote(ctx context.Context, id string) (Note, error) {
	request := getNoteRequest{ID: id}
	response, err := e.GetNoteEndpoint(ctx, request)
	if err != nil {
		return Note{}, err
	}
	resp := response.(getNoteResponse)
	return resp.Note, resp.Err
}

// PutNote implements Service. Primarily useful in a client.
func (e Endpoints) PutNote(ctx context.Context, id string, n Note) error {
	request := putNoteRequest{ID: id, Note: n}
	response, err := e.PutNoteEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(putNoteResponse)
	return resp.Err
}

// PatchNote implements Service. Primarily useful in a client.
func (e Endpoints) PatchNote(ctx context.Context, id string, n Note) error {
	request := patchNoteRequest{ID: id, Note: n}
	response, err := e.PatchNoteEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(patchNoteResponse)
	return resp.Err
}

// DeleteNote implements Service. Primarily useful in a client.
func (e Endpoints) DeleteNote(ctx context.Context, id string) error {
	request := deleteNoteRequest{ID: id}
	response, err := e.DeleteNoteEndpoint(ctx, request)
	if err != nil {
		return err
	}
	resp := response.(deleteNoteResponse)
	return resp.Err
}

// MakePostNoteEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePostNoteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(postNoteRequest)
		e := s.PostNote(ctx, req.Note)
		return postNoteResponse{Err: e}, nil
	}
}

// MakeGetNoteEndpoint ...
func MakeGetNoteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(getNoteRequest)
		n, e := s.GetNote(ctx, req.ID)
		return getNoteResponse{Note: n, Err: e}, nil
	}
}

// MakePutNoteEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePutNoteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(putNoteRequest)
		e := s.PutNote(ctx, req.ID, req.Note)
		return putNoteResponse{Err: e}, nil
	}
}

// MakePatchNoteEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakePatchNoteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(patchNoteRequest)
		e := s.PatchNote(ctx, req.ID, req.Note)
		return patchNoteResponse{Err: e}, nil
	}
}

// MakeDeleteNoteEndpoint returns an endpoint via the passed service.
// Primarily useful in a server.
func MakeDeleteNoteEndpoint(s Service) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req := request.(deleteNoteRequest)
		e := s.DeleteNote(ctx, req.ID)
		return deleteNoteResponse{Err: e}, nil
	}
}

type postNoteRequest struct {
	Note Note
}

type postNoteResponse struct {
	Err error `json:"err,omitempty"`
}

func (r postNoteResponse) error() error { return r.Err }

type getNoteRequest struct {
	ID string
}

type getNoteResponse struct {
	Note Note  `json:"note,omitempty"`
	Err  error `json:"err,omitempty"`
}

func (r getNoteResponse) error() error {
	return r.Err
}

type putNoteRequest struct {
	ID   string
	Note Note
}

type putNoteResponse struct {
	Err error `json:"err,omitempty"`
}

func (r putNoteResponse) error() error { return nil }

type patchNoteRequest struct {
	ID   string
	Note Note
}

type patchNoteResponse struct {
	Err error `json:"err,omitempty"`
}

func (r patchNoteResponse) error() error { return r.Err }

type deleteNoteRequest struct {
	ID string
}

type deleteNoteResponse struct {
	Err error `json:"err,omitempty"`
}

func (r deleteNoteResponse) error() error { return r.Err }
