package main

import (
	"context"

	"github.com/theshawa/imms/shared/protobuf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type productsService struct {
	store *productsStore
}

func NewProductsService(store *productsStore) *productsService {
	return &productsService{store}
}

func (s *productsService) CreateProduct(ctx context.Context, payload *protobuf.CreateProductRequest) (*protobuf.Product, error) {
	product, err := s.store.CreateProduct(ctx, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}

	return product, nil
}
