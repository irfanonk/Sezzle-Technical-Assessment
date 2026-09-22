export const REQUIRED_MESSAGE = 'Enter a value.';
export const INVALID_NUMBER_MESSAGE = 'Enter a valid number.';

export type OperandValidation = {
  errors: (string | null)[];
  operands: number[];
  isValid: boolean;
};

export function validateOperand(rawValue: string): string | null {
  const value = rawValue.trim();
  if (value === '') {
    return REQUIRED_MESSAGE;
  }

  const parsed = Number(value);
  if (!Number.isFinite(parsed)) {
    return INVALID_NUMBER_MESSAGE;
  }

  return null;
}

/**
 * Frontend validation only covers the request contract (present, finite
 * numbers). Domain rules such as division by zero stay with the backend.
 */
export function validateOperands(rawValues: string[]): OperandValidation {
  const errors = rawValues.map(validateOperand);

  return {
    errors,
    operands: rawValues.map((value) => Number(value.trim())),
    isValid: errors.every((error) => error === null),
  };
}
