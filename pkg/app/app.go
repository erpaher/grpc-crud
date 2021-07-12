package app

import (
	"context"
	"fmt"

	"github.com/erpaher/grpc-crud/pkg/store"
	grpc "github.com/erpaher/grpc-crud/proto/serviceexampleapi"
)

type app struct {
	grpc.UnimplementedServiceExampleServiceServer
	store Store
}

func New(s Store) *app {
	return &app{store: s}
}

func (app *app) CreateUser(ctx context.Context, in *grpc.CreateUserRequest) (*grpc.User, error) {

	user := store.User{
		Name:     in.GetName(),
		Age:      uint8(in.GetAge()),
		UserType: uint8(in.GetUserType()),
		Items:    make([]store.Item, 0, len(in.GetItems())),
	}

	for _, item := range in.GetItems() {
		user.Items = append(user.Items, store.Item{
			Name: item.Name,
		})
	}

	newUser, err := app.store.CreateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("fail to create user: %w", err)
	}

	return storeToGRPCUser(newUser), nil
}

func (app *app) UpdateUser(ctx context.Context, in *grpc.UpdateUserRequest) (*grpc.User, error) {

	user := store.User{
		ID:       in.GetId(),
		Name:     in.GetName(),
		Age:      uint8(in.GetAge()),
		UserType: uint8(in.GetUserType()),
		Items:    make([]store.Item, 0, len(in.GetItems())),
	}

	for _, item := range in.GetItems() {
		user.Items = append(user.Items, store.Item{
			ID:   item.GetId(),
			Name: item.Name,
		})
	}

	newUser, err := app.store.UpdateUser(ctx, user)
	if err != nil {
		return nil, fmt.Errorf("fail to update user: %w", err)
	}

	return storeToGRPCUser(newUser), nil
}

func (app *app) DeleteUser(ctx context.Context, in *grpc.DeleteUserRequest) (*grpc.DeleteUserResponse, error) {
	userID := in.GetId()
	err := app.store.DeleteUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fail to delete user: %w", err)
	}
	return &grpc.DeleteUserResponse{}, nil
}

func (app *app) ListUser(ctx context.Context, in *grpc.ListUserRequest) (*grpc.ListUserResponse, error) {

	var (
		limit, page uint32
	)

	pf := in.GetPageFilter()

	if pf != nil {
		limit = pf.GetLimit()
		page = pf.GetPage()
	}

	users, err := app.store.ListUser(ctx, page, limit)
	if err != nil {
		return nil, fmt.Errorf("fail to get list of users: %w", err)
	}

	grpcUsers := make([]*grpc.User, 0, len(users))
	for _, user := range users {
		grpcUsers = append(grpcUsers, storeToGRPCUser(user))
	}

	return &grpc.ListUserResponse{Users: grpcUsers}, nil
}

func (app *app) GetUser(ctx context.Context, in *grpc.GetUserRequest) (*grpc.User, error) {
	userID := in.GetId()
	user, err := app.store.GetUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("fail to get user: %w", err)
	}
	return storeToGRPCUser(user), nil

}

func (app *app) CreateItem(ctx context.Context, in *grpc.CreateItemRequest) (*grpc.Item, error) {
	name := in.GetName()
	userID := in.GetUserId()

	item, err := app.store.CreateItem(ctx, name, userID)
	if err != nil {
		return nil, fmt.Errorf("fail to add item: %w", err)
	}
	return storeToGRPCItem(item), nil
}

func (app *app) UpdateItem(ctx context.Context, in *grpc.UpdateItemRequest) (*grpc.Item, error) {
	itemID := in.GetId()
	name := in.GetName()

	item, err := app.store.UpdateItem(ctx, itemID, name)
	if err != nil {
		return nil, fmt.Errorf("fail to update item: %w", err)
	}
	return storeToGRPCItem(item), nil
}
