package httpapi

import (
	"fmt"
	"strings"

	"go-calculator/internal/calculator"
)

type validationError struct {
	code    string
	message string
}

func validateCalculateRequest(request calculateRequest) *validationError {
	if request.Operation == nil || strings.TrimSpace(string(*request.Operation)) == "" {
		return &validationError{
			code:    "missing_operation",
			message: "operation is required",
		}
	}

	if !request.Operands.present || request.Operands.null {
		return &validationError{
			code:    "missing_operands",
			message: "operands are required",
		}
	}

	for _, operation := range calculator.SupportedOperations() {
		if operation.Name != *request.Operation {
			continue
		}

		if len(request.Operands.values) != operation.Arity {
			return &validationError{
				code: "invalid_operand_count",
				message: fmt.Sprintf(
					"operation %q requires %d operands",
					operation.Name,
					operation.Arity,
				),
			}
		}
		break
	}

	return nil
}
