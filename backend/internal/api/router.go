package api

import "net/http"

// NewRouter builds the fully-wired HTTP handler for the service: the route
// table plus the middleware chain (outermost first). Method-aware patterns
// (Go 1.22+) mean unmatched methods return 405 automatically.
func NewRouter(h *Handler, allowedOrigins []string) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/calculate", h.Calculate)
	mux.HandleFunc("GET /api/v1/healthz", h.Health)

	return chain(mux,
		recoverer,
		logging,
		cors(allowedOrigins),
	)
}
