import { toApiError } from '../api/errors';

type ErrorNoticeProps = {
  error: unknown;
  onRetry?: () => void;
};

/** Renders an already-normalized API failure; no parsing happens here. */
export function ErrorNotice({ error, onRetry }: ErrorNoticeProps) {
  const apiError = toApiError(error);

  return (
    <div
      role="alert"
      className="rounded-lg border border-red-300 bg-red-50 p-3 text-sm text-red-800"
    >
      <p>{apiError.message}</p>
      {onRetry && (
        <button
          type="button"
          onClick={onRetry}
          className="mt-2 rounded-md border border-red-400 px-3 py-1.5 font-medium text-red-800 hover:bg-red-100"
        >
          Retry
        </button>
      )}
    </div>
  );
}
