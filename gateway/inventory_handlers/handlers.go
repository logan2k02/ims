package inventory_handlers

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	pb "github.com/logan2k02/ims/shared/protobuf"
	"google.golang.org/grpc/status"
)

func Supply(inventoryCLient pb.InventoryServiceClient, validate *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id", "1")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "invalid id given",
				"details": "ind must be an integer",
			})
		}

		var payload manageDto
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		if err := validate.Struct(payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "validation failed",
				"details": err.Error(),
			})
		}

		record, err := inventoryCLient.SupplyInventoryProduct(c.Context(), &pb.ManageInventoryRequest{
			ProductId: id,
			Quantity:  payload.Quantity,
			Note:      payload.Note,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to supply inventory", "details": status.Convert(err).Message()})
		}

		return c.Status(fiber.StatusCreated).JSON(record)
	}
}

func Correct(inventoryCLient pb.InventoryServiceClient, validate *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
		idParam := c.Params("id", "1")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "invalid id given",
				"details": "ind must be an integer",
			})
		}

		var payload manageDto
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		if err := validate.Struct(payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "validation failed",
				"details": err.Error(),
			})
		}

		record, err := inventoryCLient.CorrectInventoryStock(c.Context(), &pb.ManageInventoryRequest{
			ProductId: id,
			Quantity:  payload.Quantity,
			Note:      payload.Note,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to correct stock inventory", "details": status.Convert(err).Message()})
		}

		return c.Status(fiber.StatusCreated).JSON(record)
	}
}
