package service

import (
	"time"

	"github.com/go-redis/redis"
	consul "github.com/hashicorp/consul/api"
)

// Service for service type
type Service struct {
	Name        string
	TTL         time.Duration
	RedisClient redis.UniversalClient
	ConsulAgent *consul.Agent
}

// New returns a new Service instance
func New(addrs []string, ttl time.Duration) (*Service, error) {
	s := new(Service)
	s.Name = "web"
	s.TTL = ttl
	s.RedisClient = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    addrs,
		Password: "1234",
	})

	ok, err := s.Check()
	if !ok {
		return nil, err
	}

	cfg := consul.DefaultConfig()
	cfg.Address = "192.168.56.101:8500"
	c, err := consul.NewClient(cfg)
	if err != nil {
		return nil, err
	}
	s.ConsulAgent = c.Agent()

	serviceDef := &consul.AgentServiceRegistration{
		Name: s.Name,
		Port: 8888,
		Check: &consul.AgentServiceCheck{
			TTL: s.TTL.String(),
		},
	}

	if err := s.ConsulAgent.ServiceRegister(serviceDef); err != nil {
		return nil, err
	}
	go s.UpdateTTL(s.Check)

	return s, nil
}
