package main

import (
	"context"

	pb "github.com/logan2k02/ims/shared/protobuf"
)

type ordersService struct {
	store *ordersStore
}

func NewOrdersService(store *ordersStore) *ordersService {
	return &ordersService{
		store: store,
	}
}

func (s *ordersService) CreateOrder(ctx context.Context, payload *pb.CreateOrderRequest) (*pb.Order, error) {
	order, err := s.store.CreateOrder(ctx, payload)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *ordersService) GetOrder(ctx context.Context, payload *pb.OrderIdRequest) (*pb.Order, error) {
	order, err := s.store.GetOrder(ctx, payload)
	if err != nil {
		return nil, err
	}
	return order, nil
}

func (s *ordersService) ListOrders(ctx context.Context, payload *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, err := s.store.ListOrders(ctx, payload)
	if err != nil {
		return nil, err
	}
	return orders, nil
}

func (s *ordersService) DeleteOrder(ctx context.Context, payload *pb.OrderIdRequest) (*pb.DeleteOrderResponse, error) {
	if err := s.store.DeleteOrder(ctx, payload); err != nil {
		return nil, err
	}
	return &pb.DeleteOrderResponse{}, nil
}

func (s *ordersService) ChangeOrderStatus(ctx context.Context, payload *pb.ChangeOrderStatusRequest) (*pb.Order, error) {
	order, err := s.store.ChangeOrderStatus(ctx, payload)
	if err != nil {
		return nil, err
	}
	return order, nil
}
