package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/theshawa/ims/gateway/products_handlers"
	"github.com/theshawa/ims/shared/protobuf"

	_ "github.com/joho/godotenv/autoload" // Automatically load .env file
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func registerHandlers(app *fiber.App, productsClient *protobuf.ProductsServiceClient) {
	app.Post("/products/create", products_handlers.CreateProductHandler(productsClient, validate))
}
