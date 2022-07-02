// Package metrics constructs the metrics the application will track.
package metrics

import (
	"context"
	"expvar"
	"runtime"
)

// metrics represents the set of metrics we gather. These fields are
// safe to be accessed concurrently thanks to expvar. No extra abstraction is required.

type metrics struct {
	goroutines *expvar.Int
	requests *expvar.Int
	errors *expvar.Int
	panics *expvar.Int
}

var m  *metrics


// init constructs the metrics value that will be used to capture metrics.
// The metrics value is stored in a package level variable since everything
// inside of expvar is registered as a singleton. The use of once will make
// sure this initialization only happens once.
func init() {
	m = &metrics{
		goroutines: expvar.NewInt("goroutines"),
		requests: expvar.NewInt("requests"),
		errors: expvar.NewInt("errors"),
		panics: expvar.NewInt("panics"),
	}
}


// =============================================================================

// Metrics will be supported through the context.

// ctxKeyMetric represents the type of value for the context key.
type ctxKey int

// key is how metric values are stored/retrieved.
const key ctxKey = 1

// =============================================================================


// Add more of these functions when a metric needs to be collected in
// different parts of the codebase. This will keep this package the
// central authority for metrics and metrics won't get lost.

func Set(ctx context.Context)  context.Context{
	return context.WithValue(ctx, key, m)
}

func AdddGoroutine(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		if v.requests.Value() % 100 == 0{
			v.goroutines.Set(int64(runtime.NumGoroutine()))
		}
	}
}


func AddRequest(ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.requests.Add(1)
	}
}

func AddErrors (ctx context.Context) {
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.errors.Add(1)
	}
}

func AddPanick (ctx context.Context) { 
	if v, ok := ctx.Value(key).(*metrics); ok {
		v.panics.Add(1)
	}
}