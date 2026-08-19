package calculator

import "errors"

// Sentinel errors returned by the domain. They are stable, comparable values
// (checked with errors.Is) so the transport layer can map each one to a
// machine-readable error code and HTTP status without depending on message
// text. Keeping them here — beside the logic that raises them — follows
// Separation of Concerns: the domain owns what "invalid" means; the transport
// owns how it is reported.
var (
	// ErrUnknownOperation is returned when the requested operation is not in
	// the registry.
	ErrUnknownOperation = errors.New("unknown operation")

	// ErrInvalidOperandCount is returned when the number of operands does not
	// match the operation's arity.
	ErrInvalidOperandCount = errors.New("invalid operand count")

	// ErrInvalidOperand is returned when an operand is not a finite number
	// (NaN or ±Inf).
	ErrInvalidOperand = errors.New("invalid operand")

	// ErrDivisionByZero is returned when a division has a zero divisor.
	ErrDivisionByZero = errors.New("division by zero")

	// ErrNegativeSquareRoot is returned when a square root is requested for a
	// negative operand.
	ErrNegativeSquareRoot = errors.New("square root of a negative number")

	// ErrUndefinedResult is returned when a computation produces a
	// non-finite result (overflow, or an otherwise undefined value such as a
	// negative base raised to a fractional exponent).
	ErrUndefinedResult = errors.New("undefined result")
)
