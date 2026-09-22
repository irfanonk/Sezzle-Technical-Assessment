import { apiRequest } from './client';
import type { CalculationData, CalculationInput, Operation, OperationsData } from './types';

export function fetchOperations(signal?: AbortSignal): Promise<Operation[]> {
  return apiRequest<OperationsData>('/operations', { signal }).then((data) => data.operations);
}

export function postCalculation(
  input: CalculationInput,
  signal?: AbortSignal,
): Promise<number> {
  return apiRequest<CalculationData>('/calculate', {
    method: 'POST',
    body: JSON.stringify(input),
    signal,
  }).then((data) => data.result);
}
