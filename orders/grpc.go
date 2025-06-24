package main

import (
	"context"

	pb "github.com/logan2k02/ims/shared/protobuf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ordersGRPCHandler struct {
	service *ordersService
	pb.UnimplementedOrdersServiceServer
}

func NewOrdersGRPCHandler(service *ordersService) *ordersGRPCHandler {
	return &ordersGRPCHandler{
		service: service,
	}
}

func (h *ordersGRPCHandler) CreateOrder(ctx context.Context, payload *pb.CreateOrderRequest) (*pb.Order, error) {
	order, err := h.service.CreateOrder(ctx, payload)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return order, nil
}

func (h *ordersGRPCHandler) GetOrder(ctx context.Context, payload *pb.OrderIdRequest) (*pb.Order, error) {
	order, err := h.service.GetOrder(ctx, payload)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if order == nil {
		return nil, status.Error(codes.NotFound, "order not found")
	}

	return order, nil
}

func (h *ordersGRPCHandler) ListOrders(ctx context.Context, payload *pb.ListOrdersRequest) (*pb.ListOrdersResponse, error) {
	orders, err := h.service.ListOrders(ctx, payload)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return orders, nil
}

func (h *ordersGRPCHandler) DeleteOrder(ctx context.Context, payload *pb.OrderIdRequest) (*pb.DeleteOrderResponse, error) {
	orders, err := h.service.DeleteOrder(ctx, payload)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return orders, nil
}

func (h *ordersGRPCHandler) ChangeOrderStatus(ctx context.Context, payload *pb.ChangeOrderStatusRequest) (*pb.Order, error) {
	record, err := h.service.ChangeOrderStatus(ctx, payload)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return record, nil
}
