package service

import (
	"time"

	"github.com/go-redis/redis"
	consul "github.com/hashicorp/consul/api"
	"github.com/prometheus/client_golang/prometheus"
)

// Service for service type
type Service struct {
	Name        string
	Port        int
	TTL         time.Duration
	RedisClient redis.UniversalClient
	ConsulAgent *consul.Agent
	Metrics     Metrics
}

// Metrics for prometheus
type Metrics struct {
	RedisRequests  *prometheus.CounterVec
	RedisDurations prometheus.Summary
}

// New returns a new Service instance
func New(addrs []string, ttl time.Duration, port int) (*Service, error) {
	s := new(Service)
	s.Name = "web"
	s.Port = port
	s.TTL = ttl
	s.RedisClient = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs:    addrs,
		Password: "1234",
	})

	ok, err := s.Check()
	if !ok {
		return nil, err
	}

	s.metricsRegister()
	s.consulRegister()

	go s.UpdateTTL(s.Check)

	return s, nil
}

func (s *Service) metricsRegister() {
	s.Metrics.RedisRequests = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "redis_requests_total",
			Help: "How many Redis requests processed, partitioned by status",
		},
		[]string{"status"},
	)
	prometheus.MustRegister(s.Metrics.RedisRequests)

	s.Metrics.RedisDurations = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name:       "redis_request_durations",
			Help:       "Redis requests latencies in microseconds",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
		},
	)
	prometheus.MustRegister(s.Metrics.RedisDurations)
}

func (s *Service) consulRegister() {
	cfg := consul.DefaultConfig()
	cfg.Address = "192.168.56.101:8500"
	c, err := consul.NewClient(cfg)
	if err != nil {
		panic("Failed to connect to Consul agent")
	}
	s.ConsulAgent = c.Agent()

	serviceDef := &consul.AgentServiceRegistration{
		Name: s.Name,
		Port: s.Port,
		Check: &consul.AgentServiceCheck{
			TTL: s.TTL.String(),
		},
	}

	if err := s.ConsulAgent.ServiceRegister(serviceDef); err != nil {
		panic("Failed to register Service")
	}
}
