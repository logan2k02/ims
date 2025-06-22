package main

import (
	"context"

	pb "github.com/logan2k02/ims/shared/protobuf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type inventoryGRPCHandler struct {
	service *inventoryService
	pb.UnimplementedInventoryServiceServer
}

func NewInventoryGRPCHandler(service *inventoryService) *inventoryGRPCHandler {
	return &inventoryGRPCHandler{
		service: service,
	}
}

func (h *inventoryGRPCHandler) PurchaseInventoryProduct(ctx context.Context, payload *pb.PurchaseInventoryRequest) (*pb.StockMovement, error) {
	record, err := h.service.Purchase(ctx, payload)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return record, nil
}

func (h *inventoryGRPCHandler) SupplyInventoryProduct(ctx context.Context, payload *pb.ManageInventoryRequest) (*pb.StockMovement, error) {
	record, err := h.service.Supply(ctx, payload)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return record, nil
}

func (h *inventoryGRPCHandler) CorrectInventoryStock(ctx context.Context, payload *pb.ManageInventoryRequest) (*pb.StockMovement, error) {
	record, err := h.service.CorrectStockQuantity(ctx, payload)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return record, nil
}

func (h *inventoryGRPCHandler) ListStockMovements(ctx context.Context, payload *pb.ListStockMovementsRequest) (*pb.ListStockMovementsResponse, error) {
	records, err := h.service.ListStockMovements(ctx, payload.ProductId)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ListStockMovementsResponse{
		Records: records,
	}, nil
}
