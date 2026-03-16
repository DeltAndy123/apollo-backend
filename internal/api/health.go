package api

import (
	"encoding/json"
	"net/http"
	"os"
)

func (a *api) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status": "available",
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}

func (a *api) bundleIDHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"bundle_id": os.Getenv("APPLE_BUNDLE_ID"),
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}
