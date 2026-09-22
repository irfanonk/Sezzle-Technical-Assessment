# Go Calculator API

A small calculator HTTP service built with the Go standard library.

## Requirements

- Go 1.22 or newer

## Run

```sh
go run ./cmd/server
```

The server listens on port `8080` by default. Set `PORT` to use another port:

```sh
PORT=3000 go run ./cmd/server
```

The server applies read, write, header, and idle timeouts and shuts down
gracefully on `SIGINT` or `SIGTERM`.

## Test

```sh
go test ./...
```

## API response format

Every endpoint returns JSON using the same envelope. A response contains either
`data` or `error`, never both.

Success:

```json
{
  "data": {
    "result": 5
  },
  "error": null
}
```

Failure:

```json
{
  "data": null,
  "error": {
    "code": "division_by_zero",
    "message": "cannot divide by zero"
  }
}
```

Clients should use the stable error `code` for program logic. The `message` is
intended for display or diagnostics.

## List supported operations

```http
GET /api/operations
```

Example:

```sh
curl http://localhost:8080/api/operations
```

The response includes each operation's stable name, label, symbol, and required
operand count:

```json
{
  "data": {
    "operations": [
      {
        "name": "add",
        "label": "Addition",
        "symbol": "+",
        "arity": 2
      }
    ]
  },
  "error": null
}
```

The operation list and calculation dispatch use the same internal registry, so
an advertised operation is always executable.

## Calculate

```http
POST /api/calculate
Content-Type: application/json
```

Request:

```json
{
  "operation": "divide",
  "operands": [10, 2]
}
```

Example:

```sh
curl \
  -H 'Content-Type: application/json' \
  -d '{"operation":"divide","operands":[10,2]}' \
  http://localhost:8080/api/calculate
```

Supported operation names are `add`, `subtract`, `multiply`, and `divide`.
Each currently requires exactly two operands.

## Error codes

- `empty_body`: the request body is empty
- `malformed_json`: invalid JSON or multiple JSON values
- `unknown_field`: the request contains a field outside the API contract
- `invalid_field_type`: a request field has the wrong JSON type
- `missing_operation`: the operation is missing, null, empty, or whitespace
- `missing_operands`: operands are missing or null
- `invalid_operand_type`: an operand is not a finite JSON number; the message
  identifies its zero-based array index
- `unsupported_media_type`: the calculation request is not
  `application/json`
- `request_too_large`: the request body exceeds 64 KiB
- `unsupported_operation`: the operation is not in the registry
- `invalid_operand_count`: the number of operands does not match the operation
- `division_by_zero`: the divisor is zero or negative zero
- `non_finite_value`: an operand or calculated result is NaN or infinity
- `method_not_allowed`: the endpoint does not support the HTTP method
- `not_found`: the endpoint does not exist
- `internal_error`: an unexpected server error

Invalid requests return `400`, oversized bodies return `413`, unsupported
media types return `415`, and unsupported methods return `405`.

## Design

Calculation rules live in `internal/calculator` and have no dependency on HTTP.
The HTTP package keeps JSON decoding separate from request contract validation
in `validators.go`, then maps failures to the shared response envelope. Inputs
are grouped in a `Calculation` struct rather than passed as positional function
parameters.

The service uses `float64`, which is suitable for a general calculator but has
normal IEEE-754 precision behavior. It rejects NaN, infinity, and results that
overflow to infinity.
