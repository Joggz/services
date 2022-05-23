package mid

import (
	"context"
	"net/http"

	"github.com/Joggz/services/business/sys/metrics"
	"github.com/Joggz/services/foundation/web"
)


func Metrics() web.Middleware{
	m := func (handler web.Handler) web.Handler  {
		
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			// Add the metrics into the context for metric gathering.
			ctx = metrics.Set(ctx)

			err := handler(ctx, w, r)
				// Handle updating the metrics that can be handled here.
   
			// Increment the request and goroutines counter.
			metrics.AddRequest(ctx)
			metrics.AdddGoroutine(ctx)

			if err != nil {
				metrics.AddErrors(ctx)
			}

			return err
		}
		return h
	}

	return m
}