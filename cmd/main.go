package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/bigsky-park/go-http-example/api/v1/handler"
)

var (
	httpPort    = 18080
	serviceName = "go-http-example"
)

func main() {
	logger := log.New(os.Stdout, fmt.Sprintf("[%s] ", serviceName), log.LstdFlags)

	hh := handler.NewHello(logger)
	dh := handler.NewDump(logger)

	sm := http.NewServeMux()
	sm.Handle("/hello", hh)
	sm.Handle("/dump", dh)
	sm.Handle("/metrics", promhttp.Handler())

	s := &http.Server{
		Addr:         fmt.Sprintf(":%d", httpPort),
		Handler:      sm,
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			logger.Fatal(err)
		}
	}()

	// TODO: make consul register as configurable
	//consul := client.NewConsulClient(logger, "localhost:8500", "")
	//config := client.ServiceInfo{
	//	Id:      "Hello-01",
	//	Name:    serviceName,
	//	Address: "localhost",
	//	Port:    httpPort,
	//	Tags:    []string{"http", "hello"},
	//}
	//if id, err := consul.Register(&config); err == nil {
	//	defer consul.Deregister(id)
	//} else {
	//	logger.Fatalf("Failed to resiter service %v", err)
	//}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	log.Println("Got signal:", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := s.Shutdown(tc)
	if err != nil {
		logger.Fatal(err)
	}
}
