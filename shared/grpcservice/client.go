package grpcservice

import (
	"fmt"

	"github.com/logan2k02/ims/shared/consul"
	"google.golang.org/grpc"

	"google.golang.org/grpc/credentials/insecure"
)

func GetGRPCConnection(consulClient *consul.Client, serviceName string) (*grpc.ClientConn, error) {
	addr, err := consulClient.DiscoverService(serviceName)
	if err != nil {
		return nil, fmt.Errorf("failed to discover service %s: %w", serviceName, err)
	}

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`))
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc client connection: %w", err)
	}

	return conn, nil
}
