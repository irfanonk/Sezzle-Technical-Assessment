import { ApiError, NETWORK_ERROR_CODE, UNEXPECTED_RESPONSE_CODE } from './errors';
import type { ApiResponse } from './types';

const API_BASE_URL = import.meta.env.VITE_API_BASE_URL ?? '/api';

/**
 * Single entry point for backend calls: it applies the JSON headers, unwraps
 * the `data`/`error` envelope, and turns every failure into an ApiError.
 */
export async function apiRequest<T>(path: string, init: RequestInit = {}): Promise<T> {
  let response: Response;

  try {
    response = await fetch(`${API_BASE_URL}${path}`, {
      ...init,
      headers: {
        Accept: 'application/json',
        ...(init.body === undefined ? {} : { 'Content-Type': 'application/json' }),
        ...init.headers,
      },
    });
  } catch {
    throw new ApiError({
      code: NETWORK_ERROR_CODE,
      message: 'Cannot reach the calculator service. Check your connection and try again.',
    });
  }

  const envelope = await readEnvelope<T>(response);

  if (envelope?.error) {
    throw new ApiError({
      code: envelope.error.code,
      message: envelope.error.message,
      status: response.status,
    });
  }

  if (!response.ok) {
    throw new ApiError({
      code: UNEXPECTED_RESPONSE_CODE,
      message: `The calculator service responded with status ${response.status}.`,
      status: response.status,
    });
  }

  if (envelope?.data == null) {
    throw new ApiError({
      code: UNEXPECTED_RESPONSE_CODE,
      message: 'The calculator service returned an unexpected response.',
      status: response.status,
    });
  }

  return envelope.data;
}

async function readEnvelope<T>(response: Response): Promise<ApiResponse<T> | null> {
  try {
    return (await response.json()) as ApiResponse<T>;
  } catch {
    return null;
  }
}
