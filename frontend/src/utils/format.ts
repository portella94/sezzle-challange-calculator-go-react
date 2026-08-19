/**
 * Formats a numeric result for display, stripping floating-point noise (e.g.
 * 0.1 + 0.2 = 0.30000000000000004 → "0.3") while leaving integers and short
 * decimals untouched.
 */
export function formatResult(value: number): string {
  if (!Number.isFinite(value)) return String(value);
  return String(Number(value.toPrecision(12)));
}

/**
 * Parses a user-entered operand. Returns a finite number, or null if the input
 * is empty or not a valid number (fail-fast client-side validation).
 */
export function parseOperand(input: string): number | null {
  const trimmed = input.trim();
  if (trimmed === '') return null;
  const n = Number(trimmed);
  return Number.isFinite(n) ? n : null;
}
