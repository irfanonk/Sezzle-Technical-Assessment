# Calculator Frontend

Responsive React client for the Go calculator API. Supported operations and
their operand counts are loaded from the backend rather than duplicated in the
frontend.

## Stack

- React and TypeScript
- Vite
- Tailwind CSS
- TanStack Query
- Vitest and React Testing Library

## Requirements

- Node.js 22 LTS
- pnpm

## Setup

```sh
pnpm install
cp .env.example .env
```

## Development

Start the backend first (it listens on port `8080`):

```sh
cd ../backend
go run ./cmd/server
```

In another terminal, start the frontend:

```sh
pnpm dev
```

Vite prints the local frontend URL, normally `http://localhost:5173`.

During development, requests to `/api` are proxied to the `BACKEND_URL`
configured in `.env`. This keeps the browser on one origin and avoids requiring
a development CORS policy in the backend.

## API configuration

Copy `.env.example` to `.env` and adjust these values when needed:

```dotenv
BACKEND_URL=http://localhost:8080
VITE_API_BASE_URL=/api
```

- `BACKEND_URL` is read by Vite and controls the local development proxy
  target. It is not exposed to browser code.
- `VITE_API_BASE_URL` is used by the browser as the API base path. Keep it as
  `/api` when using the development proxy.

For a separately hosted backend, set `VITE_API_BASE_URL` to its public API URL
before building:

```sh
VITE_API_BASE_URL=https://calculator.example.com/api pnpm build
```

The frontend uses:

- `GET /api/operations` to load operation metadata
- `POST /api/calculate` to submit calculations

## Test 

```sh
pnpm test
pnpm run test:coverage  
```

## Type check, and build

```sh
pnpm typecheck
pnpm build
```

To preview the production build locally:

```sh
pnpm preview
```

## Docker

Build the production image:

```sh
docker build -t my-frontend .
```

### Run the frontend container

With the backend running locally on port `8080`, run the frontend by itself:

```sh
docker run --rm \
  -p 3000:80 \
  -e BACKEND_URL=http://host.docker.internal:8080 \
  my-frontend
```

Open `http://localhost:3000`. Nginx serves the built application and proxies
`/api` to `BACKEND_URL`, keeping browser requests on one origin.

When frontend and backend containers share a Docker network and the backend
service is named `backend`, the image's default
`BACKEND_URL=http://backend:8080` works without an override.

On Linux, reaching a backend running directly on the host may also require:

```sh
--add-host=host.docker.internal:host-gateway
```

### Run the full stack with Docker Compose

From the repository root:

```sh
docker compose up --build
```

Open `http://localhost:3000`. The backend is also available directly at
`http://localhost:8080`.

## Structure

- `src/api/` — `fetch` wrapper, response envelope unwrapping, and error
  normalization
- `src/queries/` — TanStack Query hooks (`useOperationsQuery`,
  `useCalculateMutation`) and shared query-client configuration
- `src/components/` — rendering, local form state, and user interaction
- `src/validation/` — request contract checks performed before submission
- `src/test/` — shared test setup and helpers

## Design notes

Operations are always loaded from `GET /api/operations`; the UI renders one
operand input per advertised `arity`, so adding a backend operation requires no
frontend change.

Frontend validation only covers what can be checked locally (required inputs
and finite numbers). Domain rules such as division by zero stay authoritative
in the backend, and their error messages are displayed as returned.

Components never call `fetch` directly. The API layer normalizes HTTP, backend,
and network failures into `ApiError` values, while dedicated query hooks expose
server state to the UI. Operations remain fresh for five minutes.

The interface is mobile-first and expands operand fields horizontally on wider
screens.

## Troubleshooting

If switching from npm to pnpm, remove the old installation before reinstalling:

```sh
rm -rf node_modules
rm -f package-lock.json
pnpm install
```

Commit `pnpm-lock.yaml` so all environments resolve the same dependency
versions.
