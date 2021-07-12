package app

import (
	"context"

	"github.com/erpaher/grpc-crud/pkg/store"
)

type Store interface {
	CreateUser(ctx context.Context, user store.User) (*store.User, error)
	UpdateUser(ctx context.Context, user store.User) (*store.User, error)
	DeleteUser(ctx context.Context, userID uint32) error
	ListUser(ctx context.Context, page, limit uint32) ([]*store.User, error)
	GetUser(ctx context.Context, userID uint32) (*store.User, error)
	CreateItem(ctx context.Context, name string, userID uint32) (*store.Item, error)
	UpdateItem(ctx context.Context, itemID uint32, name string) (*store.Item, error)
}
