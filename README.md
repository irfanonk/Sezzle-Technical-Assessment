# Calculator Application

Full-stack calculator with a Go backend and a React frontend.

## Run both applications with Docker Compose

From the repository root:

```sh
docker compose up --build
```

Once the containers are running:

- Frontend: `http://localhost:3000`
- Backend API: `http://localhost:8080`
- Supported operations: `http://localhost:8080/api/operations`

The frontend's Nginx server proxies `/api` requests to the backend service over
the internal Compose network.

To run the containers in the background:

```sh
docker compose up --build -d
```

View logs:

```sh
docker compose logs -f
```

Stop and remove the containers:

```sh
docker compose down
```

## Project documentation

- [Backend setup and API documentation](backend/README.md)
- [Frontend setup and architecture](frontend/README.md)
