# Todo backend

A small HTTP API for todos. The server listens on the port in `PORT`, or
`8080`. Todos are kept in memory, so they are lost when the process restarts.

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

```bash
PORT=3000 go run .
```

```bash
curl localhost:3000/todos
curl -X POST localhost:3000/todos -H 'Content-Type: application/json' \
  -d '{"text":"Learn Kubernetes"}'
```

## Run in Kubernetes

Apply the resources:

```bash
kubectl apply -f manifests/
```

The `todo-backend-svc` Service exposes the backend to other Pods in the
cluster on port `2345`. From inside the cluster it is reachable at
`http://todo-backend-svc:2345/todos`; an Ingress is not needed for Pod-to-Pod
traffic. To try the API from your machine:

```bash
kubectl port-forward svc/todo-backend-svc 2345:2345
curl localhost:2345/todos
```
