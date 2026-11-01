# Ping-pong

The server listens on the port in `PORT`, or `8080`. Each GET request to
`/pingpong` returns the current in-memory count and increments it. The count
starts at `0` and resets when the process restarts.

Run locally:

```bash
PORT=3000 go run .
```

Apply the Kubernetes resources:

```bash
kubectl apply -f manifests/ -R
```

Each GET to `/pingpong` (or `/`) increments the in-memory counter and returns
`pong <count>`. `GET /pings` returns the current count without
incrementing it; Log output uses this endpoint over the `ping-pong-svc`
Kubernetes Service.

The `/pingpong` path is routed to this application by the shared Ingress in
`../log-output/manifests/ingress.yaml`. The count resets when the process
restarts.
