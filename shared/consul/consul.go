package consul

import (
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
)

type Client struct {
	client     *api.Client
	Addr       string
	rrIndexMap map[string]int
}

func NewClient(consulAddr string) (*Client, error) {
	config := api.DefaultConfig()
	config.Address = consulAddr

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create consul client: %w", err)
	}

	return &Client{
		client:     client,
		Addr:       consulAddr,
		rrIndexMap: make(map[string]int),
	}, nil
}

func (c *Client) RegisterGRPCService(serviceName string, serviceHost string, servicePort int) error {
	healthCheckHost := serviceHost

	if serviceHost == "127.0.0.1" || serviceHost == "localhost" {
		log.Printf("using 'host.docker.internal' for health check address for service %s as consul is running inside docker container\n", serviceName)
		healthCheckHost = "host.docker.internal"
	}

	registration := &api.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s-%s-%d", serviceName, serviceHost, servicePort),
		Name:    serviceName,
		Address: serviceHost,
		Port:    servicePort,
		Tags:    []string{"grpc", serviceName},
		Check: &api.AgentServiceCheck{
			GRPC:                           fmt.Sprintf("%s:%d", healthCheckHost, servicePort),
			Interval:                       "10m",
			DeregisterCriticalServiceAfter: "1m",
			Timeout:                        "10s",
		},
	}

	if err := c.client.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("failed to register service with consul: %w", err)
	}

	fmt.Printf("service %s registered successfully with consul at %s:%d\n", serviceName, serviceHost, servicePort)

	return nil
}

func (c *Client) RegisterHTTPService(serviceName string, serviceHost string, servicePort int) error {
	healthCheckHost := serviceHost

	if serviceHost == "127.0.0.1" || serviceHost == "localhost" {
		log.Printf("using 'host.docker.internal' for health check address for service %s as consul is running inside docker container\n", serviceName)
		healthCheckHost = "host.docker.internal"
	}

	registration := &api.AgentServiceRegistration{
		ID:      serviceName,
		Name:    serviceName,
		Address: serviceHost,
		Port:    servicePort,
		Tags:    []string{"http", serviceName},
		Check: &api.AgentServiceCheck{
			HTTP:                           fmt.Sprintf("http://%s:%d/health", healthCheckHost, servicePort),
			Interval:                       "10s",
			DeregisterCriticalServiceAfter: "1m",
			Timeout:                        "10s",
		},
	}

	if err := c.client.Agent().ServiceRegister(registration); err != nil {
		return fmt.Errorf("failed to register HTTP service with consul: %w", err)
	}

	fmt.Printf("HTTP service %s registered successfully with consul at %s:%d\n", serviceName, serviceHost, servicePort)

	return nil
}

func (c *Client) DiscoverService(serviceName string) (string, error) {
	services, _, err := c.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return "", fmt.Errorf("failed to discover service %s: %w", serviceName, err)
	}

	if len(services) == 0 {
		return "", fmt.Errorf("no healthy instances found for service %s", serviceName)
	}

	index := c.rrIndexMap[serviceName]
	selected := services[index%len(services)]
	c.rrIndexMap[serviceName] = (index + 1) % len(services)

	target := fmt.Sprintf("%s:%d", selected.Service.Address, selected.Service.Port)

	log.Printf("discovered service %s at %s", serviceName, target)

	return target, nil
}

func (c *Client) DeregisterService(serviceName string, serviceHost string, servicePort int) error {
	return c.client.Agent().ServiceDeregister(fmt.Sprintf("%s-%s-%d", serviceName, serviceHost, servicePort))
}
