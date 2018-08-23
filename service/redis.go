package service

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

// ServeHTTP is a http handler
func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status := 200

	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000 // make microseconds
		s.Metrics.RedisDurations.Observe(us)
	}))
	defer timer.ObserveDuration()

	key := strings.Trim(r.URL.Path, "/")
	val, err := s.RedisClient.Get(key).Result()
	if err != nil {
		http.Error(w, "Key not found", http.StatusNotFound)
		status = 404
		s.Metrics.RedisRequests.WithLabelValues("fail").Inc()
	}
	s.Metrics.RedisRequests.WithLabelValues("success").Inc()

	fmt.Fprint(w, val)
	log.Printf("url=\"%s\" remote=\"%s\" key=\"%s\" status=%d\n",
		r.URL, r.RemoteAddr, key, status)
}

// Check returns the redis status
func (s *Service) Check() (bool, error) {
	_, err := s.RedisClient.Ping().Result()
	if err != nil {
		return false, err
	}
	return true, nil
}
