package inventory_handlers

type manageDto struct {
	Quantity int64  `json:"quantity" validate:"required,gt=0"`
	Note     string `json:"note"`
}
