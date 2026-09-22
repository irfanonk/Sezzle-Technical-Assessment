import { MutationCache, QueryCache, QueryClient } from '@tanstack/react-query';
import { isClientError, toApiError } from '../api/errors';

export const CACHE_TIME_MS = 5 * 60 * 1000;
const MAX_RETRIES = 2;

/**
 * Global API error handling: failures are normalized and reported in one
 * place, so components only decide how to render the message.
 */
function reportError(error: unknown) {
  const apiError = toApiError(error);
  console.error(`[api] ${apiError.code}: ${apiError.message}`);
}

export function createQueryClient(): QueryClient {
  return new QueryClient({
    queryCache: new QueryCache({ onError: reportError }),
    mutationCache: new MutationCache({ onError: reportError }),
    defaultOptions: {
      queries: {
        gcTime: CACHE_TIME_MS,
        staleTime: CACHE_TIME_MS,
        refetchOnWindowFocus: false,
        retry: (failureCount, error) => !isClientError(error) && failureCount < MAX_RETRIES,
      },
      mutations: {
        retry: false,
      },
    },
  });
}
