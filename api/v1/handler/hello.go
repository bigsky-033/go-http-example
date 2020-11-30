package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	helloCallCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sample_app_hello_call_count",
		Help: "The total number of hello call count",
	})
)

type Hello struct {
	logger *log.Logger
}

func NewHello(logger *log.Logger) *Hello {
	return &Hello{logger}
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.logger.Println("Hello is called")

	if r.Method == http.MethodGet {
		helloCallCounter.Inc()
		fmt.Fprintf(rw, "Hello")
		return
	}
	rw.WriteHeader(http.StatusMethodNotAllowed)
}
