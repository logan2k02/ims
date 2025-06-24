package orders_handlers

type orderItem struct {
	ProductId int64 `json:"product_id" validate:"required"`
	Quantity  int64 `json:"quantity" validate:"required,gt=0"`
}

type createOrderDto struct {
	CustomerName     string      `json:"customer_name" validate:"required"`
	CustomerContact  string      `json:"customer_contact" validate:"required"`
	Items            []orderItem `json:"items" validate:"required,dive"`
	PaymentReference string      `json:"payment_reference" validate:"required"`
}

type changeOrderStatusDto struct {
	Status string `json:"status" validate:"required,oneof=pending completed cancelled"`
}

type order struct {
	Id               int64       `json:"id"`
	PaymentReference string      `json:"payment_reference"`
	CustomerName     string      `json:"customer_name"`
	CustomerContact  string      `json:"customer_contact"`
	Status           string      `json:"status"`
	CreatedAt        string      `json:"created_at"`
	Items            []orderItem `json:"items"`
}
