# Todo backend

A small HTTP API for todos. Todos are stored in Postgres, so they survive
application restarts. The
server port and database connection are configured with the required `PORT`,
`POSTGRES_USER`, `POSTGRES_PASSWORD`, `POSTGRES_DB`, `DB_HOST`, `DB_PORT` and
`DB_SSLMODE` environment variables.

## Endpoints

- `GET /todos` returns all todos as a JSON array:

  ```json
  [{ "id": 1, "text": "Learn Kubernetes", "done": false }]
  ```

- `POST /todos` creates a todo from a JSON body:

  ```json
  { "text": "Learn Kubernetes" }
  ```

  It answers `201 Created` with the stored todo, including its generated `id`.
  The text is trimmed and must be between 1 and 140 characters; an empty or
  invalid body answers `400 Bad Request`.

## Run locally

Run with a Postgres database available at the configured address:

```bash
PORT=3000 DB_HOST=localhost DB_PORT=5432 DB_SSLMODE=disable \
POSTGRES_USER=postgres POSTGRES_PASSWORD=postgres POSTGRES_DB=todos \
go run .
```

```bash
curl localhost:3000/todos
curl -X POST localhost:3000/todos -H 'Content-Type: application/json' \
  -d '{"text":"Learn Kubernetes"}'
```

## Run in Kubernetes

Apply the resources:

```bash
kubectl apply -f ../namespaces/project.yaml
kubectl apply -f manifests/
```

The `postgres` StatefulSet has one replica and a persistent volume claim.
Both Postgres and the backend read the database configuration and credentials
from the `postgres-secret` Secret through `envFrom`. The backend validates this
configuration at startup
and waits for the database to become available before starting its HTTP server.

The `todo-backend-svc` Service exposes the backend to other Pods in the
cluster on port `2345`. From inside the cluster it is reachable at
`http://todo-backend-svc:2345/todos`; an Ingress is not needed for Pod-to-Pod
traffic. To try the API from your machine:

```bash
kubectl port-forward -n project svc/todo-backend-svc 2345:2345
curl localhost:2345/todos
```

## Logs in Grafana
To review the exercise, also check the `/monitoring` directory.

The monitoring configuration sends Pod logs to Loki, and Grafana already has a
Loki data source. In **Explore**, select **Loki** and use this LogQL query:

```logql
{namespace="project"}"
```
