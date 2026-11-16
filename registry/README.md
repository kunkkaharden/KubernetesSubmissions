# In-cluster Container Registry

This Helm chart deploys a private Docker Registry and its web UI inside our
Kubernetes cluster. We use our own registry instead of relying on an external
registry. Both Services are `ClusterIP`; no Ingress is created, so the
Registry and UI are intended for in-cluster access.

## Install

From this directory (`registry/`), install the chart in its own namespace:

```sh
helm upgrade --install registry . \
  --namespace registry \
  --create-namespace
```

Check that the pods and Services are ready:

```sh
kubectl get pods,services,pvc --namespace registry
```

To apply later chart or configuration changes, run the same `helm
upgrade --install` command again.

## In-cluster access

The Registry Service is available on port `5000` at:

```text
registry-docker-registry-ui-registry.registry.svc.cluster.local:5000
```

The UI Service is available on port `80` at:

```text
registry-docker-registry-ui-ui.registry.svc.cluster.local
```

Workloads in the same namespace can use the shorter service names. Workloads
in another namespace can use the full DNS names above.

To open the UI locally for a quick check, forward its Service port:

```sh
kubectl port-forward --namespace registry \
  service/registry-docker-registry-ui-ui 8080:80
```

Then visit <http://localhost:8080>.

## Security

Authentication is disabled by default, and the Registry allows image deletion.
Any workload that can reach the Registry Service can read, push, or delete
images. Use this configuration only in a trusted cluster; restrict access
with Kubernetes network policies if needed.
