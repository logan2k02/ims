package products_handlers

import (
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	pb "github.com/logan2k02/ims/shared/protobuf"
	"google.golang.org/grpc/status"
)

func CreateProductHandler(productsClient *pb.ProductsServiceClient, validate *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var payload manageProductDto
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
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create product", "details": status.Convert(err).Message()})
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

		productRes, err := (*productsClient).GetProduct(c.Context(), &pb.ProductIdRequest{
			Id: id,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get product", "details": status.Convert(err).Message()})
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

func ListProducts(productsClient *pb.ProductsServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idsParam := c.Query("ids", "")
		var ids []int64
		if idsParam != "" {
			idStrings := strings.Split(idsParam, ",")
			for _, idStr := range idStrings {
				id, err := strconv.ParseInt(idStr, 10, 64)
				if err != nil {
					return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
						"error":   "invalid id given",
						"details": "ids must be a comma-separated list of integers",
					})
				}
				ids = append(ids, id)
			}
		}

		productRes, err := (*productsClient).ListProducts(c.Context(), &pb.ListProductsRequest{
			Ids: ids,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to get products", "details": status.Convert(err).Message()})
		}

		var products []*product
		for _, p := range productRes.Products {
			product := &product{
				Id:          p.Id,
				Name:        p.Name,
				Description: p.Description,
				Price:       p.Price,
				CreatedAt:   p.CreatedAt,
			}
			products = append(products, product)
		}

		return c.Status(fiber.StatusOK).JSON(products)
	}
}

func DeleteProduct(productsClient *pb.ProductsServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id", "1")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "invalid id given",
				"details": "id must be an integer",
			})
		}

		if _, err := (*productsClient).DeleteProduct(c.Context(), &pb.ProductIdRequest{
			Id: id,
		}); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete product", "details": status.Convert(err).Message()})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func UpdateProduct(productsClient *pb.ProductsServiceClient, validate *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id", "1")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "invalid id given",
				"details": "id must be an integer",
			})
		}

		var payload manageProductDto
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		if err := validate.Struct(payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "validation failed",
				"details": err.Error(),
			})
		}

		productRes, err := (*productsClient).UpdateProduct(c.Context(), &pb.UpdateProductRequest{
			Id:          id,
			Name:        payload.Name,
			Description: payload.Description,
			Price:       payload.Price,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update product", "details": status.Convert(err).Message()})
		}

		product := &product{
			Id:          productRes.Id,
			Name:        productRes.Name,
			Description: productRes.Description,
			Price:       productRes.Price,
			CreatedAt:   productRes.CreatedAt,
		}

		return c.Status(fiber.StatusOK).JSON(product)
	}
}
