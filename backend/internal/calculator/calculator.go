// Package calculator holds the pure arithmetic domain of the service.
//
// It has zero knowledge of HTTP, JSON, or any transport concern (Separation of
// Concerns) and is therefore trivially unit-testable. Operations are declared
// in a registry keyed by Operation; each entry carries its arity and
// implementation, so supporting a new operation is a single map entry and
// requires no changes to Calculate or the transport layer (Open/Closed
// Principle).
package calculator

import (
	"math"
	"sort"
)

// Operation is the identifier of a supported arithmetic operation. It is the
// wire value clients send (e.g. "add").
type Operation string

// Supported operations.
const (
	Add        Operation = "add"
	Subtract   Operation = "subtract"
	Multiply   Operation = "multiply"
	Divide     Operation = "divide"
	Power      Operation = "power"      // operands[0] ^ operands[1]
	SquareRoot Operation = "sqrt"       // √operands[0]
	Percentage Operation = "percentage" // operands[0]% of operands[1] = a/100*b
)

// operation is the internal metadata + implementation of an Operation.
type operation struct {
	arity int
	apply func(operands []float64) (float64, error)
}

// registry is the single source of truth for supported operations. Add an
// entry here to support a new operation — nothing else needs to change.
var registry = map[Operation]operation{
	Add:      {arity: 2, apply: func(o []float64) (float64, error) { return o[0] + o[1], nil }},
	Subtract: {arity: 2, apply: func(o []float64) (float64, error) { return o[0] - o[1], nil }},
	Multiply: {arity: 2, apply: func(o []float64) (float64, error) { return o[0] * o[1], nil }},
	Divide: {arity: 2, apply: func(o []float64) (float64, error) {
		if o[1] == 0 {
			return 0, ErrDivisionByZero
		}
		return o[0] / o[1], nil
	}},
	Power: {arity: 2, apply: func(o []float64) (float64, error) {
		// A non-finite result (e.g. a negative base with a fractional
		// exponent yields NaN) is caught by the post-computation guard in
		// Calculate and reported as ErrUndefinedResult.
		return math.Pow(o[0], o[1]), nil
	}},
	SquareRoot: {arity: 1, apply: func(o []float64) (float64, error) {
		if o[0] < 0 {
			return 0, ErrNegativeSquareRoot
		}
		return math.Sqrt(o[0]), nil
	}},
	Percentage: {arity: 2, apply: func(o []float64) (float64, error) {
		return (o[0] / 100) * o[1], nil
	}},
}

// Service performs arithmetic operations. It carries no state; the struct
// exists so callers can depend on it through a narrow interface and inject it
// at their composition root (Dependency Inversion).
type Service struct{}

// New returns a ready-to-use calculator Service.
func New() *Service { return &Service{} }

// Calculate validates the request and applies the requested operation.
//
// Validation is fail-fast and ordered: unknown operation, wrong operand count,
// non-finite operands, then the operation's own domain rules, then a final
// guard rejecting non-finite results. Every failure is a stable sentinel error
// (see errors.go) so the transport can map it without inspecting message text.
func (s *Service) Calculate(op Operation, operands []float64) (float64, error) {
	def, ok := registry[op]
	if !ok {
		return 0, ErrUnknownOperation
	}
	if len(operands) != def.arity {
		return 0, ErrInvalidOperandCount
	}
	for _, v := range operands {
		if math.IsNaN(v) || math.IsInf(v, 0) {
			return 0, ErrInvalidOperand
		}
	}

	result, err := def.apply(operands)
	if err != nil {
		return 0, err
	}
	if math.IsNaN(result) || math.IsInf(result, 0) {
		return 0, ErrUndefinedResult
	}
	return result, nil
}

// Arity returns the number of operands the operation expects and whether the
// operation is supported.
func Arity(op Operation) (int, bool) {
	def, ok := registry[op]
	return def.arity, ok
}

// SupportedOperations returns the supported operation identifiers in a stable,
// sorted order. Useful for documentation and for keeping clients in sync.
func SupportedOperations() []Operation {
	ops := make([]Operation, 0, len(registry))
	for op := range registry {
		ops = append(ops, op)
	}
	sort.Slice(ops, func(i, j int) bool { return ops[i] < ops[j] })
	return ops
}
