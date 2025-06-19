package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/theshawa/imms/shared/consul"
	"github.com/theshawa/imms/shared/grpcservice"
	"github.com/theshawa/imms/shared/protobuf"
	"github.com/theshawa/imms/shared/utils"
)

var (
	port       = utils.GetEnv("PORT", "3000")                  // Default HTTP port
	consulAddr = utils.GetEnv("CONSUL_ADDR", "localhost:8500") // Default Consul address
)

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	consulClient, err := consul.NewClient(consulAddr)
	if err != nil {
		log.Fatalf("failed to create consul client: %v", err)
	}

	productsClientConn, err := grpcservice.GetGRPCConnection(consulClient, "products-grpc-service")
	if err != nil {
		log.Fatalf("failed to get gRPC connection for products-grpc-service: %v", err)
	}
	defer productsClientConn.Close()

	productsClient := protobuf.NewProductsServiceClient(productsClientConn)

	registerHandlers(app, &productsClient)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("failed to start HTTP server: %v", err)
	}
}
