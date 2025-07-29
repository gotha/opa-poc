# Setup in k8s

```sh
kind create cluster --name opa --config cluster-config.yaml
```

## Build app image and load it into kind clusterA

```sh
docker build -t myapp:latest ..
kind load docker-image myapp:latest --name opa
```

## Install Redis

```sh
helm repo add bitnami https://charts.bitnami.com/bitnami
helm repo update


helm install redis bitnami/redis -f redis.values.yaml

wait for redis to start
kubectl wait --namespace default \
  --for=condition=ready pod \
  --selector=apps.kubernetes.io/pod-index=0 \
  --timeout=90s
```

## Install Webdis

```sh
helm install webdis ./webdis
```

## Set allowlist via webdis

```sh
kubectl port-forward svc/webdis 7379:7379
```

```sh
curl -X PUT \
  -H "Content-Type: application/json" \
  --data '{"some-valid-token-1": ["/resource1", "/resource2"], "some-valid-token-2": ["/resource2"]}' \
  http://localhost:7379/SET/myapp_allowlist
```

## Install OPA

```shA
helm install opa ./opa
```

keep in mind that the rego policy is in `./opa/templates/configmap.yaml`


## Install MyApp

```sh
helm install myapp ./myapp
```

## Install ingress

based on https://kind.sigs.k8s.io/examples/ingress/deploy-ingress-nginx.yaml

```sh
kubectl apply -f deploy-ingress-nginx.yaml
```


## Test it

```/etc/hosts
127.0.0.1 myapp.local
```

```sh
curl -H "Authorization: Bearer some-valid-token-1" myapp.local:1080/resource1
curl -H "Authorization: Bearer some-invalid-token" myapp.lcaol:1080/resource1
```

