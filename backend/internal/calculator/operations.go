package calculator

import "math"

// Operation is a stable API identifier for a calculation.
type Operation string

const (
	OperationAdd      Operation = "add"
	OperationSubtract Operation = "subtract"
	OperationMultiply Operation = "multiply"
	OperationDivide   Operation = "divide"
)

// OperationInfo describes an operation exposed to API clients.
type OperationInfo struct {
	Name   Operation `json:"name"`
	Label  string    `json:"label"`
	Symbol string    `json:"symbol"`
	Arity  int       `json:"arity"`
}

type operationDefinition struct {
	info      OperationInfo
	calculate func([]float64) (float64, error)
}

// operationRegistry is the single source of truth for operation discovery and
// execution. Its order is also the order returned to API clients.
var operationRegistry = []operationDefinition{
	{
		info: OperationInfo{
			Name:   OperationAdd,
			Label:  "Addition",
			Symbol: "+",
			Arity:  2,
		},
		calculate: func(operands []float64) (float64, error) {
			return operands[0] + operands[1], nil
		},
	},
	{
		info: OperationInfo{
			Name:   OperationSubtract,
			Label:  "Subtraction",
			Symbol: "-",
			Arity:  2,
		},
		calculate: func(operands []float64) (float64, error) {
			return operands[0] - operands[1], nil
		},
	},
	{
		info: OperationInfo{
			Name:   OperationMultiply,
			Label:  "Multiplication",
			Symbol: "×",
			Arity:  2,
		},
		calculate: func(operands []float64) (float64, error) {
			return operands[0] * operands[1], nil
		},
	},
	{
		info: OperationInfo{
			Name:   OperationDivide,
			Label:  "Division",
			Symbol: "÷",
			Arity:  2,
		},
		calculate: func(operands []float64) (float64, error) {
			if operands[1] == 0 {
				return 0, ErrDivisionByZero
			}
			return operands[0] / operands[1], nil
		},
	},
}

// SupportedOperations returns a copy of the operation metadata.
func SupportedOperations() []OperationInfo {
	operations := make([]OperationInfo, 0, len(operationRegistry))
	for _, definition := range operationRegistry {
		operations = append(operations, definition.info)
	}
	return operations
}

func findOperation(operation Operation) (operationDefinition, bool) {
	for _, definition := range operationRegistry {
		if definition.info.Name == operation {
			return definition, true
		}
	}
	return operationDefinition{}, false
}

func isFinite(value float64) bool {
	return !math.IsNaN(value) && !math.IsInf(value, 0)
}
