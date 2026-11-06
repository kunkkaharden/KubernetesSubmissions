## Add the Helm repositories

These commands add the Prometheus and Grafana repositories and update the list of available charts.

```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
```

## Install the monitoring stack

`helm upgrade --install` installs each chart or updates it if it already exists. Everything is installed in the `monitoring` namespace, which is created by the first command. The `*-values.yaml` files contain each component's configuration.

- **Prometheus**: collects and stores metrics.
- **Loki**: collects and stores logs.
- **k8s-monitoring**: configures metrics and log collection for Kubernetes.
- **Grafana**: displays metrics and logs in dashboards.

```bash
helm upgrade --install prom prometheus-community/prometheus \
  --namespace monitoring \
  --values prom-values.yaml

helm upgrade --install loki grafana/loki \
  --namespace monitoring \
  --values loki-values.yaml

helm upgrade --install k8smon grafana/k8s-monitoring \
  --namespace monitoring \
  --values k8smon-values.yaml

helm upgrade --install grafana grafana/grafana \
  --namespace monitoring \
  --values grafana-values.yaml
```
