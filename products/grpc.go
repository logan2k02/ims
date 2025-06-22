package main

import (
	"context"

	pb "github.com/theshawa/ims/shared/protobuf"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type productsGRPCHandler struct {
	service *productsService
	pb.UnimplementedProductsServiceServer
}

func NewProductsGRPCHandler(service *productsService) *productsGRPCHandler {
	return &productsGRPCHandler{
		service: service,
	}
}

func (s *productsGRPCHandler) CreateProduct(ctx context.Context, payload *pb.CreateProductRequest) (*pb.Product, error) {
	product, err := s.service.CreateProduct(ctx, payload)

	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create product: %v", err)
	}

	return product, nil
}

func (s *productsGRPCHandler) GetProduct(ctx context.Context, payload *pb.GetProductRequest) (*pb.Product, error) {
	product, err := s.service.GetProduct(ctx, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get product: %v", err)
	}

	if product == nil {
		return nil, status.Errorf(codes.Internal, "product with given id does not exists")
	}

	return product, nil
}

func (s *productsGRPCHandler) ListProducts(ctx context.Context, payload *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products, err := s.service.GetProducts(ctx, payload)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get products: %v", err)
	}

	return &pb.ListProductsResponse{
		Products: products,
	}, nil
}
