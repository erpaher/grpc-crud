package app

import (
	"github.com/erpaher/grpc-crud/pkg/store"
	grpc "github.com/erpaher/grpc-crud/proto/serviceexampleapi"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func storeToGRPCUser(user *store.User) *grpc.User {

	grpcUser := &grpc.User{
		Id:       user.ID,
		Name:     user.Name,
		Age:      int32(user.Age),
		UserType: grpc.UserType(user.UserType),
		Items:    make([]*grpc.Item, 0, len(user.Items)),
		CreatedAt: &timestamppb.Timestamp{
			Seconds: user.CreatedAt.Unix(),
			Nanos:   0,
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: user.UpdatedAt.Unix(),
			Nanos:   0,
		},
	}

	for _, item := range user.Items {
		grpcUser.Items = append(grpcUser.Items, &grpc.Item{
			Id:     item.ID,
			Name:   item.Name,
			UserId: item.UserID,
			CreatedAt: &timestamppb.Timestamp{
				Seconds: item.CreatedAt.Unix(),
				Nanos:   0,
			},
			UpdatedAt: &timestamppb.Timestamp{
				Seconds: item.UpdatedAt.Unix(),
				Nanos:   0,
			},
		})
	}
	return grpcUser
}

func grpcToStoreUser(user *grpc.User) *store.User {
	storeUser := &store.User{
		ID:       user.GetId(),
		Name:     user.GetName(),
		Age:      uint8(user.GetAge()),
		UserType: uint8(user.GetUserType()),
		Items:    make([]store.Item, 0, len(user.GetItems())),
	}

	for _, item := range user.GetItems() {
		storeUser.Items = append(storeUser.Items, store.Item{
			ID:   item.GetId(),
			Name: item.Name,
		})
	}
	return storeUser
}

func storeToGRPCItem(item *store.Item) *grpc.Item {

	grpcItem := &grpc.Item{
		Id:     item.ID,
		Name:   item.Name,
		UserId: item.UserID,
		CreatedAt: &timestamppb.Timestamp{
			Seconds: item.CreatedAt.Unix(),
			Nanos:   0,
		},
		UpdatedAt: &timestamppb.Timestamp{
			Seconds: item.UpdatedAt.Unix(),
			Nanos:   0,
		},
	}

	return grpcItem
}
