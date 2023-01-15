package httpapi.authz

# io.jwt.decode_verify
# https://www.openpolicyagent.org/docs/latest/language-reference/#tokens
token = t {
  [valid, _, payload] = io.jwt.decode_verify(input.token, { "secret": "secret" })
  t := {
    "valid": valid,
    "payload": payload
  }
}

default allow = false

# Allow admin all actions
allow {
  token.valid

  input.method == [_]
  input.path = [_]
  token.payload.role == admin
}