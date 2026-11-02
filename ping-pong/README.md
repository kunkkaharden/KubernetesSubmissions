# Ping-pong

The server listens on the port in `PORT`, or `8080`. The counter is stored in
Postgres, so it survives application restarts. The connection is configured
with `DATABASE_URL`, or with the standard `POSTGRES_USER`, `POSTGRES_PASSWORD`,
`POSTGRES_DB`, `DB_HOST` and `DB_PORT` variables (defaults:
`postgres://postgres:postgres@localhost:5432/pingpong?sslmode=disable`).

Run locally with a Postgres database available at that address:

```bash
PORT=3000 go run .
```

Apply the Kubernetes resources:

```bash
kubectl apply -f manifests/ -R
```

The Postgres StatefulSet has one replica and a persistent volume claim. Both
the database and the application read their configuration from the
`postgres-secret` Secret through `envFrom`, so adding a variable to the Secret
is enough to make it available in both pods. The application waits for the
database to become available before starting its HTTP server.

Each GET to `/pingpong` (or `/`) increments the counter and returns
`pong <count>`. `GET /pings` returns the current count without incrementing it;
Log output uses this endpoint over the `ping-pong-svc` Kubernetes Service.

The `/pingpong` path is routed to this application by the shared Ingress in
`../log-output/manifests/ingress.yaml`.
