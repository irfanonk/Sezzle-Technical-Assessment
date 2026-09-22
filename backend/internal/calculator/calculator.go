package calculator

import "errors"

var (
	ErrUnsupportedOperation = errors.New("unsupported operation")
	ErrInvalidOperandCount  = errors.New("invalid operand count")
	ErrDivisionByZero       = errors.New("division by zero")
	ErrNonFiniteValue       = errors.New("non-finite value")
)

// Calculation contains all input needed to perform an operation.
type Calculation struct {
	Operation Operation
	Operands  []float64
}

// Calculate validates and executes a calculation.
func Calculate(calculation Calculation) (float64, error) {
	definition, ok := findOperation(calculation.Operation)
	if !ok {
		return 0, ErrUnsupportedOperation
	}

	if len(calculation.Operands) != definition.info.Arity {
		return 0, ErrInvalidOperandCount
	}

	for _, operand := range calculation.Operands {
		if !isFinite(operand) {
			return 0, ErrNonFiniteValue
		}
	}

	result, err := definition.calculate(calculation.Operands)
	if err != nil {
		return 0, err
	}
	if !isFinite(result) {
		return 0, ErrNonFiniteValue
	}

	return result, nil
}
