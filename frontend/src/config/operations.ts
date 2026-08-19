import type { Operation } from '../types';

/**
 * Metadata describing how each operation is presented and validated. This is
 * the single source of truth on the frontend for an operation's arity and
 * labels — the form derives how many operand inputs to show, and what to call
 * them, entirely from here (DRY). It mirrors the backend registry.
 */
export interface OperationMeta {
  value: Operation;
  /** Label shown in the operation selector. */
  label: string;
  /** Number of operands the operation requires. */
  arity: 1 | 2;
  /** Labels for each operand input, indexed by operand position. */
  operandLabels: string[];
  /** Short explanation of the operation's semantics. */
  hint: string;
}

export const OPERATIONS: OperationMeta[] = [
  { value: 'add', label: 'Add (+)', arity: 2, operandLabels: ['A', 'B'], hint: 'A + B' },
  { value: 'subtract', label: 'Subtract (−)', arity: 2, operandLabels: ['A', 'B'], hint: 'A − B' },
  { value: 'multiply', label: 'Multiply (×)', arity: 2, operandLabels: ['A', 'B'], hint: 'A × B' },
  { value: 'divide', label: 'Divide (÷)', arity: 2, operandLabels: ['A', 'B'], hint: 'A ÷ B (B ≠ 0)' },
  { value: 'power', label: 'Power (xʸ)', arity: 2, operandLabels: ['Base', 'Exponent'], hint: 'Base raised to Exponent' },
  { value: 'sqrt', label: 'Square root (√)', arity: 1, operandLabels: ['Value'], hint: 'Square root of Value (≥ 0)' },
  { value: 'percentage', label: 'Percentage (%)', arity: 2, operandLabels: ['Percent', 'Of'], hint: 'Percent% of Of' },
];

export const OPERATION_MAP: Record<Operation, OperationMeta> = Object.fromEntries(
  OPERATIONS.map((op) => [op.value, op]),
) as Record<Operation, OperationMeta>;
