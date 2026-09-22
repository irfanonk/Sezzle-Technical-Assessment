# Calculator Frontend

React + TypeScript client for the Go calculator API, built with Vite, Tailwind
CSS, and TanStack Query.

## Requirements

- Node.js 20 or newer

## Setup

```sh
npm install
```

## Run

Start the backend first (it listens on port `8080`):

```sh
cd ../backend && go run ./cmd/server
```

Then start the dev server:

```sh
npm run dev
```

The dev server proxies `/api` to `http://localhost:8080`, so the browser stays
on a single origin and the backend needs no CORS policy. To point at another
backend without the proxy, set `VITE_API_BASE_URL` to an absolute URL.

## Test, type check, and build

```sh
npm test
npm run typecheck
npm run build
```

## Structure

- `src/api/` — `fetch` wrapper, response envelope unwrapping, and error
  normalization. Every failure leaves this layer as an `ApiError`.
- `src/queries/` — TanStack Query hooks (`useOperationsQuery`,
  `useCalculateMutation`) and the shared client with a 5 minute cache.
- `src/components/` — rendering, local form state, and user interaction.
- `src/validation/` — request contract checks performed before a request.

## Design notes

Operations are always loaded from `GET /api/operations`; the UI renders one
operand input per advertised `arity`, so adding a backend operation requires no
frontend change.

Frontend validation only covers what can be checked locally (required inputs
and finite numbers). Domain rules such as division by zero stay authoritative
in the backend, and their error messages are displayed as returned.
