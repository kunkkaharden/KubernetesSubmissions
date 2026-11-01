# Todo app

This HTTP server listens on the port in `PORT`, or `8080`. It displays a
random image from Lorem Picsum, cached for 10 minutes in `IMAGE_CACHE_DIR`.
The image is stored as a plain file and its modification time marks when it
was fetched, so the cache also survives a restart. Once the 10 minutes are
up, the next request downloads and persists a new image. If Picsum is
unavailable, the last saved image is served instead.

The page lists the todos stored by the todo backend and has a todo input
limited to 140 characters. Submitting the form creates the todo through the
backend and reloads the list.

Requests to `/todos` are proxied by this server to the URL in
`TODO_BACKEND_URL`, or `http://localhost:8080`. That way the browser only talks
to the todo app, so the backend does not need its own Ingress.

Run the backend and then the app locally:

```bash
cd ../todo-backend && PORT=8080 go run .
```

```bash
PORT=3000 IMAGE_CACHE_DIR=./data go run .
```

The cache is stored in `todo-app/manifests/pv.yaml` and
`todo-app/manifests/pvc.yaml`. Before applying the manifests, make sure the
directory `/mnt/data/todo-app-image-cache` exists on node
`k3d-k3s-default-agent-0`. Update the PV's node affinity and path if your
cluster uses a different node or location.

Then deploy the backend and the app:

```bash
kubectl apply -f ../namespaces/project.yaml
kubectl apply -f ../todo-backend/manifests/
kubectl apply -f manifests/
```

The app reaches the backend through the `todo-backend-svc` Service, configured
with `TODO_BACKEND_URL` in the Deployment.

To follow the Pod logs in real time:

```bash
kubectl logs -n project -l app=todo-app --follow
```

To verify that the image survives a Pod replacement, note the current image,
delete the Pod, wait for its replacement, and reload the app:

```bash
kubectl delete pod -n project -l app=todo-app
```
