# Todo app

This HTTP server listens on the port in `PORT`, or `8080`.

Run locally:

```bash
PORT=3000 go run .
```

```bash
kubectl apply -f manifests/deployment.yaml
```

To follow the Pod logs in real time:

```bash
kubectl logs -l app=todo-app --follow
```

