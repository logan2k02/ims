package products_handlers

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	pb "github.com/theshawa/ims/shared/protobuf"
)

func CreateProductHandler(productsClient *pb.ProductsServiceClient, validate *validator.Validate) fiber.Handler {
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

		productRes, err := (*productsClient).CreateProduct(c.Context(), &pb.CreateProductRequest{
			Name:        payload.Name,
			Description: payload.Description,
			Price:       payload.Price,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create product", "details": err.Error()})
		}

		product := &product{
			Id:          productRes.Id,
			Name:        productRes.Name,
			Description: productRes.Description,
			Price:       productRes.Price,
			CreatedAt:   productRes.CreatedAt,
		}

		return c.Status(fiber.StatusCreated).JSON(product)
	}
}

func GetProduct(productsClient *pb.ProductsServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id", "1")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "invalid id given",
				"details": "ind must be an integer",
			})
		}

		productRes, err := (*productsClient).GetProduct(c.Context(), &pb.GetProductRequest{
			Id: id,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get product", "details": err.Error()})
		}

		product := &product{
			Id:          productRes.Id,
			Name:        productRes.Name,
			Description: productRes.Description,
			Price:       productRes.Price,
			CreatedAt:   productRes.CreatedAt,
		}

		return c.Status(fiber.StatusCreated).JSON(product)
	}
}
