import { describe, expect, it } from 'vitest';
import { INVALID_NUMBER_MESSAGE, REQUIRED_MESSAGE, validateOperands } from './operands';

describe('validateOperands', () => {
  it('accepts finite numbers, including negatives and decimals', () => {
    const validation = validateOperands(['-4', ' 2.5 ']);

    expect(validation.isValid).toBe(true);
    expect(validation.operands).toEqual([-4, 2.5]);
  });

  it('requires every operand to be filled in', () => {
    const validation = validateOperands(['', '  ']);

    expect(validation.isValid).toBe(false);
    expect(validation.errors).toEqual([REQUIRED_MESSAGE, REQUIRED_MESSAGE]);
  });

  it('rejects values that are not finite numbers', () => {
    const validation = validateOperands(['abc', 'Infinity']);

    expect(validation.isValid).toBe(false);
    expect(validation.errors).toEqual([INVALID_NUMBER_MESSAGE, INVALID_NUMBER_MESSAGE]);
  });
});
