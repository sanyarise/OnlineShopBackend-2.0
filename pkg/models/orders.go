package models

import (
	"time"

	"github.com/google/uuid"
)

type Order struct {
	ID            uuid.UUID `json:"id,omitempty"`
	Shipment_time time.Time `json:"shipment___time,omitempty"`
	User          *User     `json:"user,omitempty"`
	Address       string    `json:"address,omitempty"`
	Items         []Item    `json:"items,omitempty"`
}
