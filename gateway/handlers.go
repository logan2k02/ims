package main

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/logan2k02/ims/gateway/products_handlers"
	"github.com/logan2k02/ims/shared/protobuf"

	_ "github.com/joho/godotenv/autoload"
)

var validate = validator.New(validator.WithRequiredStructEnabled())

func registerHandlers(app *fiber.App, productsClient *protobuf.ProductsServiceClient) {
	app.Post("/products/create", products_handlers.CreateProductHandler(productsClient, validate))
	app.Get("/products/:id", products_handlers.GetProduct(productsClient))
	app.Get("/products", products_handlers.ListProducts(productsClient))
	app.Delete("/products/:id", products_handlers.DeleteProduct(productsClient))
	app.Put("/products/:id", products_handlers.UpdateProduct(productsClient, validate))
}
