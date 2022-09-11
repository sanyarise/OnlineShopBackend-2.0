package models

import (
	"time"

	"github.com/google/uuid"
)

type Cart struct {
	ID         uuid.UUID `json:"id,omitempty"`
	Expires_at time.Time `json:"expires___at,omitempty"`
	User       *User     `json:"user,omitempty"`
	Items      []Item    `json:"items,omitempty"`
}
