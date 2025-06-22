package grpcservice

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/theshawa/ims/shared/consul"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type server struct {
	gRPCServer   *grpc.Server
	consulClient *consul.Client
	ServiceName  string
	host         string
	port         int
}

func NewServer(consulClient *consul.Client, serviceName string, host string, port int) (*server, error) {
	gRPCServer := grpc.NewServer()
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(gRPCServer, healthServer)

	healthServer.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	healthServer.SetServingStatus(serviceName, grpc_health_v1.HealthCheckResponse_SERVING)

	reflection.Register(gRPCServer)

	return &server{
		gRPCServer:   gRPCServer,
		consulClient: consulClient,
		ServiceName:  serviceName,
		host:         host,
		port:         port,
	}, nil
}

func (s *server) Start() error {
	if err := s.consulClient.RegisterGRPCService(s.ServiceName, s.host, s.port); err != nil {
		return fmt.Errorf("failed to register service with consul: %w", err)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.host, s.port))
	if err != nil {
		return fmt.Errorf("failed to listen on %s:%d: %w", s.host, s.port, err)
	}

	go s.handleGracefulShutdown()

	return s.gRPCServer.Serve(lis)
}

func (s *server) handleGracefulShutdown() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	log.Println("shutting down server gracefully...")

	if err := s.consulClient.DeregisterService(s.ServiceName, s.host, s.port); err != nil {
		log.Printf("failed to deregister service with consul: %v", err)
	} else {
		log.Printf("service %s deregistered from consul", s.ServiceName)
	}

	s.gRPCServer.GracefulStop()
}

func (s *server) RegisterService(desc *grpc.ServiceDesc, impl any) {
	s.gRPCServer.RegisterService(desc, impl)
}
