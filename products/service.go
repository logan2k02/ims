package main

import (
	"context"

	"github.com/theshawa/ims/shared/protobuf"
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
		return nil, err
	}

	return product, nil
}

func (s *productsService) GetProducts(ctx context.Context, payload *protobuf.ListProductsRequest) ([]*protobuf.Product, error) {
	products, err := s.store.GetProducts(ctx, payload.Ids)
	if err != nil {
		return nil, err
	}
	return products, nil
}

func (s *productsService) GetProduct(ctx context.Context, payload *protobuf.GetProductRequest) (*protobuf.Product, error) {
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
