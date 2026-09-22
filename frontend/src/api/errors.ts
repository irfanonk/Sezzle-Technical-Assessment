export const NETWORK_ERROR_CODE = 'network_error';
export const UNEXPECTED_RESPONSE_CODE = 'unexpected_response';
export const UNKNOWN_ERROR_CODE = 'unknown_error';

const FALLBACK_MESSAGE = 'Something went wrong. Please try again.';

/**
 * Every failure leaving the api layer is an ApiError, so components and hooks
 * never have to inspect raw responses or transport exceptions.
 */
export class ApiError extends Error {
  readonly code: string;
  readonly status: number | null;

  constructor({
    code,
    message,
    status = null,
  }: {
    code: string;
    message: string;
    status?: number | null;
  }) {
    super(message);
    this.name = 'ApiError';
    this.code = code;
    this.status = status;
  }
}

export function isApiError(error: unknown): error is ApiError {
  return error instanceof ApiError;
}

/** Normalizes anything thrown by the api layer or React Query into an ApiError. */
export function toApiError(error: unknown): ApiError {
  if (isApiError(error)) {
    return error;
  }

  return new ApiError({
    code: UNKNOWN_ERROR_CODE,
    message: error instanceof Error && error.message ? error.message : FALLBACK_MESSAGE,
  });
}

/** Client errors are caused by the request itself, so retrying cannot help. */
export function isClientError(error: unknown): boolean {
  const status = toApiError(error).status;
  return status !== null && status >= 400 && status < 500;
}
