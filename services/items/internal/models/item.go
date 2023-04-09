package models

import "github.com/google/uuid"

type Item struct {
	Id          uuid.UUID
	Title       string
	Description string
	Price       int32
	Category    Category
	Vendor      string
	Images      []string
}

type ItemWithQuantity struct {
	Item
	Quantity int
}
