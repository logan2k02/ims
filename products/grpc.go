package main

import (
	"context"

	pb "github.com/logan2k02/ims/shared/protobuf"
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

func (h *productsGRPCHandler) CreateProduct(ctx context.Context, payload *pb.CreateProductRequest) (*pb.Product, error) {
	product, err := h.service.CreateProduct(ctx, payload)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return product, nil
}

func (h *productsGRPCHandler) GetProduct(ctx context.Context, payload *pb.ProductIdRequest) (*pb.Product, error) {
	product, err := h.service.GetProduct(ctx, payload)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	if product == nil {
		return nil, status.Errorf(codes.Internal, "product with given id does not exists")
	}

	return product, nil
}

func (h *productsGRPCHandler) ListProducts(ctx context.Context, payload *pb.ListProductsRequest) (*pb.ListProductsResponse, error) {
	products, err := h.service.GetProducts(ctx, payload)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.ListProductsResponse{
		Products: products,
	}, nil
}

func (h *productsGRPCHandler) DeleteProduct(ctx context.Context, payload *pb.ProductIdRequest) (*pb.DeleteProductResponse, error) {
	if err := h.service.DeleteProduct(ctx, payload); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.DeleteProductResponse{}, nil
}

func (h *productsGRPCHandler) UpdateProduct(ctx context.Context, payload *pb.UpdateProductRequest) (*pb.Product, error) {
	product, err := h.service.UpdateProduct(ctx, payload)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return product, nil
}
