package client

import (
	"fmt"
	"log"

	api "github.com/hashicorp/consul/api"
)

type Consul struct {
	client *api.Client
	logger *log.Logger
}

type ServiceInfo struct {
	Id      string
	Name    string
	Address string
	Port    int
	Tags    []string
}

func NewConsulClient(logger *log.Logger, agentAddress string, token string) *Consul {
	config := api.DefaultConfig()
	config.Address = agentAddress
	config.Token = token

	consul, err := api.NewClient(config)
	if err != nil {
		log.Fatalf("[ERROR] Failed to create consul client %v\n", err)
	}
	return &Consul{
		logger: logger,
		client: consul,
	}
}

func (c *Consul) Register(config *ServiceInfo) (string, error) {
	registration := &api.AgentServiceRegistration{
		ID:      config.Id,
		Name:    config.Name,
		Address: config.Address,
		Tags:    config.Tags,
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/hello", config.Address, config.Port),
			Interval: "5s",
			Timeout:  "3s",
		},
	}
	err := c.client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatalf("[ERROR] Failed to register service %v\n", err)
	}
	log.Printf("[INFO] Success to register %s\n", config.Id)
	return config.Id, nil
}

func (c *Consul) Deregister(id string) {
	if err := c.client.Agent().ServiceDeregister(id); err != nil {
		log.Printf("[WARN] Failed to deregister %s\n", id)
	}
	log.Printf("[INFO] Success to deregister %s\n", id)
}
