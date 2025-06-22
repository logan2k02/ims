package products_handlers

import (
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	pb "github.com/logan2k02/ims/shared/protobuf"
	"google.golang.org/grpc/status"
)

func CreateProductHandler(productsClient pb.ProductsServiceClient, validate *validator.Validate) fiber.Handler {
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

		productRes, err := productsClient.CreateProduct(c.Context(), &pb.CreateProductRequest{
			Name:            payload.Name,
			Sku:             payload.Sku,
			Description:     payload.Description,
			Price:           payload.Price,
			ReorderLevel:    payload.ReorderLevel,
			ReorderQuantity: payload.ReorderQuantity,
			InitialQuantity: payload.InitialQuantity,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create product", "details": status.Convert(err).Message()})
		}

		product := &product{
			Id:              productRes.Id,
			Name:            productRes.Name,
			Sku:             productRes.Sku,
			Description:     productRes.Description,
			Price:           productRes.Price,
			CreatedAt:       productRes.CreatedAt,
			ReorderLevel:    productRes.ReorderLevel,
			ReorderQuantity: payload.ReorderQuantity,
			StockQuantity:   productRes.StockQuantity,
		}

		return c.Status(fiber.StatusCreated).JSON(product)
	}
}

func GetProduct(productsClient pb.ProductsServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id", "1")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "invalid id given",
				"details": "ind must be an integer",
			})
		}

		productRes, err := productsClient.GetProduct(c.Context(), &pb.ProductIdRequest{
			Id: id,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get product", "details": status.Convert(err).Message()})
		}

		product := &product{
			Id:              productRes.Id,
			Name:            productRes.Name,
			Sku:             productRes.Sku,
			Description:     productRes.Description,
			Price:           productRes.Price,
			CreatedAt:       productRes.CreatedAt,
			ReorderLevel:    productRes.ReorderLevel,
			ReorderQuantity: productRes.ReorderQuantity,
			StockQuantity:   productRes.StockQuantity,
		}

		return c.Status(fiber.StatusCreated).JSON(product)
	}
}

func ListProducts(productsClient pb.ProductsServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idsParam := c.Query("ids", "")
		var ids []int64
		if idsParam != "" {
			idStrings := strings.SplitSeq(idsParam, ",")
			for idStr := range idStrings {
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

		productRes, err := productsClient.ListProducts(c.Context(), &pb.ListProductsRequest{
			Ids: ids,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get products", "details": status.Convert(err).Message()})
		}

		var products []*product
		for _, p := range productRes.Products {
			product := &product{
				Id:              p.Id,
				Name:            p.Name,
				Sku:             p.Sku,
				Description:     p.Description,
				Price:           p.Price,
				CreatedAt:       p.CreatedAt,
				ReorderLevel:    p.ReorderLevel,
				ReorderQuantity: p.ReorderQuantity,
				StockQuantity:   p.StockQuantity,
			}
			products = append(products, product)
		}

		return c.Status(fiber.StatusOK).JSON(products)
	}
}

func DeleteProduct(productsClient pb.ProductsServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id", "1")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "invalid id given",
				"details": "id must be an integer",
			})
		}

		if _, err := productsClient.DeleteProduct(c.Context(), &pb.ProductIdRequest{
			Id: id,
		}); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete product", "details": status.Convert(err).Message()})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}

func UpdateProduct(productsClient pb.ProductsServiceClient, validate *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id", "1")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "invalid id given",
				"details": "id must be an integer",
			})
		}

		var payload updateProductDto
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		if err := validate.Struct(payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "validation failed",
				"details": err.Error(),
			})
		}

		productRes, err := productsClient.UpdateProduct(c.Context(), &pb.UpdateProductRequest{
			Id:              id,
			Name:            payload.Name,
			Sku:             payload.Sku,
			Description:     payload.Description,
			Price:           payload.Price,
			ReorderLevel:    payload.ReorderLevel,
			ReorderQuantity: payload.ReorderQuantity,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to update product", "details": status.Convert(err).Message()})
		}

		product := &product{
			Id:              productRes.Id,
			Name:            productRes.Name,
			Sku:             payload.Sku,
			Description:     productRes.Description,
			Price:           productRes.Price,
			CreatedAt:       productRes.CreatedAt,
			ReorderLevel:    productRes.ReorderLevel,
			ReorderQuantity: productRes.ReorderQuantity,
			StockQuantity:   productRes.StockQuantity,
		}

		return c.Status(fiber.StatusOK).JSON(product)
	}
}
