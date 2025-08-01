package myapp

default allow = false

# Load allowlist from Redis
resp := http.send({
  "method": "GET",
  "url": "http://localhost:7379/GET/myapp_allowlist",
  "headers": {"Content-Type": "application/json"},
  "force_json_decode": true
})

allowlist := json.unmarshal(resp.body.GET) if {
   resp.status_code == 200
   resp.body != ""
}

key_allowed(key) if {
  keys := {k | k := key; allowlist[key]}
  key in keys
}

path_allowed(token, request_path) if {
    allowed_paths := allowlist[token]
    some allowed_path in allowed_paths
    startswith(request_path, allowed_path)
}

# Main allow rule
allow if {
  authz_header := input.headers["authorization"]
  startswith(authz_header, "Bearer ")
  api_key := substring(authz_header, 7, -1)
  key_allowed(api_key)
  path_allowed(api_key, input.request.path)
}

allow if {
  input.path == "/"
}

# Debug helpers - remove these in production
# ==========================================
key_valid if {
  startswith(input.request.headers.authorization, "Bearer ")
  api_key := substring(input.request.headers.authorization, 7, -1)
  key_allowed(api_key)
}

path_valid if {
  startswith(input.request.headers.authorization, "Bearer ")
  api_key := substring(input.request.headers.authorization, 7, -1)
  path_allowed(api_key, input.request.path)
}

debug_info := {
  "request_path": input.request.path,
}
