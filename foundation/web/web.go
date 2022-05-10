// Package web contains a small web framework extension.
package web

import (
	// "context"
	// "net/http"
	"os"
	// "syscall"
	// "time"

	"github.com/dimfeld/httptreemux/v5"
	// "go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	// "go.opentelemetry.io/otel/trace"
)




// App is the entrypoint into our application and what configures our context
// object for each of our http handlers. Feel free to add any configuration
// data/logic on this App struct.
type App struct {
    *httptreemux.ContextMux
	shutdown chan os.Signal
}