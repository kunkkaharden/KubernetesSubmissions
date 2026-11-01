# Log output

The `log-writer` container generates a random string when it starts and appends
it with a UTC timestamp to `/data/log.txt` every five seconds. The `log-reader`
container reads that file and gets the Ping-pong request count from
`http://ping-pong-svc:3000/pings`, returning both from `GET /`.
The log writer and reader share a PersistentVolumeClaim for the log file; the
Ping-pong Pod does not use this volume.

Prepare and apply the shared PersistentVolume and claim first. Follow the
instructions in [`persistent-volume/README.md`](../persistent-volume/README.md).

```bash
kubectl apply -f ../persistent-volume/
```

Apply Ping-pong and Log output:

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
