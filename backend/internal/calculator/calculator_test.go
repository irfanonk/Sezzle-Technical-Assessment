package calculator

import (
	"errors"
	"math"
	"testing"
)

func TestCalculate(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name        string
		calculation Calculation
		want        float64
		wantErr     error
	}{
		{
			name:        "addition",
			calculation: Calculation{Operation: OperationAdd, Operands: []float64{2, 3}},
			want:        5,
		},
		{
			name:        "subtraction",
			calculation: Calculation{Operation: OperationSubtract, Operands: []float64{7, 2.5}},
			want:        4.5,
		},
		{
			name:        "multiplication",
			calculation: Calculation{Operation: OperationMultiply, Operands: []float64{-4, 2}},
			want:        -8,
		},
		{
			name:        "division",
			calculation: Calculation{Operation: OperationDivide, Operands: []float64{7, 2}},
			want:        3.5,
		},
		{
			name:        "exponentiation",
			calculation: Calculation{Operation: OperationExponent, Operands: []float64{2, 8}},
			want:        256,
		},
		{
			name:        "square root",
			calculation: Calculation{Operation: OperationSqrt, Operands: []float64{9}},
			want:        3,
		},
		{
			name:        "percentage",
			calculation: Calculation{Operation: OperationPercent, Operands: []float64{25}},
			want:        0.25,
		},
		{
			name:        "unsupported operation",
			calculation: Calculation{Operation: "modulo", Operands: []float64{7, 2}},
			wantErr:     ErrUnsupportedOperation,
		},
		{
			name:        "too few operands",
			calculation: Calculation{Operation: OperationAdd, Operands: []float64{2}},
			wantErr:     ErrInvalidOperandCount,
		},
		{
			name:        "too many operands",
			calculation: Calculation{Operation: OperationAdd, Operands: []float64{2, 3, 4}},
			wantErr:     ErrInvalidOperandCount,
		},
		{
			name:        "too many operands for unary operation",
			calculation: Calculation{Operation: OperationSqrt, Operands: []float64{9, 2}},
			wantErr:     ErrInvalidOperandCount,
		},
		{
			name:        "division by zero",
			calculation: Calculation{Operation: OperationDivide, Operands: []float64{7, 0}},
			wantErr:     ErrDivisionByZero,
		},
		{
			name: "division by negative zero",
			calculation: Calculation{
				Operation: OperationDivide,
				Operands:  []float64{7, math.Copysign(0, -1)},
			},
			wantErr: ErrDivisionByZero,
		},
		{
			name:        "NaN operand",
			calculation: Calculation{Operation: OperationAdd, Operands: []float64{math.NaN(), 1}},
			wantErr:     ErrNonFiniteValue,
		},
		{
			name:        "infinite operand",
			calculation: Calculation{Operation: OperationAdd, Operands: []float64{math.Inf(1), 1}},
			wantErr:     ErrNonFiniteValue,
		},
		{
			name:        "infinite result",
			calculation: Calculation{Operation: OperationMultiply, Operands: []float64{math.MaxFloat64, 2}},
			wantErr:     ErrNonFiniteValue,
		},
		{
			name:        "negative square root",
			calculation: Calculation{Operation: OperationSqrt, Operands: []float64{-1}},
			wantErr:     ErrNonFiniteValue,
		},
		{
			name:        "exponentiation overflow",
			calculation: Calculation{Operation: OperationExponent, Operands: []float64{math.MaxFloat64, 2}},
			wantErr:     ErrNonFiniteValue,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			got, err := Calculate(test.calculation)
			if !errors.Is(err, test.wantErr) {
				t.Fatalf("Calculate() error = %v, want %v", err, test.wantErr)
			}
			if test.wantErr == nil && got != test.want {
				t.Errorf("Calculate() = %v, want %v", got, test.want)
			}
		})
	}
}

func TestSupportedOperationsReturnsCopy(t *testing.T) {
	t.Parallel()

	first := SupportedOperations()
	if len(first) == 0 {
		t.Fatal("SupportedOperations() returned no operations")
	}

	first[0].Name = "changed"
	second := SupportedOperations()
	if second[0].Name == "changed" {
		t.Fatal("SupportedOperations() exposed mutable registry state")
	}
}
