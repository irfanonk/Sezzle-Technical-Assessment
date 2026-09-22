import { QueryClientProvider } from '@tanstack/react-query';
import { render } from '@testing-library/react';
import { vi } from 'vitest';
import type { ReactElement } from 'react';
import { createQueryClient } from '../queries/queryClient';

/** The real client, minus retries so failure paths resolve immediately. */
function createTestQueryClient() {
  const client = createQueryClient();
  client.setDefaultOptions({
    queries: { retry: false },
    mutations: { retry: false },
  });
  return client;
}

export function renderWithClient(ui: ReactElement) {
  const client = createTestQueryClient();
  return render(<QueryClientProvider client={client}>{ui}</QueryClientProvider>);
}

type FetchHandler = (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>;

export function stubFetch(handler: FetchHandler) {
  const mock = vi.fn(handler);
  vi.stubGlobal('fetch', mock);
  return mock;
}

export function jsonResponse(body: unknown, status = 200): Response {
  return {
    ok: status >= 200 && status < 300,
    status,
    json: async () => body,
  } as unknown as Response;
}

export function deferred<T>() {
  let resolve!: (value: T) => void;
  const promise = new Promise<T>((resolvePromise) => {
    resolve = resolvePromise;
  });
  return { promise, resolve };
}

export function calculateCalls(mock: { mock: { calls: unknown[][] } }) {
  return mock.mock.calls.filter((call) => String(call[0]).endsWith('/calculate'));
}
