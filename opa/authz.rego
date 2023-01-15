package httpapi.authz

default allow = false

# Allow admin all actions
allow {
    input.role == "Admin"
}