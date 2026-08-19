# Calculator Frontend (React + TypeScript)

Strict-TypeScript React app built with Vite and MUI. Consumes the Go backend's
`POST /api/v1/calculate` endpoint through a two-operand form.

- **Vite:** 8.2.1 · **React:** 19 · **MUI:** 9.3.1 · **Vitest:** 4.1.10
- **Node:** 24 LTS recommended
- **Version:** 0.0.1 (SemVer)

## Layout

```
src/
  api/          HTTP client + calculate() call
  components/   OperationSelect, ResultDisplay (presentational), CalculatorForm (container)
  config/       operations.ts — single source of truth for operation labels + arity
  hooks/        useCalculator — state, client-side validation, submit workflow
  utils/        number parsing/formatting
  theme.ts      MUI theme
```

`config/operations.ts` drives everything arity-related: the form renders one or
two operand inputs (and their labels) purely from the selected operation's
metadata, keeping the UI in sync with the backend contract (DRY).

## Develop

```bash
npm install          # if npm errors on peer resolution, use: npm install --legacy-peer-deps
npm run dev          # http://localhost:5173, proxies /api → http://localhost:8080
```

Run the [backend](../backend/README.md) on port 8080 alongside it. Configure a
different API base via `VITE_API_BASE_URL` (see `.env.example`).

> **Note:** Node 24 LTS ships an npm that resolves the dependency tree cleanly.
> On older npm (10.x) add `--legacy-peer-deps` to `npm install`.

## Test, coverage & build

```bash
npm test             # run unit + component tests once (Vitest)
npm run coverage     # tests with a V8 coverage report (text + HTML in ./coverage)
npm run build        # type-check (tsc -b) then production build to ./dist
npm run lint         # type-check only
```

Tests use Vitest + React Testing Library and cover the API client, number
utilities, the form's validation/arity behaviour, and success/error rendering.
