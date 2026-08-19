import type { CalculateRequest, CalculateResponse } from '../types';
import { postJson } from './client';

/** Calls the backend to perform a single calculation. */
export function calculate(request: CalculateRequest): Promise<CalculateResponse> {
  return postJson<CalculateRequest, CalculateResponse>('/v1/calculate', request);
}
