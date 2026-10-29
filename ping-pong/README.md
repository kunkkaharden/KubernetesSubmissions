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

The `/pingpong` path is routed to this application by the shared Ingress in
`../log-output/manifests/ingress.yaml`. In Kubernetes, the request count is
stored in `/data/pingpong.txt` on the shared PersistentVolumeClaim.
