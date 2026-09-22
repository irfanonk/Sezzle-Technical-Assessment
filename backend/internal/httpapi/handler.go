package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"

	"go-calculator/internal/calculator"
)

const maxRequestBodyBytes int64 = 64 << 10

// Handler routes calculator API requests.
type Handler struct{}

// NewHandler creates the calculator API handler.
func NewHandler() http.Handler {
	return &Handler{}
}

func (handler *Handler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	switch request.URL.Path {
	case "/api/operations":
		handler.operations(writer, request)
	case "/api/calculate":
		handler.calculate(writer, request)
	default:
		writeError(writer, http.StatusNotFound, "not_found", "endpoint not found")
	}
}

type operationsData struct {
	Operations []calculator.OperationInfo `json:"operations"`
}

func (handler *Handler) operations(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		writer.Header().Set("Allow", http.MethodGet)
		writeError(writer, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}

	writeSuccess(writer, http.StatusOK, operationsData{
		Operations: calculator.SupportedOperations(),
	})
}

type calculateRequest struct {
	Operation calculator.Operation `json:"operation"`
	Operands  []float64            `json:"operands"`
}

type calculateData struct {
	Result float64 `json:"result"`
}

func (handler *Handler) calculate(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		writer.Header().Set("Allow", http.MethodPost)
		writeError(writer, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}

	mediaType, _, err := mime.ParseMediaType(request.Header.Get("Content-Type"))
	if err != nil || mediaType != "application/json" {
		writeError(
			writer,
			http.StatusUnsupportedMediaType,
			"unsupported_media_type",
			"Content-Type must be application/json",
		)
		return
	}

	var payload calculateRequest
	if !decodeRequest(writer, request, &payload) {
		return
	}

	result, err := calculator.Calculate(calculator.Calculation{
		Operation: payload.Operation,
		Operands:  payload.Operands,
	})
	if err != nil {
		writeCalculationError(writer, err)
		return
	}

	writeSuccess(writer, http.StatusOK, calculateData{Result: result})
}

func decodeRequest(writer http.ResponseWriter, request *http.Request, destination any) bool {
	request.Body = http.MaxBytesReader(writer, request.Body, maxRequestBodyBytes)

	decoder := json.NewDecoder(request.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(destination); err != nil {
		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError) {
			writeError(writer, http.StatusRequestEntityTooLarge, "request_too_large", "request body is too large")
			return false
		}
		writeError(writer, http.StatusBadRequest, "malformed_json", "request body must contain valid JSON")
		return false
	}

	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		var maxBytesError *http.MaxBytesError
		if errors.As(err, &maxBytesError) {
			writeError(writer, http.StatusRequestEntityTooLarge, "request_too_large", "request body is too large")
			return false
		}
		writeError(writer, http.StatusBadRequest, "malformed_json", "request body must contain one JSON object")
		return false
	}

	return true
}

func writeCalculationError(writer http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, calculator.ErrUnsupportedOperation):
		writeError(writer, http.StatusBadRequest, "unsupported_operation", "operation is not supported")
	case errors.Is(err, calculator.ErrInvalidOperandCount):
		writeError(writer, http.StatusBadRequest, "invalid_operand_count", "operation has an invalid number of operands")
	case errors.Is(err, calculator.ErrDivisionByZero):
		writeError(writer, http.StatusBadRequest, "division_by_zero", "cannot divide by zero")
	case errors.Is(err, calculator.ErrNonFiniteValue):
		writeError(writer, http.StatusBadRequest, "non_finite_value", "operands and result must be finite numbers")
	default:
		writeError(writer, http.StatusInternalServerError, "internal_error", "an unexpected error occurred")
	}
}
