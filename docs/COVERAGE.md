# Coverage

How to generate coverage reports for both layers, and the latest snapshot.

## Backend (Go)

```bash
cd backend
go test ./... -coverprofile=coverage.out -covermode=atomic
go tool cover -func=coverage.out            # per-function + total, in terminal
go tool cover -html=coverage.out -o coverage.html   # browsable HTML report
```

Latest snapshot:

| Package                | Coverage |
| ---------------------- | -------: |
| `internal/calculator`  |   100.0% |
| `internal/config`      |   100.0% |
| `internal/api`         |    97.3% |
| `cmd/server`           |     0.0% (composition root, intentionally untested) |
| **total**              | **82.8%** |

The meaningful logic — arithmetic, validation, error mapping, middleware — is
fully covered. Only the `main` wiring function is excluded.

## Frontend (Vitest + V8)

```bash
cd frontend
npm run coverage        # text summary + HTML report in ./coverage
```

Latest snapshot (V8 provider):

| Metric      | Coverage |
| ----------- | -------: |
| Statements  |   95.7% |
| Branches    |   82.1% |
| Functions   |   95.7% |
| Lines       |   96.6% |

Covered: the API client (success / backend-error / network-error paths), number
parsing & formatting, the `useCalculator` validation and submit workflow, and
the form's render / arity-switching / success / error behaviour.

> Coverage artifacts (`coverage.out`, `coverage.html`, `frontend/coverage/`) are
> git-ignored; regenerate them with the commands above.
