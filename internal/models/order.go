package models

import (
	"time"

	"github.com/google/uuid"
)

// CREATE TABLE orders (
//     id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
//     shipment_time timestamp not NULL,
//     user_id UUID,
//     address TEXT,
//     CONSTRAINT fk_user_id
//         FOREIGN KEY(user_id) REFERENCES users(id)
// );

type Status string

const (
	Processing Status = "order processing"
	Processed  Status = "order processed"
	Ready      Status = "ready for shipment"
	Courier    Status = "picked by courier"
	Shipped    Status = "delivered"
)

type Order struct {
	ID           uuid.UUID
	ShipmentTime time.Time
	User         User
	Address      UserAddress
	Status       Status
}
