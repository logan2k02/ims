package main

import (
	"context"

	pb "github.com/logan2k02/ims/shared/protobuf"
)

type productsService struct {
	store *productsStore
}

func NewProductsService(store *productsStore) *productsService {
	return &productsService{store}
}

func (s *productsService) CreateProduct(ctx context.Context, payload *pb.CreateProductRequest) (*pb.Product, error) {
	product, err := s.store.CreateProduct(ctx, payload)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productsService) GetProducts(ctx context.Context, payload *pb.ListProductsRequest) ([]*pb.Product, error) {
	products, err := s.store.GetProducts(ctx, payload.Ids)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *productsService) GetProduct(ctx context.Context, payload *pb.ProductIdRequest) (*pb.Product, error) {
	products, err := s.store.GetProducts(ctx, []int64{
		payload.Id,
	})

	if err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return nil, nil
	}

	return products[0], nil
}

func (s *productsService) DeleteProduct(ctx context.Context, payload *pb.ProductIdRequest) error {
	return s.store.DeleteProduct(ctx, payload.Id)
}

func (s *productsService) UpdateProduct(ctx context.Context, payload *pb.UpdateProductRequest) (*pb.Product, error) {
	product, err := s.store.UpdateProduct(ctx, payload)
	if err != nil {
		return nil, err
	}

	return product, nil
}
