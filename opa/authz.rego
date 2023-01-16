package httpapi.authz

default allow = false

# Allow customer actions
allow {
    input.endpoint == "updateUserPassword"
    input.role == "customer"
}
allow {
    input.endpoint == "index"
    input.role == "customer"
}

#allow {
#    input.role == "deliveryman"
#}

#allow {
#    input.role == "vendor"
#}

# Allow seller actions
allow {
    input.endpoint == "index"
    input.role == "seller"
}
allow {
    input.endpoint == "updateUserPassword"
    input.role == "seller"
}

allow {
    input.endpoint == "createCategory"
    input.role == "seller"
}

allow {
    input.endpoint == "upCatImg"
    input.role == "seller"
}

allow {
    input.endpoint == "createItem"
    input.role == "seller"
}

allow {
    input.endpoint == "upItImg"
    input.role == "seller"
}

allow {
    input.endpoint == "updateCategory"
    input.role == "seller"
}

allow {
    input.endpoint == "updateItem"
    input.role == "seller"
}

allow {
    input.endpoint == "catImgDel"
    input.role == "seller"
}

allow {
    input.endpoint == "categoryDelete"
    input.role == "seller"
}

allow {
    input.endpoint == "itImgDel"
    input.role == "seller"
}

allow {
    input.endpoint == "itemDelete"
    input.role == "seller"
}

# Allow admin all actions
allow {
    input.role == "admin"
}