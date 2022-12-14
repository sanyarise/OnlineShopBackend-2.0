package models

import "github.com/google/uuid"

// id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
//     name VARCHAR(256),
//     rules text[]

// TODO: create type
const (
	Admin    = "Admin"
	Customer = "Customer"
	Seller   = "Seller"
)

type Rights struct {
	ID    uuid.UUID
	Name  string
	Rules []string
}
