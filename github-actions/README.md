# GitHub Actions Runner Controller

This directory contains the configuration values for this repository's
autoscaling GitHub Actions runner. ARC creates ephemeral runners in Kubernetes
when jobs are pending and removes them when they finish.

## Requirements
- A fine-grained personal access token (PAT) limited to the
  `kunkkaharden/KubernetesSubmissions` repository, with the **Administration:
  Read and write** permission. A GitHub App is a recommended alternative for
  production installations; see the
  [ARC authentication guide](https://docs.github.com/en/actions/how-tos/manage-runners/use-actions-runner-controller/authenticate-to-the-api).

## Installation

Run the following commands from the repository root.

Install the ARC controller in its own namespace:

```bash
helm upgrade --install arc \
  --namespace arc-systems \
  --create-namespace \
  --wait \
  oci://ghcr.io/actions/actions-runner-controller-charts/gha-runner-scale-set-controller
```

First create the `arc-runners` namespace, then create the authentication
secret. The secret must be in the same namespace where the runner will be
installed. The following command reads the PAT without displaying it or
passing it as a process argument:

```bash
kubectl create namespace arc-runners

read -rsp "GitHub PAT: " GITHUB_TOKEN
echo
printf '%s' "$GITHUB_TOKEN" |
  kubectl create secret generic github-actions-auth \
    --namespace arc-runners \
    --from-file=github_token=/dev/stdin \
    --dry-run=client \
    --output=yaml |
  kubectl apply -f -
unset GITHUB_TOKEN
```

Install the scale set using the values in this directory:

```bash
helm upgrade --install kubernetes-submissions \
  --namespace arc-runners \
  --create-namespace \
  --values github-actions/values.yaml \
  oci://ghcr.io/actions/actions-runner-controller-charts/gha-runner-scale-set
```

The included values keep zero idle runners and allow scaling up to two
concurrent jobs. Adjust `maxRunners` in `values.yaml` to match the cluster's
capacity.

## Usage and verification

In a workflow in this repository, select the runner by its scale set name:

```yaml
jobs:
  example:
    runs-on: kubernetes-submissions
    steps:
      - uses: actions/checkout@v4
      - run: echo "Running on Kubernetes"
```

Check the installed resources and watch the runner pods:

```bash
helm list --all-namespaces
kubectl get pods --namespace arc-systems
kubectl get pods --namespace arc-runners --watch
```

Runners are created only while there are jobs to run, so there may be no
runner pods when the scale set is idle.

## Uninstallation

```bash
helm uninstall kubernetes-submissions --namespace arc-runners
helm uninstall arc --namespace arc-systems
kubectl delete namespace arc-runners arc-systems
```

Deleting the namespaces also deletes the authentication secret.
