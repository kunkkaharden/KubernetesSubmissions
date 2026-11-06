## Using Floci as a Local Cloud Emulator

Floci emulates cloud services locally. To route traffic to applications running in Kubernetes, install the NGINX Ingress Controller.


```bash
kubectl apply -f \
https://raw.githubusercontent.com/kubernetes/ingress-nginx/controller-v1.12.1/deploy/static/provider/cloud/deploy.yaml
```