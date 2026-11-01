# Todo app

This HTTP server listens on the port in `PORT`, or `8080`. It displays a
random image from Lorem Picsum, cached for 10 minutes in `IMAGE_CACHE_DIR`.
The image is stored as a plain file and its modification time marks when it
was fetched, so the cache also survives a restart. Once the 10 minutes are
up, the next request downloads and persists a new image. If Picsum is
unavailable, the last saved image is served instead.

The page includes a todo input limited to 140 characters, a Send button, and
example todos. The button is not wired up yet, so todos are not submitted or
saved.

Run locally:

```bash
PORT=3000 IMAGE_CACHE_DIR=./data go run .
```

The cache is stored in `todo-app/manifests/pv.yaml` and
`todo-app/manifests/pvc.yaml`. Before applying the manifests, make sure the
directory `/mnt/data/todo-app-image-cache` exists on node
`k3d-k3s-default-agent-0`. Update the PV's node affinity and path if your
cluster uses a different node or location.

Then apply the resources:

```bash
kubectl apply -f manifests/
```

To follow the Pod logs in real time:

```bash
kubectl logs -l app=todo-app --follow
```

To verify that the image survives a Pod replacement, note the current image,
delete the Pod, wait for its replacement, and reload the app:

```bash
kubectl delete pod -l app=todo-app
```
