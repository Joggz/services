// Package checkgrp maintains the group of handlers for health checking.
package checkgrp

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/Joggz/services/business/sys/database"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

type Handlers  struct{
	Build string
	Log *zap.SugaredLogger
	DB *sqlx.DB
}


// Readiness checks if the database is ready and if not will return a 500 status.
// Do not respond by just returning an error because further up in the call
// stack it will interpret that as a non-trusted error.
func (h Handlers) Readiness(w http.ResponseWriter, r *http.Request) {
	// data := struct{
	// 	Status string `json:"status"`
	// }{
	// 	Status: "OK",
	// }
	// statusCode := http.StatusOK

	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	status := "ok"
	statusCode := http.StatusOK // 200
	// err := database.StatusCheck(ctx, h.DB , h.Log); 
	// h.Log.Error("debugging error", status, " err", err.Error(), " request context", r.Context() )

	if err := database.StatusCheck(ctx, h.DB , h.Log); err != nil {
	
		status = "db not ready"
		statusCode = http.StatusInternalServerError //500

		h.Log.Error("status", status, " StatusCode", statusCode, " err", err.Error() )
	}

	

	data :=  struct {
		Status string `json:"status"`
	}{
		Status: status,
	}

	if  err := response(w, statusCode, data); err != nil{
		h.Log.Error("readiness", "Error", err)
	}

	h.Log.Infow("readiness", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)

}

// Liveness returns simple status info if the service is alive. If the
// app is deployed to a Kubernetes cluster, it will also return pod, node, and
// namespace details via the Downward API. The Kubernetes environment variables
// need to be set within your Pod/Deployment manifest.
func (h Handlers) Liveness(w http.ResponseWriter, r *http.Request) {
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	data := struct {
		Status    string `json:"status,omitempty"`
		Build     string `json:"build,omitempty"`
		Host      string `json:"host,omitempty"`
		Pod       string `json:"pod,omitempty"`
		PodIP     string `json:"podIP,omitempty"`
		Node      string `json:"node,omitempty"`
		Namespace string `json:"namespace,omitempty"`
	}{
		Status:    "up",
		Build:     h.Build,
		Host:      host,
		Pod:       os.Getenv("KUBERNETES_PODNAME"),
		PodIP:     os.Getenv("KUBERNETES_NAMESPACE_POD_IP"),
		Node:      os.Getenv("KUBERNETES_NODENAME"),
		Namespace: os.Getenv("KUBERNETES_NAMESPACE"),
	}

	statusCode := http.StatusOK
	if err := response(w, statusCode, data); err != nil {
		h.Log.Errorw("liveness", "ERROR", err)
	}

	// THIS IS A FREE TIMER. WE COULD UPDATE THE METRIC GOROUTINE COUNT HERE.

	h.Log.Infow("liveness", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)

}

func response(w http.ResponseWriter, statusCode int, data any) error {

	// Convert the response value to JSON.
	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	// Set the content type and headers once we know marshaling has succeeded.
	w.Header().Set("Content-Type", "application/json")

	// Write the status code to the response.
	w.WriteHeader(statusCode)

	// Send the result back to the client.
	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}
