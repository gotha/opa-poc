package myapp

default allow = false


allowlist := {
  "some-valid-token-1": [
    "/resource1",
    "/resource2"
  ],
  "some-valid-token-2": [
    "/resource2"
  ]
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
  startswith(input.request.headers.authorization, "Bearer ")
  api_key := substring(input.request.headers.authorization, 7, -1)
  key_allowed(api_key)
  path_allowed(api_key, input.request.path)
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
