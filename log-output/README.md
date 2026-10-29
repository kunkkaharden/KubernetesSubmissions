# Log output

The `log-writer` container generates a random string when it starts and appends
it with a UTC timestamp to `/data/log.txt` every five seconds. The `log-reader`
container reads that file and the Ping-pong request count from
`/data/pingpong.txt`, returning both from `GET /`. Both containers and the
Ping-pong Pod share a PersistentVolumeClaim.

Prepare and apply the shared PersistentVolume and claim first. Follow the
instructions in [`persistent-volume/README.md`](../persistent-volume/README.md).

```bash
kubectl apply -f ../persistent-volume/
```

Apply Ping-pong before Log output so the latter can be scheduled on the same
node and share the local hostPath volume:

```bash
kubectl apply -f ../ping-pong/manifests/
kubectl apply -f manifests/
```

To see the latest file contents:

```bash
kubectl port-forward service/log-output-svc 3000:3000
# In another terminal:
curl http://localhost:3000/
```
