# OPA PoC for API key management and more ...

Run OPA server

```sh
opa run apikey-policy.rego --server --addr :8181
```

test if it is working fine

```sh
curl -X POST http://localhost:8181/v1/data/myapp/allow \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "request": {
        "headers": {
          "authorization": "Bearer some-valid-token-1"
        }
      }
    }
  }'
```

start go web server

```sh
go run main
```

and test auth 

```sh
curl -H "Authorization: Bearer some-valid-token-1" localhost:8080/resource1
curl -H "Authorization: Bearer some-invalid-token" localhost:8080/resource1
```

# TODO 

[ ] allow a key to be used only on a specific list of resources
[ ] load valid api keys from external store - vault or redis
[ ] setup OPA policy replication and configure go server to use replica
[ ] integrate OPA server with ingress directly and remove auth code from go server
