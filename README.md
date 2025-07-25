# OPA PoC for API key management and more ...


## Start services

### With process compose

```sh
process-compose up --no-server
```

### Manually

start redis

```sh
redis-server
```

put access policy in redis:

```sh
redis-cli SET myapp_allowlist '{"some-valid-token-1": ["/resource1", "/resource2"], "some-valid-token-2": ["/resource2"]}'
```

Start webdis:

```sh
webdis
```

Run OPA server

```sh
opa run apikey-policy.rego --server --addr :8181
```


start go web server

```sh
go run main
```

## Test service

and test auth 

```sh
curl -H "Authorization: Bearer some-valid-token-1" localhost:8080/resource1
curl -H "Authorization: Bearer some-invalid-token" localhost:8080/resource1
```

## Debug OPA policies 

test if it is working fine

```sh
curl -X POST http://localhost:8181/v1/data/myapp/allow \
  -H "Content-Type: application/json" \
  -d '{
    "input": {
      "request": {
        "headers": {
          "authorization": "Bearer some-valid-token-1"
        },
        "path": "/resource1"
      }
    }
  }'
```

or debug the policy with you can do

```sh
cat EOF > input.json
{
  "input": {
    "request": {
      "headers": {
        "authorization": "Bearer some-valid-token-1"
      },
      "path": "/resource1"
    }
  }
}
EOF
opa eval -d apikey-policy.rego -i input.json "data.myapp.path_valid"
opa eval -d apikey-policy.rego -i input.json "data.myapp.key_valid"
```


# TODO 

[x] allow a key to be used only on a specific list of resources
[x] load valid api keys from external store - vault or redis
[ ] setup OPA policy replication and configure go server to use replica
[ ] integrate OPA server with ingress directly and remove auth code from go server
