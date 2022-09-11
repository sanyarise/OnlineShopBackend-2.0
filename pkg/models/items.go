package models

import "github.com/google/uuid"

type Item struct {
	ID          uuid.UUID  `json:"id,omitempty"`
	Name        string     `json:"name,omitempty"`
	Categories  []Category `json:"categories,omitempty"`
	Description string     `json:"description,omitempty"`
	Price       int        `json:"price,omitempty"`
	Vendor      string     `json:"vendor,omitempty"`
	Pictures    []string   `json:"pictures,omitempty"`
}
