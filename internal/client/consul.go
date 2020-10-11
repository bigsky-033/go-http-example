package client

import (
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
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

func (c *Consul) Register(serviceInfo *ServiceInfo) (string, error) {
	registration := &api.AgentServiceRegistration{
		ID:      serviceInfo.Id,
		Name:    serviceInfo.Name,
		Address: serviceInfo.Address,
		Tags:    serviceInfo.Tags,
		Port:    serviceInfo.Port,
		Check: &api.AgentServiceCheck{
			HTTP:     fmt.Sprintf("http://%s:%d/hello", serviceInfo.Address, serviceInfo.Port),
			Interval: "5s",
			Timeout:  "3s",
		},
	}
	err := c.client.Agent().ServiceRegister(registration)
	if err != nil {
		log.Fatalf("[ERROR] Failed to register service %v\n", err)
	}
	log.Printf("[INFO] Success to register %s\n", serviceInfo.Id)
	return serviceInfo.Id, nil
}

func (c *Consul) Deregister(id string) {
	if err := c.client.Agent().ServiceDeregister(id); err != nil {
		log.Printf("[WARN] Failed to deregister %s\n", id)
	}
	log.Printf("[INFO] Success to deregister %s\n", id)
}
