// Package datasvc provides ...
package datasvc

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

// LoggingMiddleware ...
func LoggingMiddleware(logger log.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

type loggingMiddleware struct {
	next   Service
	logger log.Logger
}

func (mw loggingMiddleware) PostNote(ctx context.Context, n Note) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PostNote", "id", n.ID, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PostNote(ctx, n)
}

func (mw loggingMiddleware) GetNote(ctx context.Context, id string) (n Note, err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "GetNote", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.GetNote(ctx, id)
}

func (mw loggingMiddleware) PutNote(ctx context.Context, id string, n Note) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PutNote", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PutNote(ctx, id, n)
}

func (mw loggingMiddleware) PatchNote(ctx context.Context, id string, n Note) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "PatchNote", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.PatchNote(ctx, id, n)
}

func (mw loggingMiddleware) DeleteNote(ctx context.Context, id string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "DeleteNote", "id", id, "took", time.Since(begin), "err", err)
	}(time.Now())
	return mw.next.DeleteNote(ctx, id)
}
