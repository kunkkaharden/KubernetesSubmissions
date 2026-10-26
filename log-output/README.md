# Log output

The `log-writer` container generates a random string when it starts and appends
it with a UTC timestamp to `/data/log.txt` every five seconds. The `log-reader`
container reads that shared file and returns its contents from `GET /`. Both
containers run in the same Pod and share an `emptyDir` volume.


Apply the Deployment to your Kubernetes cluster:

```bash
kubectl apply -f manifests/ -R
```

To see the latest file contents:

```bash
kubectl port-forward service/log-output-svc 3000:3000
# In another terminal:
curl http://localhost:3000/
```
