package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/portella94/sezzle-challange-calculator-go-react/backend/internal/calculator"
)

// newTestServer returns a router backed by the real calculator service.
func newTestServer() http.Handler {
	return NewRouter(NewHandler(calculator.New()), []string{"*"})
}

func do(t *testing.T, srv http.Handler, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)
	return rec
}

func decodeError(t *testing.T, rec *httptest.ResponseRecorder) errorDetail {
	t.Helper()
	var body errorResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("could not decode error response %q: %v", rec.Body.String(), err)
	}
	return body.Error
}

func TestCalculate_Success(t *testing.T) {
	srv := newTestServer()
	rec := do(t, srv, http.MethodPost, "/api/v1/calculate", `{"operation":"add","operands":[2,3]}`)

	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200 (body %q)", rec.Code, rec.Body.String())
	}
	var resp calculateResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Result != 5 {
		t.Errorf("result = %v, want 5", resp.Result)
	}
	if resp.Operation != calculator.Add || len(resp.Operands) != 2 {
		t.Errorf("response did not echo request: %+v", resp)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Errorf("Content-Type = %q, want application/json", ct)
	}
}

func TestCalculate_DomainErrors(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		wantCode string
	}{
		{"division by zero", `{"operation":"divide","operands":[1,0]}`, "division_by_zero"},
		{"negative sqrt", `{"operation":"sqrt","operands":[-4]}`, "negative_square_root"},
		{"unknown operation", `{"operation":"modulo","operands":[1,2]}`, "unknown_operation"},
		{"wrong operand count", `{"operation":"add","operands":[1]}`, "invalid_operand_count"},
		{"undefined result", `{"operation":"power","operands":[-2,0.5]}`, "undefined_result"},
	}

	srv := newTestServer()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := do(t, srv, http.MethodPost, "/api/v1/calculate", tc.body)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400", rec.Code)
			}
			if got := decodeError(t, rec).Code; got != tc.wantCode {
				t.Errorf("error code = %q, want %q", got, tc.wantCode)
			}
		})
	}
}

func TestCalculate_RequestErrors(t *testing.T) {
	tests := []struct {
		name string
		body string
	}{
		{"malformed json", `{"operation":`},
		{"unknown field", `{"operation":"add","operands":[1,2],"extra":true}`},
		{"empty body", ``},
		{"trailing data", `{"operation":"add","operands":[1,2]}{}`},
	}

	srv := newTestServer()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := do(t, srv, http.MethodPost, "/api/v1/calculate", tc.body)
			if rec.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, want 400", rec.Code)
			}
			if got := decodeError(t, rec).Code; got != "invalid_request" {
				t.Errorf("error code = %q, want invalid_request", got)
			}
		})
	}
}

func TestHealth(t *testing.T) {
	srv := newTestServer()
	rec := do(t, srv, http.MethodGet, "/api/v1/healthz", "")
	if rec.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rec.Code)
	}
	var resp healthResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Status != "ok" || resp.Version == "" {
		t.Errorf("health = %+v, want status ok and non-empty version", resp)
	}
}

func TestMethodNotAllowed(t *testing.T) {
	srv := newTestServer()
	rec := do(t, srv, http.MethodGet, "/api/v1/calculate", "")
	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("status = %d, want 405", rec.Code)
	}
}

func TestCORSPreflight(t *testing.T) {
	srv := newTestServer()
	req := httptest.NewRequest(http.MethodOptions, "/api/v1/calculate", nil)
	req.Header.Set("Origin", "http://localhost:5173")
	rec := httptest.NewRecorder()
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want 204", rec.Code)
	}
	if got := rec.Header().Get("Access-Control-Allow-Origin"); got != "*" {
		t.Errorf("Allow-Origin = %q, want *", got)
	}
}

// stubCalculator lets us exercise the unexpected-error (500) branch that the
// real service never returns.
type stubCalculator struct{ err error }

func (s stubCalculator) Calculate(calculator.Operation, []float64) (float64, error) {
	return 0, s.err
}

func TestCalculate_InternalError(t *testing.T) {
	srv := NewRouter(NewHandler(stubCalculator{err: errors.New("boom")}), []string{"*"})
	rec := do(t, srv, http.MethodPost, "/api/v1/calculate", `{"operation":"add","operands":[1,2]}`)
	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want 500", rec.Code)
	}
	if got := decodeError(t, rec).Code; got != "internal_error" {
		t.Errorf("error code = %q, want internal_error", got)
	}
}
