// Shared types mirroring the backend API contract.

/** Supported operations — must stay in sync with the Go backend registry. */
export type Operation =
  | 'add'
  | 'subtract'
  | 'multiply'
  | 'divide'
  | 'power'
  | 'sqrt'
  | 'percentage';

/** Request body for POST /api/v1/calculate. */
export interface CalculateRequest {
  operation: Operation;
  operands: number[];
}

/** Success response body from POST /api/v1/calculate. */
export interface CalculateResponse {
  operation: Operation;
  operands: number[];
  result: number;
}

/** Error envelope returned by the backend on failure. */
export interface ApiErrorBody {
  error: {
    code: string;
    message: string;
  };
}
