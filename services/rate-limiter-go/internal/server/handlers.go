package server

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// handleHealthz reports that the process is up.
func (s *Server) handleHealthz(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

// handleReadyz reports whether the server is accepting traffic. It returns
// 503 while shutting down so the pod is taken out of the Service first.
func (s *Server) handleReadyz(w http.ResponseWriter, _ *http.Request) {
	if !s.ready.Load() {
		writeJSON(w, http.StatusServiceUnavailable, map[string]string{"status": "not ready"})
		return
	}
	writeJSON(w, http.StatusOK, map[string]string{"status": "ready"})
}

type checkRequest struct {
	Key string `json:"key"`
}

type checkResponse struct {
	Allowed bool `json:"allowed"`
}

// handleCheck answers {"key": "alice"} with 200 {"allowed": true} or
// 429 {"allowed": false}.
func (s *Server) handleCheck(w http.ResponseWriter, r *http.Request) {
	var req checkRequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, 4096)).Decode(&req); err != nil || req.Key == "" {
		writeJSON(w, http.StatusBadRequest, map[string]string{"error": `body must be JSON like {"key": "alice"}`})
		return
	}

	if s.limiter.Allow(req.Key) {
		writeJSON(w, http.StatusOK, checkResponse{Allowed: true})
		return
	}
	writeJSON(w, http.StatusTooManyRequests, checkResponse{Allowed: false})
}

// recoverPanics turns a panic in a handler (such as the unimplemented limiter)
// into a 500 instead of a dropped connection.
func (s *Server) recoverPanics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				s.logger.Error("handler panic", "path", r.URL.Path, "panic", fmt.Sprint(v))
				writeJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprint(v)})
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}
