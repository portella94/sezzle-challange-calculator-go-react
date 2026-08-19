import { useCallback, useMemo, useState } from 'react';
import type { Operation } from '../types';
import { OPERATION_MAP, type OperationMeta } from '../config/operations';
import { calculate } from '../api/calculator';
import { CalculatorApiError } from '../api/client';
import { parseOperand } from '../utils/format';

/** Per-operand validation messages, indexed by operand position. */
export type FieldErrors = Record<number, string>;

export interface UseCalculator {
  operation: Operation;
  meta: OperationMeta;
  /** Raw operand input strings, one per operand position (max arity = 2). */
  operands: string[];
  result: number | null;
  error: string | null;
  fieldErrors: FieldErrors;
  loading: boolean;
  setOperation: (op: Operation) => void;
  setOperand: (index: number, value: string) => void;
  submit: () => Promise<void>;
}

const MAX_OPERANDS = 2;

/**
 * Owns all calculator interaction state and the submit workflow. Validation is
 * fail-fast and client-side first (required + numeric), but the backend remains
 * the source of truth for domain errors (division by zero, etc.).
 */
export function useCalculator(initial: Operation = 'add'): UseCalculator {
  const [operation, setOperationState] = useState<Operation>(initial);
  const [operands, setOperands] = useState<string[]>(Array(MAX_OPERANDS).fill(''));
  const [result, setResult] = useState<number | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [fieldErrors, setFieldErrors] = useState<FieldErrors>({});
  const [loading, setLoading] = useState(false);

  const meta = OPERATION_MAP[operation];

  const setOperation = useCallback((op: Operation) => {
    setOperationState(op);
    // Changing the operation invalidates prior output and validation state.
    setResult(null);
    setError(null);
    setFieldErrors({});
  }, []);

  const setOperand = useCallback((index: number, value: string) => {
    setOperands((prev) => {
      const next = [...prev];
      next[index] = value;
      return next;
    });
    // A fresh edit clears that field's error and any stale result.
    setFieldErrors((prev) => {
      if (prev[index] === undefined) return prev;
      const { [index]: _removed, ...rest } = prev;
      return rest;
    });
    setResult(null);
  }, []);

  const submit = useCallback(async () => {
    setError(null);
    setResult(null);

    // Client-side validation for the operands this operation actually uses.
    const errors: FieldErrors = {};
    const parsed: number[] = [];
    for (let i = 0; i < meta.arity; i++) {
      const value = parseOperand(operands[i] ?? '');
      if (value === null) {
        errors[i] = operands[i]?.trim() ? 'Enter a valid number' : 'This field is required';
      } else {
        parsed.push(value);
      }
    }
    if (Object.keys(errors).length > 0) {
      setFieldErrors(errors);
      return; // fail fast — do not hit the network with invalid input
    }
    setFieldErrors({});

    setLoading(true);
    try {
      const response = await calculate({ operation, operands: parsed });
      setResult(response.result);
    } catch (err) {
      const message =
        err instanceof CalculatorApiError ? err.message : 'Something went wrong. Please try again.';
      setError(message);
    } finally {
      setLoading(false);
    }
  }, [meta.arity, operands, operation]);

  return useMemo(
    () => ({
      operation,
      meta,
      operands,
      result,
      error,
      fieldErrors,
      loading,
      setOperation,
      setOperand,
      submit,
    }),
    [operation, meta, operands, result, error, fieldErrors, loading, setOperation, setOperand, submit],
  );
}
