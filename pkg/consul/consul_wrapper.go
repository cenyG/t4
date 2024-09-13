package consul

import (
	"fmt"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
	"github.com/samber/lo"
	"log/slog"
	"net"
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

func (c *consulWrapper) RegisterService(serviceName, servicePort, postfix, healthPort string) error {
	// get local ip
	ip, err := getLocalIP()
	if err != nil {
		return fmt.Errorf("failed to get local IP: %v", err)
	}

	servicePortInt, err := strconv.Atoi(servicePort)
	if err != nil {
		return fmt.Errorf("failed to parse port %s: %v", servicePort, err)
	}

	postfix = lo.Ternary(postfix == "", servicePort, postfix)
	healthUrl := fmt.Sprintf("http://%s:%s/health", ip, healthPort)
	err = c.client.Agent().ServiceRegister(&consulapi.AgentServiceRegistration{
		ID:      fmt.Sprintf("%s.%s", serviceName, postfix),
		Name:    serviceName,
		Address: ip,
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

func (c *consulWrapper) GetServices(serviceName string) ([]*consulapi.ServiceEntry, error) {
	services, _, err := c.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		return nil, err
	}

	return services, nil
}

func getLocalIP() (string, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "", err
	}

	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if ok && !ipNet.IP.IsLoopback() && ipNet.IP.To4() != nil {
			return ipNet.IP.String(), nil
		}
	}
	return "", fmt.Errorf("local IP address not found")
}
