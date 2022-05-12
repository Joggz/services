package mid

import (
	"context"
	"net/http"
	"time"

	"github.com/Joggz/services/foundation/web"
	"go.uber.org/zap"
)

// Logger writes some information about the request to the logs in the
// format: TraceID : (200) GET /foo -> IP ADDR (latency)
func Logger(log *zap.SugaredLogger) web.Middleware {

	// This is the actual middleware function to be executed.
	m := func(handler web.Handler) web.Handler {

		// Create the handler that will be attached in the middleware chain.
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {

			traceID := "0000000000"
			StatusCode := http.StatusOK
			now := time.Now()
			
			log.Infow("request started", "traceid", traceID, "method", r.Method, "path", r.URL.Path,
				"remoteaddr", r.RemoteAddr)

			// Call the next handler.
			err := handler(ctx, w, r)

			log.Infow("request completed", "traceid", traceID, "method", r.Method, "path", r.URL.Path,
				"remoteaddr", r.RemoteAddr, "statuscode", StatusCode, "since", time.Since(now))

			// Return the error so it can be handled further up the chain.
			return err
		}

		return h
	}

	return m
}
