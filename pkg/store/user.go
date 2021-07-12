package store

import (
	"time"
)

type User struct {
	ID        uint32
	Name      string
	Age       uint8
	UserType  uint8
	Items     []Item
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
