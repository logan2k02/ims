package orders_handlers

import (
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	pb "github.com/logan2k02/ims/shared/protobuf"
	"google.golang.org/grpc/status"
)

func CreateOrderHandler(ordersClient pb.OrdersServiceClient, validate *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var payload createOrderDto
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		if err := validate.Struct(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "validation failed",
				"details": err.Error(),
			})
		}

		var items []*pb.OrderItem
		for _, item := range payload.Items {
			items = append(items, &pb.OrderItem{
				ProductId: item.ProductId,
				Quantity:  item.Quantity,
			})
		}

		orderRes, err := ordersClient.CreateOrder(c.Context(), &pb.CreateOrderRequest{
			Items:            items,
			PaymentReference: payload.PaymentReference,
			CustomerName:     payload.CustomerName,
			CustomerContact:  payload.CustomerContact,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to create order", "details": status.Convert(err).Message()})
		}

		var orderResItems []orderItem
		for _, item := range orderRes.Items {
			orderResItems = append(orderResItems, orderItem{
				ProductId: item.ProductId,
				Quantity:  item.Quantity,
			})
		}

		return c.Status(fiber.StatusCreated).JSON(order{
			Id:               orderRes.Id,
			Items:            orderResItems,
			PaymentReference: orderRes.PaymentReference,
			CustomerName:     orderRes.CustomerName,
			CustomerContact:  orderRes.CustomerContact,
			Status:           orderRes.Status,
			CreatedAt:        orderRes.CreatedAt,
		})
	}
}

func GetOrderHandler(ordersClient pb.OrdersServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id", "")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "missing order ID",
				"details": "order ID is required",
			})
		}

		orderId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "invalid order ID",
				"details": "order ID must be an integer",
			})
		}

		orderRes, err := ordersClient.GetOrder(c.Context(), &pb.OrderIdRequest{
			Id: orderId,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to get order", "details": status.Convert(err).Message()})
		}

		var items []orderItem
		for _, item := range orderRes.Items {
			items = append(items, orderItem{
				ProductId: item.ProductId,
				Quantity:  item.Quantity,
			})
		}

		return c.Status(fiber.StatusOK).JSON(order{
			Id:               orderRes.Id,
			Items:            items,
			PaymentReference: orderRes.PaymentReference,
			CustomerName:     orderRes.CustomerName,
			CustomerContact:  orderRes.CustomerContact,
			Status:           orderRes.Status,
			CreatedAt:        orderRes.CreatedAt,
		})
	}
}

func ListOrdersHandler(ordersClient pb.OrdersServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		pageStr := c.Query("page", "1")
		pageSizeStr := c.Query("page_size", "10")

		page, err := strconv.ParseInt(pageStr, 10, 64)
		if err != nil || page < 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "invalid page number",
				"details": "page must be a positive integer",
			})
		}

		pageSize, err := strconv.ParseInt(pageSizeStr, 10, 64)
		if err != nil || pageSize < 1 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "invalid limit number",
				"details": "limit must be a positive integer",
			})
		}

		listRes, err := ordersClient.ListOrders(c.Context(), &pb.ListOrdersRequest{
			Page:     page,
			PageSize: pageSize,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to list orders", "details": status.Convert(err).Message()})
		}

		var orders []order
		for _, o := range listRes.Orders {
			var items []orderItem
			for _, item := range o.Items {
				items = append(items, orderItem{
					ProductId: item.ProductId,
					Quantity:  item.Quantity,
				})
			}
			orders = append(orders, order{
				Id:               o.Id,
				Items:            items,
				PaymentReference: o.PaymentReference,
				CustomerName:     o.CustomerName,
				CustomerContact:  o.CustomerContact,
				Status:           o.Status,
				CreatedAt:        o.CreatedAt,
			})
		}

		return c.Status(fiber.StatusOK).JSON(orders)
	}
}

func ChangeOrderStatusHandler(ordersClient pb.OrdersServiceClient, validate *validator.Validate) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var payload changeOrderStatusDto
		if err := c.BodyParser(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid request body"})
		}

		if err := validate.Struct(&payload); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "validation failed",
				"details": err.Error(),
			})
		}

		orderIdstr := c.Params("id", "")
		if orderIdstr == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "missing order ID",
				"details": "order ID is required",
			})
		}

		orderId, err := strconv.ParseInt(orderIdstr, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "invalid order ID",
				"details": "order ID must be an integer",
			})
		}

		orderRes, err := ordersClient.ChangeOrderStatus(c.Context(), &pb.ChangeOrderStatusRequest{
			Id:     orderId,
			Status: payload.Status,
		})
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to change order status", "details": status.Convert(err).Message()})
		}

		var items []orderItem
		for _, item := range orderRes.Items {
			items = append(items, orderItem{
				ProductId: item.ProductId,
				Quantity:  item.Quantity,
			})
		}

		return c.Status(fiber.StatusOK).JSON(order{
			Id:               orderRes.Id,
			Items:            items,
			PaymentReference: orderRes.PaymentReference,
			CustomerName:     orderRes.CustomerName,
			CustomerContact:  orderRes.CustomerContact,
			Status:           orderRes.Status,
			CreatedAt:        orderRes.CreatedAt,
		})
	}
}

func DeleteOrderHandler(ordersClient pb.OrdersServiceClient) fiber.Handler {
	return func(c *fiber.Ctx) error {
		id := c.Params("id", "")
		if id == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "missing order ID",
				"details": "order ID is required",
			})
		}

		orderId, err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "invalid order ID",
				"details": "order ID must be an integer",
			})
		}

		if _, err := ordersClient.DeleteOrder(c.Context(), &pb.OrderIdRequest{
			Id: orderId,
		}); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to delete order", "details": status.Convert(err).Message()})
		}

		return c.SendStatus(fiber.StatusNoContent)
	}
}
