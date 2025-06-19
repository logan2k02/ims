package main

import (
	"context"

	"github.com/theshawa/imms/shared/protobuf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type productsGRPCHandler struct {
	service *productsService
	protobuf.UnimplementedProductsServiceServer
}

func NewProductsGRPCHandler(service *productsService) *productsGRPCHandler {
	return &productsGRPCHandler{
		service: service,
	}
}

func (s *productsGRPCHandler) CreateProduct(ctx context.Context, payload *protobuf.CreateProductRequest) (*protobuf.Product, error) {
	product, err := s.service.CreateProduct(ctx, payload)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}

	return product, nil
}
