package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/logan2k02/ims/shared/consul"
	"github.com/logan2k02/ims/shared/grpcservice"
	"github.com/logan2k02/ims/shared/logger"
	"github.com/logan2k02/ims/shared/protobuf"
	"github.com/logan2k02/ims/shared/utils"
)

var (
	port       = utils.GetEnv("PORT", "3000")                  // Default HTTP port
	consulAddr = utils.GetEnv("CONSUL_ADDR", "localhost:8500") // Default Consul address

	Logger = logger.NewLogger("gateway service")
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	consulClient, err := consul.NewClient(consulAddr)
	if err != nil {
		Logger.FatalLog("consul init", "failed to create client: %v", err)
	}

	productsClientConn, err := grpcservice.GetGRPCConnection(consulClient, "products-grpc-service")
	if err != nil {
		Logger.FatalLog("get products client connection", "failed to get gRPC connection: %v", err)
	}
	defer productsClientConn.Close()

	productsClient := protobuf.NewProductsServiceClient(productsClientConn)

	inventoryClientConn, err := grpcservice.GetGRPCConnection(consulClient, "inventory-grpc-service")
	if err != nil {
		Logger.FatalLog("get inventory client connection", "failed to get gRPC connection: %v", err)
	}
	defer inventoryClientConn.Close()

	inventoryClient := protobuf.NewInventoryServiceClient(inventoryClientConn)

	ordersClientConn, err := grpcservice.GetGRPCConnection(consulClient, "orders-grpc-service")
	if err != nil {
		Logger.FatalLog("get orders client connection", "failed to get gRPC connection: %v", err)
	}
	defer ordersClientConn.Close()

	ordersClient := protobuf.NewOrdersServiceClient(ordersClientConn)

	registerHandlers(app, productsClient, inventoryClient, ordersClient)

	if err := app.Listen(":" + port); err != nil {
		Logger.FatalLog("http server init", "failed to start HTTP server: %v", err)
	}
}
