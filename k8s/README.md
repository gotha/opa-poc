# Setup in k8s

```sh
kind create cluster --name opa

```

## Build app image and load it into kind clusterA

```sh
docker build -t myapp:latest ..
kind load docker-image myapp:latest --name opa
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

## Test it

```sh
kubectl port-forward svc/myapp 8080:8080
```

```sh
curl -H "Authorization: Bearer some-valid-token-1" localhost:8080/resource1
curl -H "Authorization: Bearer some-invalid-token" localhost:8080/resource1
```
