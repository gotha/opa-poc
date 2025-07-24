package myapp

default allow = false

allow if {
  startswith(input.request.headers.authorization, "Bearer ")
  token := substring(input.request.headers.authorization, 7, -1)
  token in allowed_tokens
}

allowed_tokens := {
  "some-valid-token-1",
  "some-valid-token-2",
}
