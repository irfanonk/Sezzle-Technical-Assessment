import { useMutation } from '@tanstack/react-query';
import { postCalculation } from '../api/calculator';
import type { ApiError } from '../api/errors';
import type { CalculationInput } from '../api/types';

export function useCalculateMutation() {
  return useMutation<number, ApiError, CalculationInput>({
    mutationFn: (input) => postCalculation(input),
  });
}
