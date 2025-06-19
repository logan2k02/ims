package products_handlers

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/theshawa/imms/shared/protobuf"
)

func CreateProductHandler(productsClient *protobuf.ProductsServiceClient, validate *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var payload createProductDto
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		if err := validate.Struct(payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "validation failed",
				"details": err.Error(),
			})
		}

		productRes, err := (*productsClient).CreateProduct(c.Context(), &protobuf.CreateProductRequest{
			Name:        payload.Name,
			Description: payload.Description,
			Price:       payload.Price,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create product"})
		}

		product := &product{
			Id:          productRes.Id,
			Name:        productRes.Name,
			Description: productRes.Description,
			Price:       productRes.Price,
		}

		return c.Status(fiber.StatusCreated).JSON(product)
	}
}
