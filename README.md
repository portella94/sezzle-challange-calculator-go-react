# Full-Stack Calculator (Go + React/TypeScript)

This application has two parts: a Go REST microservice and a React
(TypeScript) frontend. The frontend collects two operands and an operation.
The frontend sends this data to the backend. The backend does all arithmetic,
data checks, and error handling.


## Software components

| Layer    | Software                                                    |
| -------- | ------------------------------------------------------------ |
| Frontend | React 19, TypeScript, Vite 8.2.1, MUI 9.3.1, Vitest 4.1.10   |
| Backend  | Go 1.26.6, standard library `net/http` (no framework)        |
| Tools    | Node 24 LTS, Docker, Docker Compose                           |

## Repository structure

```
backend/    Go microservice (domain, transport, config). Refer to backend/README.md.
frontend/   React app (Vite + MUI). Refer to frontend/README.md.
docker-compose.yml   Runs both services together.
```

## How to start the application with Docker

Do this to start the full application:

```bash
docker compose up --build
```

Then, open this address in a browser: http://localhost:8080.

nginx serves the built frontend on port 8080. nginx also sends `/api`
requests to the Go backend. Because of this, the frontend and the backend
use the same origin.

## How to develop the application on your local computer

Open two terminal windows. Run one service in each window.

### Backend

Use Go 1.26.x. The Go toolchain gets this version automatically from the
`go.mod` file.

```bash
cd backend
go run ./cmd/server
```

This command starts the backend at this address: http://localhost:8080.

### Frontend

```bash
cd frontend
npm install
npm run dev
```

This command starts the frontend at this address: http://localhost:5173.
The Vite dev server sends API calls to the backend. Because of this, you do
not need to configure CORS in the development environment.

## API reference

The base path is `/api/v1`. All requests and responses use the JSON format.

### `POST /api/v1/calculate`

This function does one arithmetic operation.

**Request body**

| Field       | Type      | Description                                              |
| ----------- | --------- | ---------------------------------------------------------- |
| `operation` | string    | Give one of the operations in the table below.              |
| `operands`  | number[]  | Give the operands. The number of operands must agree with the arity of the operation. |

**Operations**

| Operation    | Arity | Operands       | Function                     |
| ------------ | :---: | -------------- | ----------------------------- |
| `add`        |   2   | `[a, b]`       | Calculates `a + b`             |
| `subtract`   |   2   | `[a, b]`       | Calculates `a - b`             |
| `multiply`   |   2   | `[a, b]`       | Calculates `a × b`             |
| `divide`     |   2   | `[a, b]`       | Calculates `a ÷ b`. `b` must not be `0`. |
| `power`      |   2   | `[base, exp]`  | Calculates `base ^ exp`        |
| `sqrt`       |   1   | `[a]`          | Calculates `√a`. `a` must not be less than `0`. |
| `percentage` |   2   | `[a, b]`       | Calculates `a% of b`, which is `(a / 100) × b` |

**Success response — `200 OK`**

```json
{ "operation": "divide", "operands": [10, 2], "result": 5 }
```

**Error response — `400 Bad Request`**

```json
{ "error": { "code": "division_by_zero", "message": "division by zero" } }
```

**Error codes**

| Code                    | Cause                                                    |
| ----------------------- | ----------------------------------------------------------- |
| `invalid_request`       | The JSON data is not correct. Or, it has an unknown field. Or, the body is empty. |
| `unknown_operation`     | The system does not support this operation.                 |
| `invalid_operand_count` | The number of operands is not correct for this operation.   |
| `invalid_operand`       | An operand is not a finite number.                          |
| `division_by_zero`      | The divisor is `0`.                                          |
| `negative_square_root`  | The operand for the square-root operation is less than `0`. |
| `undefined_result`      | The result is not a finite number. This can occur from an overflow or from an undefined operation. |
| `internal_error`        | The server had an unexpected error (`500`).                  |

### `GET /api/v1/healthz`

```json
{ "status": "ok", "version": "0.0.1" }
```

### Examples

```bash
# Addition
curl -X POST localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"add","operands":[2,3]}'
# Result: {"operation":"add","operands":[2,3],"result":5}

# Square root
curl -X POST localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"sqrt","operands":[144]}'
# Result: {"operation":"sqrt","operands":[144],"result":12}

# Percentage (50% of 200)
curl -X POST localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"percentage","operands":[50,200]}'
# Result: {"operation":"percentage","operands":[50,200],"result":100}

# Division by zero (error)
curl -X POST localhost:8080/api/v1/calculate \
  -H 'Content-Type: application/json' \
  -d '{"operation":"divide","operands":[1,0]}'
# Result: 400 {"error":{"code":"division_by_zero","message":"division by zero"}}
```

## Design decisions

- **Layered backend.** The `calculator` package contains only arithmetic
  logic. This package has no knowledge of HTTP. The `api` package manages
  transport: it decodes requests, validates data, maps errors, and encodes
  responses. The `cmd/server` package is the composition root. You can test
  each layer by itself.
- **Operation registry.** The operations exist in one map. Each entry
  contains the arity and the function for one operation. To add an
  operation, add one entry to this map. You do not need to change the
  handler or the validation logic. The validation logic (arity checks,
  finiteness checks) exists once, for all operations.
- **One `/calculate` endpoint.** The API uses one request schema and one
  validation path, not a separate handler for each operation. A two-operand
  form matches this design. The design does not need an expression parser.
- **Fail-fast behavior.** The backend rejects these conditions before it
  starts a calculation: unknown JSON fields, unknown operations, an
  incorrect number of operands, and non-finite input values. The frontend
  checks that values are present and have a correct numeric format before
  it sends a request. However, the backend remains the authority for all
  validation.
- **Typed, stable errors.** Each domain error maps to a machine-readable
  code and to the correct HTTP status. Client applications must not parse
  the text of the error message.
- **Dependency inversion.** The HTTP handler depends on a `Calculator`
  interface, not on a concrete type. This keeps the transport layer
  independent. This also makes the handler easy to test with a stub.
- **Single source of truth in the frontend.** The file
  `config/operations.ts` declares the arity and the labels for each
  operation. The form uses this data to determine the number of input
  fields to show. This matches the backend contract.
- **Standard library only (backend).** The backend does not use a web
  framework. The `net/http` package, with the method-aware routing in Go, 
  is sufficient. This keeps the number of dependencies low.

## Tests and coverage

Use these commands to run tests:

```bash
# Backend tests
cd backend && go test ./... -cover

# Frontend tests
cd frontend && npm run coverage
```

Refer to [docs/COVERAGE.md](docs/COVERAGE.md) for the current coverage
numbers. The latest results are:

- **Backend:** The domain and config packages have 100% coverage. The
  transport package has approximately 97% coverage. Only the `main`
  function is not covered.
- **Frontend:** The API client, the hook, the utility functions, and the
  components have approximately 96% statement coverage together.

Refer to [docs/prompt-calculator-repo-review.md](docs/prompt-calculator-repo-review.md) 
for the prompt for the engineering-quality review.

Refer to [docs/prompt-calculator-repo-tests-plan.md](docs/prompt-calculator-repo-tests-plan.md) 
for the prompt for the test-suite plan.
