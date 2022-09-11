package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Email    string    `json:"email,omitempty"`
	Password string    `json:"password,omitempty"`
	Address  string    `json:"address,omitempty"`
	Rights   Right     `json:"rights,omitempty"` // ? The highest level, or the list of rights
}
