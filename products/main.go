package main

import (
	"strconv"

	"github.com/logan2k02/ims/shared/consul"
	"github.com/logan2k02/ims/shared/grpcservice"
	"github.com/logan2k02/ims/shared/logger"
	"github.com/logan2k02/ims/shared/protobuf"
	"github.com/logan2k02/ims/shared/utils"

	_ "github.com/joho/godotenv/autoload"
)

var (
	gRPCPort   = utils.GetEnv("GRPC_PORT", "50051")
	gRPCHost   = utils.GetEnv("GRPC_HOST", "localhost")
	consulAddr = utils.GetEnv("CONSUL_ADDR", "localhost:8500")

	Logger = logger.NewLogger("products-gRPC")
)

func main() {
	store, err := NewProductsStore()
	if err != nil {
		Logger.FatalLog("store init", "failed to create products store: %v", err)
	}
	defer func() {
		if err := store.Close(); err != nil {
			Logger.FatalLog("store close", "%v", err)
		}
		Logger.Log("store close", "products store closed successfully")
	}()

	if err := store.Init(); err != nil {
		Logger.FatalLog("store init", "failed to init: %v", err)
	}

	Logger.Log("store init", "initialized successfully")

	service := NewProductsService(store)

	consulCient, err := consul.NewClient(consulAddr)
	if err != nil {
		Logger.FatalLog("consul init", "failed to create client: %v", err)
	}

	_gRPCPort, _ := strconv.Atoi(gRPCPort)

	gRPCServiceServer, err := grpcservice.NewServer(consulCient, "products-grpc-service", gRPCHost, _gRPCPort)
	if err != nil {
		Logger.FatalLog("grpc server init", "failed to create server: %v", err)
	}

	productsGRPCHandler := NewProductsGRPCHandler(service)
	gRPCServiceServer.RegisterService(&protobuf.ProductsService_ServiceDesc, productsGRPCHandler)

	Logger.Log("grpc server init", "starting server on port %s", gRPCPort)

	if err := gRPCServiceServer.Start(); err != nil {
		Logger.FatalLog("grpc server init", "failed to start: %v", err)
	}
}
