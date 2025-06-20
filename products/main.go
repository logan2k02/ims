package main

import (
	"log"
	"strconv"

	"github.com/theshawa/imms/shared/consul"
	"github.com/theshawa/imms/shared/grpcservice"
	"github.com/theshawa/imms/shared/protobuf"
	"github.com/theshawa/imms/shared/utils"

	_ "github.com/joho/godotenv/autoload"
)

var (
	gRPCPort   = utils.GetEnv("GRPC_PORT", "50051")
	gRPCHost   = utils.GetEnv("GRPC_HOST", "localhost")
	consulAddr = utils.GetEnv("CONSUL_ADDR", "localhost:8500")
)

func main() {
	store, err := NewProductsStore()
	if err != nil {
		log.Fatalf("failed to create products store: %v", err)
	}
	defer func() {
		if err := store.Close(); err != nil {
			log.Fatalf("failed to close products store: %v", err)
		}
		log.Println("products store closed successfully")
	}()

	if err := store.Init(); err != nil {
		log.Fatalf("failed to initialize products store: %v", err)
	}

	log.Println("products store initialized successfully")

	service := NewProductsService(store)

	consulCient, err := consul.NewClient(consulAddr)
	if err != nil {
		log.Fatalf("failed to create consul client: %v", err)
	}

	_gRPCPort, _ := strconv.Atoi(gRPCPort)

	gRPCServiceServer, err := grpcservice.NewServer(consulCient, "products-grpc-service", gRPCHost, _gRPCPort)
	if err != nil {
		log.Fatalf("failed to create gRPC server for products-service: %v", err)
	}

	productsGRPCHandler := NewProductsGRPCHandler(service)
	gRPCServiceServer.RegisterService(&protobuf.ProductsService_ServiceDesc, productsGRPCHandler)

	log.Printf("starting gRPC server for products-service on port %s", gRPCPort)

	if err := gRPCServiceServer.Start(); err != nil {
		log.Fatalf("failed to start gRPC server for products-service: %v", err)
	}
}
