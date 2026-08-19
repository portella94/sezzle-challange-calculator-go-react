package api

import "github.com/portella94/sezzle-challange-calculator-go-react/backend/internal/calculator"

// calculateRequest is the JSON body accepted by POST /api/v1/calculate.
type calculateRequest struct {
	Operation calculator.Operation `json:"operation"`
	Operands  []float64            `json:"operands"`
}

// calculateResponse is the JSON body returned on a successful calculation. It
// echoes the inputs alongside the result so the response is self-describing
// (Principle of Least Astonishment).
type calculateResponse struct {
	Operation calculator.Operation `json:"operation"`
	Operands  []float64            `json:"operands"`
	Result    float64              `json:"result"`
}

// healthResponse is the JSON body returned by GET /api/v1/healthz.
type healthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
}

// errorResponse is the single, consistent error envelope for every failure.
type errorResponse struct {
	Error errorDetail `json:"error"`
}

// errorDetail carries a stable machine-readable code and a human-readable
// message.
type errorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
