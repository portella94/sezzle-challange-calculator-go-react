import type { ApiErrorBody } from '../types';

/** Base URL of the API. In dev/prod, "/api" is proxied to the Go backend. */
const BASE_URL = (import.meta.env.VITE_API_BASE_URL ?? '/api').replace(/\/+$/, '');

/**
 * Error thrown for any failed API call. `code` is the backend's stable,
 * machine-readable error code (e.g. "division_by_zero"), or a client-side code
 * such as "network_error".
 */
export class CalculatorApiError extends Error {
  readonly code: string;

  constructor(code: string, message: string) {
    super(message);
    this.name = 'CalculatorApiError';
    this.code = code;
  }
}

/**
 * POSTs a JSON body to the given API path and returns the parsed JSON response.
 * Non-2xx responses and network failures are surfaced as CalculatorApiError.
 */
export async function postJson<TReq, TRes>(path: string, body: TReq): Promise<TRes> {
  let res: Response;
  try {
    res = await fetch(`${BASE_URL}${path}`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(body),
    });
  } catch {
    throw new CalculatorApiError('network_error', 'Unable to reach the calculator service.');
  }

  const data = (await res.json().catch(() => null)) as (TRes & Partial<ApiErrorBody>) | null;

  if (!res.ok) {
    const err = data?.error;
    throw new CalculatorApiError(err?.code ?? 'unknown_error', err?.message ?? 'The request failed.');
  }
  if (data === null) {
    throw new CalculatorApiError('invalid_response', 'The server returned an unreadable response.');
  }
  return data as TRes;
}
