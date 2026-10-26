# Log output

Apply the Deployment to your Kubernetes cluster:

```bash
kubectl apply -f manifests/ -R
```

To follow the Pod logs in real time:

```bash
kubectl logs -l app=log-output --follow
```

