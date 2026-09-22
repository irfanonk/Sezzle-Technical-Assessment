package httpapi

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go-calculator/internal/calculator"
)

type responseEnvelope struct {
	Data  json.RawMessage `json:"data"`
	Error *APIError       `json:"error"`
}

func TestCalculateEndpointSuccess(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		body       string
		wantResult float64
	}{
		{name: "addition", body: `{"operation":"add","operands":[2,3]}`, wantResult: 5},
		{name: "subtraction", body: `{"operation":"subtract","operands":[7,2.5]}`, wantResult: 4.5},
		{name: "multiplication", body: `{"operation":"multiply","operands":[-4,2]}`, wantResult: -8},
		{name: "division", body: `{"operation":"divide","operands":[7,2]}`, wantResult: 3.5},
		{name: "exponentiation", body: `{"operation":"exponentiate","operands":[2,8]}`, wantResult: 256},
		{name: "square root", body: `{"operation":"square_root","operands":[9]}`, wantResult: 3},
		{name: "percentage", body: `{"operation":"percentage","operands":[25]}`, wantResult: 0.25},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			response := performRequest(
				t,
				http.MethodPost,
				"/api/calculate",
				"application/json",
				test.body,
			)
			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
			}

			envelope := decodeEnvelope(t, response)
			var data calculateData
			if err := json.Unmarshal(envelope.Data, &data); err != nil {
				t.Fatalf("decode response data: %v", err)
			}
			if data.Result != test.wantResult {
				t.Errorf("result = %v, want %v", data.Result, test.wantResult)
			}
		})
	}
}

func TestCalculateEndpointFailures(t *testing.T) {
	t.Parallel()

	oversizedBody := `{"operation":"add","operands":[` + strings.Repeat("1,", 40_000) + `1]}`
	tests := []struct {
		name        string
		method      string
		path        string
		contentType string
		body        string
		wantStatus  int
		wantCode    string
		wantMessage string
		wantAllow   string
	}{
		{
			name:        "empty body",
			method:      http.MethodPost,
			path:        "/api/calculate",
			contentType: "application/json",
			wantStatus:  http.StatusBadRequest,
			wantCode:    "empty_body",
		},
		{
			name:        "malformed JSON",
			method:      http.MethodPost,
			path:        "/api/calculate",
			contentType: "application/json",
			body:        `{"operation":`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "malformed_json",
		},
		{
			name:        "unknown field",
			method:      http.MethodPost,
			path:        "/api/calculate",
			contentType: "application/json",
			body:        `{"operation":"add","operands":[1,2],"extra":true}`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "unknown_field",
		},
		{
			name:        "multiple JSON values",
			method:      http.MethodPost,
			path:        "/api/calculate",
			contentType: "application/json",
			body:        `{"operation":"add","operands":[1,2]} {}`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "malformed_json",
		},
		{
			name:        "missing operation",
			method:      http.MethodPost,
			path:        "/api/calculate",
			contentType: "application/json",
			body:        `{"operands":[1,2]}`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "missing_operation",
		},
		{
			name:        "null operation",
			method:      http.MethodPost,
			path:        "/api/calculate",
			contentType: "application/json",
			body:        `{"operation":null,"operands":[1,2]}`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "missing_operation",
		},
		{
			name:        "empty operation",
			method:      http.MethodPost,
			path:        "/api/calculate",
			contentType: "application/json",
			body:        `{"operation":" ","operands":[1,2]}`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "missing_operation",
		},
		{
			name:        "operation has wrong type",
			method:      http.MethodPost,
			path:        "/api/calculate",
			contentType: "application/json",
			body:        `{"operation":42,"operands":[1,2]}`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "invalid_field_type",
		},
		{
			name:        "missing operands",
			method:      http.MethodPost,
			path:        "/api/calculate",
			contentType: "application/json",
			body:        `{"operation":"add"}`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "missing_operands",
		},
		{
			name:        "null operands",
			method:      http.MethodPost,
			path:        "/api/calculate",
			contentType: "application/json",
			body:        `{"operation":"add","operands":null}`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "missing_operands",
		},
		{
			name:        "operands has wrong type",
			method:      http.MethodPost,
			path:        "/api/calculate",
			contentType: "application/json",
			body:        `{"operation":"add","operands":"1,2"}`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "invalid_field_type",
		},
		{
			name:        "operand has wrong type",
			method:      http.MethodPost,
			path:        "/api/calculate",
			contentType: "application/json",
			body:        `{"operation":"add","operands":[1,"two"]}`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "invalid_operand_type",
			wantMessage: "operand at index 1 must be a number",
		},
		{
			name:        "operand is null",
			method:      http.MethodPost,
			path:        "/api/calculate",
			contentType: "application/json",
			body:        `{"operation":"add","operands":[null,2]}`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "invalid_operand_type",
			wantMessage: "operand at index 0 must be a number",
		},
		{
			name:        "unsupported operation",
			method:      http.MethodPost,
			path:        "/api/calculate",
			contentType: "application/json",
			body:        `{"operation":"modulo","operands":[1,2]}`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "unsupported_operation",
		},
		{
			name:        "invalid operand count",
			method:      http.MethodPost,
			path:        "/api/calculate",
			contentType: "application/json",
			body:        `{"operation":"add","operands":[1]}`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "invalid_operand_count",
		},
		{
			name:        "empty operands",
			method:      http.MethodPost,
			path:        "/api/calculate",
			contentType: "application/json",
			body:        `{"operation":"add","operands":[]}`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "invalid_operand_count",
		},
		{
			name:        "invalid unary operand count",
			method:      http.MethodPost,
			path:        "/api/calculate",
			contentType: "application/json",
			body:        `{"operation":"square_root","operands":[9,2]}`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "invalid_operand_count",
		},
		{
			name:        "division by zero",
			method:      http.MethodPost,
			path:        "/api/calculate",
			contentType: "application/json",
			body:        `{"operation":"divide","operands":[1,0]}`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "division_by_zero",
		},
		{
			name:        "non-finite result",
			method:      http.MethodPost,
			path:        "/api/calculate",
			contentType: "application/json",
			body:        `{"operation":"multiply","operands":[1.7976931348623157e308,2]}`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "non_finite_value",
		},
		{
			name:        "negative square root",
			method:      http.MethodPost,
			path:        "/api/calculate",
			contentType: "application/json",
			body:        `{"operation":"square_root","operands":[-1]}`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "non_finite_value",
		},
		{
			name:        "exponentiation overflow",
			method:      http.MethodPost,
			path:        "/api/calculate",
			contentType: "application/json",
			body:        `{"operation":"exponentiate","operands":[1.7976931348623157e308,2]}`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "non_finite_value",
		},
		{
			name:        "non-finite numeric input",
			method:      http.MethodPost,
			path:        "/api/calculate",
			contentType: "application/json",
			body:        `{"operation":"add","operands":[1e999,2]}`,
			wantStatus:  http.StatusBadRequest,
			wantCode:    "invalid_operand_type",
			wantMessage: "operand at index 0 must be a number",
		},
		{
			name:        "unsupported media type",
			method:      http.MethodPost,
			path:        "/api/calculate",
			contentType: "text/plain",
			body:        `{"operation":"add","operands":[1,2]}`,
			wantStatus:  http.StatusUnsupportedMediaType,
			wantCode:    "unsupported_media_type",
		},
		{
			name:       "missing media type",
			method:     http.MethodPost,
			path:       "/api/calculate",
			body:       `{"operation":"add","operands":[1,2]}`,
			wantStatus: http.StatusUnsupportedMediaType,
			wantCode:   "unsupported_media_type",
		},
		{
			name:        "oversized request",
			method:      http.MethodPost,
			path:        "/api/calculate",
			contentType: "application/json",
			body:        oversizedBody,
			wantStatus:  http.StatusRequestEntityTooLarge,
			wantCode:    "request_too_large",
		},
		{
			name:       "unsupported method",
			method:     http.MethodGet,
			path:       "/api/calculate",
			wantStatus: http.StatusMethodNotAllowed,
			wantCode:   "method_not_allowed",
			wantAllow:  http.MethodPost,
		},
		{
			name:       "unknown endpoint",
			method:     http.MethodGet,
			path:       "/api/missing",
			wantStatus: http.StatusNotFound,
			wantCode:   "not_found",
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			response := performRequest(t, test.method, test.path, test.contentType, test.body)
			if response.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d; body = %s", response.Code, test.wantStatus, response.Body.String())
			}
			if allow := response.Header().Get("Allow"); allow != test.wantAllow {
				t.Errorf("Allow header = %q, want %q", allow, test.wantAllow)
			}

			envelope := decodeEnvelope(t, response)
			if envelope.Error == nil {
				t.Fatal("error response did not contain an error")
			}
			if envelope.Error.Code != test.wantCode {
				t.Errorf("error code = %q, want %q", envelope.Error.Code, test.wantCode)
			}
			if test.wantMessage != "" && envelope.Error.Message != test.wantMessage {
				t.Errorf("error message = %q, want %q", envelope.Error.Message, test.wantMessage)
			}
		})
	}
}

func TestOperationsEndpoint(t *testing.T) {
	t.Parallel()

	response := performRequest(t, http.MethodGet, "/api/operations", "", "")
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d; body = %s", response.Code, http.StatusOK, response.Body.String())
	}

	envelope := decodeEnvelope(t, response)
	var data operationsData
	if err := json.Unmarshal(envelope.Data, &data); err != nil {
		t.Fatalf("decode operations data: %v", err)
	}
	expectedArities := map[calculator.Operation]int{
		calculator.OperationAdd:      2,
		calculator.OperationSubtract: 2,
		calculator.OperationMultiply: 2,
		calculator.OperationDivide:   2,
		calculator.OperationExponent: 2,
		calculator.OperationSqrt:     1,
		calculator.OperationPercent:  1,
	}
	if len(data.Operations) != len(expectedArities) {
		t.Fatalf("operation count = %d, want %d", len(data.Operations), len(expectedArities))
	}

	seen := make(map[calculator.Operation]bool, len(data.Operations))
	for _, operation := range data.Operations {
		if seen[operation.Name] {
			t.Fatalf("operation %q was advertised more than once", operation.Name)
		}
		seen[operation.Name] = true
		expectedArity, ok := expectedArities[operation.Name]
		if !ok {
			t.Fatalf("unexpected advertised operation %q", operation.Name)
		}
		if operation.Arity != expectedArity {
			t.Errorf("operation %q arity = %d, want %d", operation.Name, operation.Arity, expectedArity)
		}

		operands := make([]float64, operation.Arity)
		for index := range operands {
			operands[index] = 1
		}

		payload, err := json.Marshal(map[string]any{
			"operation": operation.Name,
			"operands":  operands,
		})
		if err != nil {
			t.Fatalf("encode calculation for %q: %v", operation.Name, err)
		}

		calculationResponse := performRequest(
			t,
			http.MethodPost,
			"/api/calculate",
			"application/json",
			string(payload),
		)
		if calculationResponse.Code != http.StatusOK {
			t.Errorf(
				"advertised operation %q is not executable: status = %d, body = %s",
				operation.Name,
				calculationResponse.Code,
				calculationResponse.Body.String(),
			)
		}
	}
}

func TestOperationsEndpointRejectsUnsupportedMethod(t *testing.T) {
	t.Parallel()

	response := performRequest(t, http.MethodPost, "/api/operations", "application/json", `{}`)
	if response.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusMethodNotAllowed)
	}
	if allow := response.Header().Get("Allow"); allow != http.MethodGet {
		t.Errorf("Allow header = %q, want %q", allow, http.MethodGet)
	}

	envelope := decodeEnvelope(t, response)
	if envelope.Error == nil || envelope.Error.Code != "method_not_allowed" {
		t.Fatalf("error = %#v, want method_not_allowed", envelope.Error)
	}
}

func performRequest(
	t *testing.T,
	method string,
	path string,
	contentType string,
	body string,
) *httptest.ResponseRecorder {
	t.Helper()

	request := httptest.NewRequest(method, path, bytes.NewBufferString(body))
	if contentType != "" {
		request.Header.Set("Content-Type", contentType)
	}

	response := httptest.NewRecorder()
	NewHandler().ServeHTTP(response, request)
	return response
}

func decodeEnvelope(t *testing.T, response *httptest.ResponseRecorder) responseEnvelope {
	t.Helper()

	if contentType := response.Header().Get("Content-Type"); contentType != "application/json" {
		t.Errorf("Content-Type = %q, want application/json", contentType)
	}

	var fields map[string]json.RawMessage
	if err := json.Unmarshal(response.Body.Bytes(), &fields); err != nil {
		t.Fatalf("decode response envelope: %v", err)
	}
	if len(fields) != 2 || fields["data"] == nil || fields["error"] == nil {
		t.Fatalf("response must contain only data and error fields: %s", response.Body.String())
	}

	var envelope responseEnvelope
	if err := json.Unmarshal(response.Body.Bytes(), &envelope); err != nil {
		t.Fatalf("decode typed response envelope: %v", err)
	}

	if response.Code >= 200 && response.Code < 300 {
		if string(fields["data"]) == "null" || envelope.Error != nil {
			t.Errorf("success response must contain data and no error: %s", response.Body.String())
		}
	} else if string(fields["data"]) != "null" || envelope.Error == nil {
		t.Errorf("error response must contain an error and no data: %s", response.Body.String())
	}

	return envelope
}
