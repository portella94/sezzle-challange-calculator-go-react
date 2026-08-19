# Prompt: Write the Test Suite for an Existing Full-Stack Calculator Repo

I have an existing repository for a full-stack calculator (Go backend +
React/TypeScript frontend). Do not change any production code. Write the
test suite only.

## Context (existing code, do not modify)

```
backend/
  cmd/server/main.go
  internal/calculator/calculator.go   pure arithmetic domain + operation registry
  internal/calculator/errors.go       sentinel errors (errors.New), checked via errors.Is
  internal/api/dto.go                 request/response/error JSON structs
  internal/api/handler.go             decode → validate → calculate → map errors → encode
  internal/api/middleware.go          e.g. logging/recovery
  internal/api/router.go
  internal/config/config.go           reads config (e.g. PORT) from env with defaults
  internal/version/version.go
  go.mod                              go 1.26

frontend/src/
  api/client.ts        postJson<TReq,TRes>() fetch wrapper; CalculatorApiError (has `code`)
  api/calculator.ts     calculate() → POST /v1/calculate
  config/operations.ts  arity/labels/hints per operation (source of truth)
  hooks/useCalculator.ts state + submit logic
  components/CalculatorForm.tsx, OperationSelect.tsx, ResultDisplay.tsx
  utils/format.ts
  types.ts
  App.tsx
```

API contract: `POST /api/v1/calculate` takes `{ operation, operands: number[] }`
and returns `{ operation, operands, result }` on `200`, or
`{ error: { code, message } }` on `400`/`500`. Supported operations and
their arity: `add`(2), `subtract`(2), `multiply`(2), `divide`(2, `b`≠0),
`power`(2), `sqrt`(1, operand ≥ 0), `percentage`(2). Error codes:
`invalid_request`, `unknown_operation`, `invalid_operand_count`,
`invalid_operand`, `division_by_zero`, `negative_square_root`,
`undefined_result`, `internal_error`. `GET /api/v1/healthz` returns
`{ status, version }`.

## Task

Generate a complete, high-coverage test suite for both stacks, matching the
patterns and conventions already used in this codebase (table-driven Go
tests, `httptest` for HTTP, React Testing Library + Vitest for the
frontend). Do not restructure or refactor production code to make it more
testable — the code is already designed for testability (calculator has no
HTTP knowledge, handler depends on a `Calculator` interface, single
operation registry).

### Backend (Go, `testing` + `net/http/httptest`, table-driven style)

1. **`internal/calculator/calculator_test.go`**
   - `TestCalculate_Success`: table-driven, one case per operation, plus
     edge cases (multiply by zero, divide to a fraction, power with a zero
     or negative exponent, sqrt of zero, percentage of zero). Compare
     floats with a small epsilon, not `==`.
   - `TestCalculate_Errors`: unknown operation; too few/too many operands
     for a 2-arity op; wrong arity for `sqrt`; `NaN` operand; `+Inf`/`-Inf`
     operand; division by zero; negative-number square root; a case that
     produces a non-finite result (e.g. negative base with a fractional
     exponent, and an overflow case) mapped to `ErrUndefinedResult`. Assert
     with `errors.Is`, never by comparing error strings.
   - `TestArity`: known operations return the correct arity and `ok=true`;
     an unknown operation returns `ok=false`.
   - `TestSupportedOperations`: returns exactly the expected count, in
     stable sorted order, and every returned operation is a valid key via
     `Arity`.

2. **`internal/api/handler_test.go`** (using `httptest.NewRecorder` /
   `httptest.NewRequest`, and a stub implementing the handler's
   `Calculator` interface so domain logic is not re-tested here)
   - Happy path: valid request for each arity (1 and 2 operands) returns
     `200` with the correct JSON body.
   - Malformed JSON body → `400 invalid_request`.
   - Unknown JSON field in the body (decoder must reject it) → `400 invalid_request`.
   - Empty body → `400 invalid_request`.
   - Each domain sentinel error surfaced by the stub calculator maps to its
     documented HTTP status and error code (`unknown_operation`,
     `invalid_operand_count`, `invalid_operand`, `division_by_zero`,
     `negative_square_root`, `undefined_result`).
   - An unexpected/unmapped error from the calculator maps to `500 internal_error`.
   - Response `Content-Type` header is `application/json`.
   - `GET /api/v1/healthz` returns `200` with `status: "ok"` and a
     non-empty `version`.

3. **`internal/api/middleware_test.go`**
   - Whatever middleware exists (e.g. logging, panic recovery): verify a
     panicking handler is recovered and converted to a `500` instead of
     crashing the process; verify any request logging middleware calls the
     next handler and passes through the response unchanged.

4. **`internal/config/config_test.go`**
   - Loading config with the environment variable(s) unset falls back to
     the documented default(s) (e.g. port `8080`).
   - Loading config with the environment variable(s) set uses the
     provided value.
   - Any invalid value (e.g. non-numeric port) is handled per the existing
     behavior (error or fallback — verify whichever the code does).

5. Router wiring: a lightweight test (can live in `router_test.go` or
   inside `handler_test.go`) confirming `POST /api/v1/calculate` and
   `GET /api/v1/healthz` are registered and an unknown path returns `404`.

Target coverage: 100% for `internal/calculator`, `internal/config` and `internal/api` 
(only `main.go`'s process bootstrap is excluded).
Run with `go test ./... -cover`.

### Frontend (Vitest + React Testing Library)

1. **`src/api/client.test.ts`**
   - `postJson` resolves with the parsed body on a `2xx` response.
   - A non-`2xx` response with a valid `{ error: { code, message } }` body
     throws `CalculatorApiError` carrying that `code`/`message`.
   - A non-`2xx` response with an unparsable body throws
     `CalculatorApiError` with a fallback code.
   - A `2xx` response with an unparsable/empty body throws
     `CalculatorApiError` with code `invalid_response`.
   - A network failure (fetch rejects) throws `CalculatorApiError` with
     code `network_error`.
   - The request is sent to `${BASE_URL}${path}` with method `POST`,
     `Content-Type: application/json`, and the body JSON-serialized.
   Mock `global.fetch` for all cases.

2. **`src/api/calculator.test.ts`**
   - `calculate(request)` calls `postJson` with `'/v1/calculate'` and the
     given request, and returns its result (mock `./client`).

3. **`src/utils/format.test.ts`**
   - Each exported formatting helper: typical values, zero, negative
     numbers, very small/large numbers, and any documented rounding or
     locale behavior.

4. **`src/components/CalculatorForm.test.tsx`** (React Testing Library +
   `user-event`)
   - Renders with the default operation selected and the correct number of
     operand inputs for its arity.
   - Changing the operation changes the number/labels of operand inputs to
     match `config/operations.ts` (test at least one 1-arity and one
     2-arity operation).
   - Submitting with valid input calls the calculation function/hook with
     the correctly typed/parsed operands.
   - Submitting with missing or non-numeric input shows a client-side
     validation message and does not call the API.
   - On a successful API response, the result is displayed via
     `ResultDisplay`.
   - On an API error (mock `useCalculator`/`api/calculator` to reject with
     a `CalculatorApiError`), the error message is displayed instead of a
     result.
   - A loading state disables the submit button / shows a spinner while
     the request is in flight.

5. **`src/App.test.tsx`**
   - Smoke test: the app renders without crashing and shows the
     calculator form.

Mock network calls at the `api/client` or `global.fetch` boundary — do not
let component tests make real HTTP requests. Target 100% statement
coverage across `api`, `hooks`, `utils`, and `components`. Run with
`npm run coverage`.

## Constraints

- Do not modify any file outside test files (`*_test.go`, `*.test.ts`,
  `*.test.tsx`) and test-only fixtures/mocks.
- Follow the existing project conventions: Go table-driven tests with
  `t.Run` subtests and `errors.Is` for error assertions; Vitest `describe`/`it`
  blocks with React Testing Library queries by role/label, not by test ID,
  where practical.
- Every test must assert a specific, documented behavior (status code,
  error code, computed value, or rendered text) — no snapshot-only tests.
- After writing the tests, run `go test ./... -cover` and
  `npm run coverage` and report the resulting coverage numbers.
