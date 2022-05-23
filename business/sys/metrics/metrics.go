// Package metrics constructs the metrics the application will track.
package metrics

import (
	"expvar"
)

// metrics represents the set of metrics we gather. These fields are
// safe to be accessed concurrently thanks to expvar. No extra abstraction is required.

type metrics struct {
	goroutines *expvar.Int
	request *expvar.Int
	errors *expvar.Int
	panicks *expvar.Int
}

var m  *metrics


// init constructs the metrics value that will be used to capture metrics.
// The metrics value is stored in a package level variable since everything
// inside of expvar is registered as a singleton. The use of once will make
// sure this initialization only happens once.
func init() {
	m = &metrics{
		goroutines: expvar.NewInt("goroutines"),
		request: expvar.NewInt("request"),
		errors: expvar.NewInt("errors"),
		panicks: expvar.NewInt("panicks"),
	}
}


