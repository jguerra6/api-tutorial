package middleware

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/jguerra6/api-tutorial/internal/transport/http/writer"
)

func Timeout(d time.Duration) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), d)
			defer cancel()

			tw := &timeoutWriter{ResponseWriter: w}
			done := make(chan struct{})

			go func() {
				next.ServeHTTP(tw, r.WithContext(ctx))
				close(done)
			}()

			select {
			case <-done:
				return
			case <-ctx.Done():
				tw.markTimedOut()
				writer.JSON(w, http.StatusGatewayTimeout, map[string]any{
					"error":   "timeout",
					"message": "request exceeded time limit",
				})
				return
			}
		})
	}
}

type timeoutWriter struct {
	http.ResponseWriter
	mu        sync.Mutex
	timedOut  bool
	wroteHead bool
}

func (tw *timeoutWriter) markTimedOut() {
	tw.mu.Lock()
	tw.timedOut = true
	tw.mu.Unlock()
}

func (tw *timeoutWriter) WriteHeader(code int) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut || tw.wroteHead {
		return
	}
	tw.wroteHead = true
	tw.ResponseWriter.WriteHeader(code)
}

func (tw *timeoutWriter) Write(b []byte) (int, error) {
	tw.mu.Lock()
	defer tw.mu.Unlock()
	if tw.timedOut {
		return 0, http.ErrHandlerTimeout
	}
	return tw.ResponseWriter.Write(b)
}
