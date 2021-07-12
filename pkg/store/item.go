package store

import (
	"time"
)

type Item struct {
	ID        uint32
	Name      string
	UserID    uint32
	CreatedAt *time.Time
	UpdatedAt *time.Time
}
