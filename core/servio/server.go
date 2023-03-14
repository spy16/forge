package servio

import (
	"context"
	"net/http"
	"time"

	"google.golang.org/appengine/log"

	"github.com/spy16/forge/core/errors"
)

const gracePeriod = 5 * time.Second

type ResponseWriterCapture struct {
	http.ResponseWriter

	Status int
}

func (rwc *ResponseWriterCapture) WriteHeader(status int) {
	rwc.Status = status
	rwc.ResponseWriter.WriteHeader(status)
}

// Serve serves the given handler at addr.
func Serve(ctx context.Context, addr string, handler http.Handler) error {
	errCh := make(chan error)

	httpServer := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	go func() {
		log.Infof(ctx, "listening on %s...", httpServer.Addr)
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err

	case <-ctx.Done():
		log.Infof(ctx, "context cancelled, shutting down")

		graceCtx, cancel := context.WithTimeout(context.Background(), gracePeriod)
		defer cancel()
		return httpServer.Shutdown(graceCtx)
	}
}
