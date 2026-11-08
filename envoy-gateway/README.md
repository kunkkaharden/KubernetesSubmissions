# Envoy Gateway

Install Envoy Gateway and wait until it is ready:

```bash
kubectl apply --server-side -f \
  https://github.com/envoyproxy/gateway/releases/latest/download/install.yaml
```

Apply the `GatewayClass` manifest in this directory:

```bash
kubectl apply -f gatewayclass.yaml
```




Check the resources:

```bash
kubectl get gatewayclass
kubectl get gateway,httproute -n default
```

Test locally:

```bash
export ENVOY_SERVICE=$(kubectl get svc -n envoy-gateway-system \
  -l gateway.envoyproxy.io/owning-gateway-namespace=exercises \
  -o jsonpath='{.items[0].metadata.name}')

# 2. Check the service name
echo $ENVOY_SERVICE
# 1. Get the proxy service name (by finding the Gateway in 'exercises')
export ENVOY_SERVICE=$(kubectl get svc -n envoy-gateway-system \
  -l gateway.envoyproxy.io/owning-gateway-namespace=exercises \
  -o jsonpath='{.items[0].metadata.name}')

# 2. Check the service name
echo $ENVOY_SERVICE

# 3. Set up port forwarding
kubectl port-forward -n envoy-gateway-system service/$ENVOY_SERVICE 8888:80
```

In another terminal:

```bash
curl -H "Host: www.example.com" http://localhost:8888/get
```
