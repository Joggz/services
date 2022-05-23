// testgroup package v1
package testgrp

import (
	// "encoding/json"
	"context"
	"math/rand"
	"net/http"

	"go.uber.org/zap"

	"github.com/Joggz/services/foundation/web"
)

type Handlers  struct{
	Log *zap.SugaredLogger
}

// Test Handler for development
func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	if n := rand.Intn(100); n%2 == 0 {
		// return errors.New("untrusted error")
		return web.NewShutdownError("shutdown service")
		// return validate.NewRequestError(errors.New("testing error"), http.StatusBadRequest)
	}
	data := struct{
		Status string 
	}{
		Status: "OK",
	}
	
	return  web.Respond(ctx, w, data, http.StatusOK)	
	// json.NewEncoder(w).Encode(data)
}