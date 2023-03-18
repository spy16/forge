package servio

import (
	"context"
	"net/http"
	"time"

	"github.com/spy16/forge/core/errors"
	"github.com/spy16/forge/core/log"
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
		log.Info(ctx, "starting server", map[string]any{
			"addr": httpServer.Addr,
		})
		err := httpServer.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err

	case <-ctx.Done():
		log.Info(ctx, "context cancelled, shutting down", map[string]any{
			"reason": ctx.Err(),
		})

		graceCtx, cancel := context.WithTimeout(context.Background(), gracePeriod)
		defer cancel()
		return httpServer.Shutdown(graceCtx)
	}
}
