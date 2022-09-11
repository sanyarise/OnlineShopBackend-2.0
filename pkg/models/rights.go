package models

import "github.com/google/uuid"

type Right struct {
	ID    uuid.UUID `json:"id,omitempty"`
	Name  string    `json:"name,omitempty"`
	Rules []string  `json:"rules,omitempty"`
}
