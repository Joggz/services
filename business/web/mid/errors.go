package mid

import (
	"context"
	"net/http"

	"github.com/Joggz/services/business/sys/validate"
	"github.com/Joggz/services/foundation/web"

	"go.uber.org/zap"
)

// Errors handles errors coming out of the call chain. It detects normal
// application errors which are used to respond to the client in a uniform way.
// Unexpected errors (status >= 500) are logged.

func Error (log *zap.SugaredLogger) web.Middleware {

	// Create the handler that will be attached in the middleware chain.
	m := func(handler web.Handler) web.Handler{

	
	h := func (ctx context.Context, w http.ResponseWriter, r *http.Request) error {

		// If the context is missing this value, request the service
			// to be shutdown gracefully.
			v, err := web.GetValues(ctx)
			if err != nil {
				return web.NewShutdownError("web value missing")
			}

			if err := handler(ctx, w, r); err != nil {
				// Log the error.
				log.Errorw("ERROR", "traceid", v.TraceID, "message", err)

				// Build out the error response.
				// Build out the error response.
				var er validate.ErrorResponse
				var status int
			

				switch { 
				case validate.IsFieldErrors(err):
					fieldErrors := validate.GetFieldErrors(err)
					er = validate.ErrorResponse{
						Error:  "data validation error",
						Fields: fieldErrors.Fields(),
					}
					status = http.StatusBadRequest

				case validate.IsRequestError(err):
					reqErr := validate.GetRequestError(err)
					er = validate.ErrorResponse{
						Error: reqErr.Error(),
					}
					status = reqErr.Status

				default:
					er = validate.ErrorResponse{
						Error: http.StatusText(http.StatusInternalServerError),
					}
					status = http.StatusInternalServerError
				}

				// Respond with the error back to the client.
				if err := web.Respond(ctx, w, er, status); err != nil {
					return err
				}

				// If we receive the shutdown err we need to return it
				// back to the base handler to shut down the service.
				if web.IsShutdown(err) {
					return err
				}
			}

		return nil
	}
	return h
}

	return m
	
}