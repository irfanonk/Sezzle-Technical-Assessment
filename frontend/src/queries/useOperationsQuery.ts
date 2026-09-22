import { useQuery } from '@tanstack/react-query';
import { fetchOperations } from '../api/calculator';
import type { Operation } from '../api/types';

export const operationsQueryKey = ['operations'] as const;

/** Supported operations come from the backend and are never hardcoded here. */
export function useOperationsQuery() {
  return useQuery<Operation[]>({
    queryKey: operationsQueryKey,
    queryFn: ({ signal }) => fetchOperations(signal),
  });
}
