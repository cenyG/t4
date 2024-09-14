package consul

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"log/slog"
	"strconv"
)

type ConsulWrapper interface {
	RegisterService(serviceName, servicePort, postfix, healthPort string) error
	GetServices(serviceName string) ([]*consulapi.ServiceEntry, error)
}

type consulWrapper struct {
	client *consulapi.Client
}

func NewConsulWrapper() (ConsulWrapper, error) {
	config := consulapi.DefaultConfig()
	client, err := consulapi.NewClient(config)
	if err != nil {
		return nil, errors.Wrap(err, "Error creating Consul client")
	}

	return &consulWrapper{client: client}, nil
}

// RegisterService - register service to Consul
func (c *consulWrapper) RegisterService(serviceName, servicePort, postfix, healthPort string) error {
	servicePortInt, err := strconv.Atoi(servicePort)
	if err != nil {
		return fmt.Errorf("failed to parse port %s: %v", servicePort, err)
	}

	dockerName := fmt.Sprintf("%s%s", serviceName, postfix)
	postfix = lo.Ternary(postfix == "", servicePort, postfix)
	healthUrl := fmt.Sprintf("http://%s:%s/health", dockerName, healthPort)
	err = c.client.Agent().ServiceRegister(&consulapi.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s.%s", serviceName, postfix),
		Name:    serviceName,
		Address: dockerName,
		Port:    servicePortInt,
		Check: &consulapi.AgentServiceCheck{
			HTTP:                           healthUrl,
			Interval:                       "10s",
			Timeout:                        "5s",
			DeregisterCriticalServiceAfter: "1m",
		},
	})
	if err != nil {
		return fmt.Errorf("failed to register service with Consul: %v", err)
	}

	slog.Info(fmt.Sprintf("service %s.%s:%s registered with Consul", serviceName, postfix, servicePort))
	slog.Info(fmt.Sprintf("health url: %s", healthUrl))

	return nil
}

// GetServices - return healthy services by name
func (c *consulWrapper) GetServices(serviceName string) ([]*consulapi.ServiceEntry, error) {
	services, _, err := c.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}

	return services, nil
}
