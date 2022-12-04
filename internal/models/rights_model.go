package models

import "github.com/google/uuid"

// id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
//     name VARCHAR(256),
//     rules text[]

// TODO: create type
const (
	Admin    = "admin"
	Customer = "customer"
	Seller   = "seller"
	Vendor   = "vendor"
)

type Rights struct {
	UUID  uuid.UUID
	Name  string
	Rules []string
}
