// testgroup package v1
package testgrp

import (
	// "encoding/json"
	"net/http"
	"context"
	"go.uber.org/zap"

	"github.com/Joggz/services/foundation/web"
)

type Handlers  struct{
	Log *zap.SugaredLogger
}

// Test Handler for development
func (h Handlers) Test(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
	data := struct{
		Status string 
	}{
		Status: "OK",
	}
	
	return  web.Respond(ctx, w, data, http.StatusOK)	
	// json.NewEncoder(w).Encode(data)
}