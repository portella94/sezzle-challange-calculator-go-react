package calculator

import (
	"errors"
	"math"
	"testing"
)

func TestCalculate_Success(t *testing.T) {
	tests := []struct {
		name     string
		op       Operation
		operands []float64
		want     float64
	}{
		{"add", Add, []float64{2, 3}, 5},
		{"add negatives", Add, []float64{-2, -3}, -5},
		{"subtract", Subtract, []float64{10, 4}, 6},
		{"multiply", Multiply, []float64{6, 7}, 42},
		{"multiply by zero", Multiply, []float64{6, 0}, 0},
		{"divide", Divide, []float64{10, 2}, 5},
		{"divide fractional", Divide, []float64{1, 4}, 0.25},
		{"power", Power, []float64{2, 10}, 1024},
		{"power zero exponent", Power, []float64{5, 0}, 1},
		{"power negative exponent", Power, []float64{2, -1}, 0.5},
		{"sqrt", SquareRoot, []float64{144}, 12},
		{"sqrt of zero", SquareRoot, []float64{0}, 0},
		{"percentage", Percentage, []float64{50, 200}, 100}, // 50% of 200
		{"percentage zero", Percentage, []float64{0, 200}, 0},
	}

	svc := New()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, err := svc.Calculate(tc.op, tc.operands)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if math.Abs(got-tc.want) > 1e-9 {
				t.Errorf("Calculate(%q, %v) = %v, want %v", tc.op, tc.operands, got, tc.want)
			}
		})
	}
}

func TestCalculate_Errors(t *testing.T) {
	tests := []struct {
		name     string
		op       Operation
		operands []float64
		wantErr  error
	}{
		{"unknown operation", Operation("modulo"), []float64{1, 2}, ErrUnknownOperation},
		{"too few operands", Add, []float64{1}, ErrInvalidOperandCount},
		{"too many operands", Add, []float64{1, 2, 3}, ErrInvalidOperandCount},
		{"sqrt wrong arity", SquareRoot, []float64{1, 2}, ErrInvalidOperandCount},
		{"nan operand", Add, []float64{math.NaN(), 1}, ErrInvalidOperand},
		{"inf operand", Add, []float64{math.Inf(1), 1}, ErrInvalidOperand},
		{"division by zero", Divide, []float64{1, 0}, ErrDivisionByZero},
		{"negative sqrt", SquareRoot, []float64{-1}, ErrNegativeSquareRoot},
		{"undefined power result", Power, []float64{-2, 0.5}, ErrUndefinedResult},
		{"overflow to undefined", Power, []float64{math.MaxFloat64, 2}, ErrUndefinedResult},
	}

	svc := New()
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			_, err := svc.Calculate(tc.op, tc.operands)
			if !errors.Is(err, tc.wantErr) {
				t.Errorf("Calculate(%q, %v) error = %v, want %v", tc.op, tc.operands, err, tc.wantErr)
			}
		})
	}
}

func TestArity(t *testing.T) {
	if n, ok := Arity(SquareRoot); !ok || n != 1 {
		t.Errorf("Arity(sqrt) = %d, %v; want 1, true", n, ok)
	}
	if n, ok := Arity(Add); !ok || n != 2 {
		t.Errorf("Arity(add) = %d, %v; want 2, true", n, ok)
	}
	if _, ok := Arity(Operation("nope")); ok {
		t.Errorf("Arity(nope) ok = true; want false")
	}
}

func TestSupportedOperations(t *testing.T) {
	got := SupportedOperations()
	if len(got) != 7 {
		t.Fatalf("SupportedOperations() returned %d ops, want 7", len(got))
	}
	// Verify stable sorted order and that every returned op is usable.
	for i := 1; i < len(got); i++ {
		if got[i-1] >= got[i] {
			t.Errorf("SupportedOperations() not sorted: %v", got)
			break
		}
	}
	for _, op := range got {
		if _, ok := Arity(op); !ok {
			t.Errorf("SupportedOperations() returned unknown op %q", op)
		}
	}
}
