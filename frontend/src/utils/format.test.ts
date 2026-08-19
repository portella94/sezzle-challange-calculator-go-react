import { describe, expect, it } from 'vitest';
import { formatResult, parseOperand } from './format';

describe('formatResult', () => {
  it('leaves integers and short decimals unchanged', () => {
    expect(formatResult(5)).toBe('5');
    expect(formatResult(0.25)).toBe('0.25');
    expect(formatResult(-3)).toBe('-3');
  });

  it('strips floating-point noise', () => {
    expect(formatResult(0.1 + 0.2)).toBe('0.3');
  });

  it('passes through non-finite values', () => {
    expect(formatResult(Number.POSITIVE_INFINITY)).toBe('Infinity');
  });
});

describe('parseOperand', () => {
  it('parses valid numbers, trimming whitespace', () => {
    expect(parseOperand('  42 ')).toBe(42);
    expect(parseOperand('-3.5')).toBe(-3.5);
  });

  it('returns null for empty input', () => {
    expect(parseOperand('')).toBeNull();
    expect(parseOperand('   ')).toBeNull();
  });

  it('returns null for non-numeric input', () => {
    expect(parseOperand('abc')).toBeNull();
    expect(parseOperand('1,2')).toBeNull();
  });
});
