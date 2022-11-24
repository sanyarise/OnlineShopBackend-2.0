package models

import "github.com/google/uuid"

type Category struct {
	Id uuid.UUID `json:"id,omitempty"`

	Name string `json:"name,omitempty"`

	Description string `json:"description,omitempty"`
}
