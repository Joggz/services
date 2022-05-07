// testgroup package v1
package testgrp

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

type Handlers  struct{
	Log *zap.SugaredLogger
}

// Test Handler for development
func (h Handlers) Test(w http.ResponseWriter, r *http.Request) {
	data := struct{
		Status string 
	}{
		Status: "OK",
	}
	statusCode := http.StatusOK
	json.NewEncoder(w).Encode(data)
	
	h.Log.Infow(" Test  readiness", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)

}