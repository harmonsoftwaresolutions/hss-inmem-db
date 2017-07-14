// Package datasvc provides ...
package datasvc

import (
	"context"
	"errors"
	"sync"
)

// Service is a simple CRUD interface for notes
type Service interface {
	PostNote(ctx context.Context, n Note) error
	GetNote(ctx context.Context, id string) (Note, error)
	PutNote(ctx context.Context, id string, n Note) error
	PatchNote(ctx context.Context, id string, n Note) error
	DeleteNote(ctx context.Context, id string) error
}

// Note represents a note
type Note struct {
	ID     string `json:"id"`
	Author string `json:"author"`
	Title  string `json:"title"`
	Text   string `json:"text"`
}

// Errors
var (
	ErrInconsistentIDs = errors.New("inconsisent IDs")
	ErrAlreadyExists   = errors.New("already exists")
	ErrNotFound        = errors.New("not found")
)

type inmemService struct {
	mtx sync.RWMutex
	m   map[string]Note
}

// NewInmemService ...
func NewInmemService() Service {
	return &inmemService{
		m: map[string]Note{},
	}
}

func (s *inmemService) PostNote(ctx context.Context, n Note) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.m[n.ID]; ok {
		return ErrAlreadyExists // POST = create, don't overwrite
	}
	s.m[n.ID] = n
	return nil
}

func (s *inmemService) GetNote(ctx context.Context, id string) (Note, error) {
	s.mtx.RLock()
	// defer s.mtx.Unlock()
	n, ok := s.m[id]
	if !ok {
		return Note{}, ErrNotFound
	}
	return n, nil
}

func (s *inmemService) PutNote(ctx context.Context, id string, n Note) error {
	if id != n.ID {
		return ErrInconsistentIDs
	}
	s.mtx.Lock()
	defer s.mtx.Unlock()
	s.m[id] = n // PUT = create or update
	return nil
}

func (s *inmemService) PatchNote(ctx context.Context, id string, n Note) error {
	if n.ID != "" && id != n.ID {
		return ErrInconsistentIDs
	}

	s.mtx.Lock()
	defer s.mtx.Unlock()

	existing, ok := s.m[id]
	if !ok {
		return ErrNotFound // PATCH = update existing, don't create
	}

	// We assume that it's not possible to PATCH the ID, and that it's not
	// possible to PATCH any field to its zero value. That is, the zero value
	// means not specified. The way around this is to use e.g. Name *string in
	// the Note definition. But since this is just a demonstrative example,
	// I'm leaving that out.

	if n.Author != "" {
		existing.Author = n.Author
	}
	s.m[id] = existing
	return nil
}

func (s *inmemService) DeleteNote(ctx context.Context, id string) error {
	s.mtx.Lock()
	defer s.mtx.Unlock()
	if _, ok := s.m[id]; !ok {
		return ErrNotFound
	}
	delete(s.m, id)
	return nil
}
