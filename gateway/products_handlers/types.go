package products_handlers

type createProductDto struct {
	Name            string  `json:"name" validate:"required"`
	Sku             string  `json:"sku" validate:"required"`
	Description     string  `json:"description"`
	Price           float64 `json:"price" validate:"required,gt=0"`
	ReorderLevel    int64   `json:"reorder_level" validate:"required,gt=0"`
	ReorderQuantity int64   `json:"reorder_quantity" validate:"required,gt=0"`
	InitialQuantity int64   `json:"initial_quantity" validate:"required,gt=0"`
}

type updateProductDto struct {
	Name            string  `json:"name" validate:"required"`
	Sku             string  `json:"sku" validate:"required"`
	Description     string  `json:"description"`
	Price           float64 `json:"price" validate:"required,gt=0"`
	ReorderLevel    int64   `json:"reorder_level" validate:"required,gt=0"`
	ReorderQuantity int64   `json:"reorder_quantity" validate:"required,gt=0"`
}

type product struct {
	Id              int64   `json:"id"`
	Name            string  `json:"name"`
	Sku             string  `json:"sku"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	CreatedAt       string  `json:"created_at"`
	ReorderLevel    int64   `json:"reorder_level"`
	ReorderQuantity int64   `json:"reorder_quantity"`
	StockQuantity   int64   `json:"stock_quantity"`
}
