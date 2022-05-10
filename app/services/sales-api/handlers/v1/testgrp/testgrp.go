// testgroup package v1
package testgrp

import (
	"encoding/json"
	"net/http"
	"context"
	"go.uber.org/zap"
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
	statusCode := http.StatusOK

	
	h.Log.Infow(" Test  readiness", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)
	return 	json.NewEncoder(w).Encode(data)
}