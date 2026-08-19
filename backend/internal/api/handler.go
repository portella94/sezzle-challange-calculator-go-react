package api

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/portella94/sezzle-challange-calculator-go-react/backend/internal/calculator"
	"github.com/portella94/sezzle-challange-calculator-go-react/backend/internal/version"
)

// maxBodyBytes caps request bodies to a small size; a calculation payload is
// tiny, so anything larger is rejected outright (fail fast, defensive).
const maxBodyBytes = 1 << 16 // 64 KiB

// Calculator is the narrow behaviour the handler depends on. Declaring the
// interface here — at the consumer — lets the transport layer be tested with a
// stub and keeps it decoupled from the concrete domain type (Dependency
// Inversion, Interface Segregation).
type Calculator interface {
	Calculate(op calculator.Operation, operands []float64) (float64, error)
}

// Handler serves the calculator HTTP endpoints.
type Handler struct {
	calc Calculator
}

// NewHandler wires a Handler to a Calculator implementation.
func NewHandler(calc Calculator) *Handler {
	return &Handler{calc: calc}
}

// Calculate handles POST /api/v1/calculate.
func (h *Handler) Calculate(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)

	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields() // fail fast on unexpected fields

	var req calculateRequest
	if err := dec.Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid_request", decodeErrorMessage(err))
		return
	}
	// Reject trailing data after the JSON object (e.g. two concatenated bodies).
	if dec.More() {
		writeError(w, http.StatusBadRequest, "invalid_request", "request body must contain a single JSON object")
		return
	}

	result, err := h.calc.Calculate(req.Operation, req.Operands)
	if err != nil {
		status, code := mapDomainError(err)
		writeError(w, status, code, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, calculateResponse{
		Operation: req.Operation,
		Operands:  req.Operands,
		Result:    result,
	})
}

// Health handles GET /api/v1/healthz.
func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, healthResponse{Status: "ok", Version: version.Version})
}

// mapDomainError translates a domain sentinel error into an HTTP status and a
// stable, machine-readable error code. All domain errors stem from client
// input and are therefore 400s.
func mapDomainError(err error) (status int, code string) {
	switch {
	case errors.Is(err, calculator.ErrUnknownOperation):
		return http.StatusBadRequest, "unknown_operation"
	case errors.Is(err, calculator.ErrInvalidOperandCount):
		return http.StatusBadRequest, "invalid_operand_count"
	case errors.Is(err, calculator.ErrInvalidOperand):
		return http.StatusBadRequest, "invalid_operand"
	case errors.Is(err, calculator.ErrDivisionByZero):
		return http.StatusBadRequest, "division_by_zero"
	case errors.Is(err, calculator.ErrNegativeSquareRoot):
		return http.StatusBadRequest, "negative_square_root"
	case errors.Is(err, calculator.ErrUndefinedResult):
		return http.StatusBadRequest, "undefined_result"
	default:
		return http.StatusInternalServerError, "internal_error"
	}
}

// decodeErrorMessage returns a concise, safe message for a JSON decode failure.
func decodeErrorMessage(err error) string {
	switch {
	case errors.Is(err, io.EOF):
		return "request body is empty"
	default:
		return "request body is not valid JSON matching the expected schema"
	}
}

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		// The status/headers are already sent; nothing to do but log.
		log.Printf("api: failed to encode response: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, code, message string) {
	writeJSON(w, status, errorResponse{Error: errorDetail{Code: code, Message: message}})
}
