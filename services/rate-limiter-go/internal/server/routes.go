package server

import "net/http"

// routes lists every endpoint the service exposes.
func (s *Server) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", s.handleHealthz) // liveness
	mux.HandleFunc("GET /readyz", s.handleReadyz)   // readiness
	mux.HandleFunc("POST /v1/check", s.handleCheck) // rate limit check

	return s.recoverPanics(mux)
}
