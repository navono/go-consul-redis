package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/navono/go-consul/service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	port := flag.Int("port", 8080, "Port to listen on")
	addrStr := flag.String("addrs", "", "(Required) Redis addrs (may be delimited by ;)")
	ttl := flag.Duration("ttl", time.Second*15, "Service TTL check duration")
	flag.Parse()

	if len(*addrStr) == 0 {
		fmt.Fprintln(os.Stderr, "addrs argument is required")
		flag.PrintDefaults()
		os.Exit(1)
	}

	addrs := strings.Split(*addrStr, ";")

	s, err := service.New(addrs, *ttl, *port)
	if err != nil {
		log.Fatal(err)
	}
	http.Handle("/", prometheus.InstrumentHandler("web", s))
	http.Handle("/metrics", promhttp.Handler())

	l := fmt.Sprintf(":%d", *port)
	log.Print("Listening on ", l)
	log.Fatal(http.ListenAndServe(l, nil))
}
