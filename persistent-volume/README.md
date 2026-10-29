# Persistent volume

The local PersistentVolume uses `/mnt/data` on the Kubernetes node named in
`pv.yaml`. Make sure `nodeAffinity` in that file matches a node in your cluster
(`kubectl get nodes -o wide`).

For the included k3d node, create the directory and make it writable by the
non-root container user (UID 65532):

```bash
docker exec k3d-k3s-default-agent-0 mkdir -p /mnt/data
docker exec k3d-k3s-default-agent-0 chown 65532:65532 /mnt/data
```

Apply the PersistentVolume and PersistentVolumeClaim from this directory:

```bash
kubectl apply -f .
```
