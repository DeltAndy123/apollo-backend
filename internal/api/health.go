package api

import (
	"encoding/json"
	"net/http"
	"os"
)

func (a *api) tokenAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		secret := os.Getenv("APOLLO_SECRET_TOKEN")
		if secret == "" {
			next.ServeHTTP(w, r)
			return
		}
		// Health check and bundle ID are always public so the client can
		// discover the server and verify credentials before registering.
		if r.RequestURI == "/v1/health" || r.RequestURI == "/v1/bundle_id" {
			next.ServeHTTP(w, r)
			return
		}
		if r.Header.Get("X-Apollo-Token") != secret {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

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
