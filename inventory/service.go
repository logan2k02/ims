package main

import (
	"context"

	pb "github.com/logan2k02/ims/shared/protobuf"
)

type inventoryService struct {
	store *inventoryStore
}

func NewInventoryService(store *inventoryStore) *inventoryService {
	return &inventoryService{store}
}

func (s *inventoryService) Purchase(ctx context.Context, payload *pb.PurchaseInventoryRequest) (*pb.StockMovement, error) {
	return s.store.UpdateStockQuantity(ctx, &UpdateStockDto{
		ProductId: payload.ProductId,
		Change:    -payload.Quantity,
		Reference: payload.Reference,
		Type:      "purchase",
	})
}

func (s *inventoryService) Supply(ctx context.Context, payload *pb.ManageInventoryRequest) (*pb.StockMovement, error) {
	return s.store.UpdateStockQuantity(ctx, &UpdateStockDto{
		ProductId: payload.ProductId,
		Change:    payload.Quantity,
		Note:      payload.Note,
		Type:      "supply",
	})
}

func (s *inventoryService) CorrectStockQuantity(ctx context.Context, payload *pb.ManageInventoryRequest) (*pb.StockMovement, error) {
	return s.store.UpdateStockQuantity(ctx, &UpdateStockDto{
		ProductId: payload.ProductId,
		Change:    payload.Quantity,
		Note:      payload.Note,
		Type:      "correction",
	})
}

func (s *inventoryService) ListStockMovements(ctx context.Context, productId int64) ([]*pb.StockMovement, error) {
	return s.store.ListStockMovements(ctx, productId)
}
