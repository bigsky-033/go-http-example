package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Dump struct {
	logger *log.Logger
}

var (
	dumpCallCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "sample_app_dump_call_count",
		Help: "The total number of dump call count",
	})
)

func NewDump(logger *log.Logger) *Dump {
	return &Dump{logger}
}

func (d *Dump) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	d.logger.Println("Dump is called")
	dumpCallCounter.Inc()
	dump, err := httputil.DumpRequest(r, true)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	d.logger.Println(string(dump))
	fmt.Fprintf(rw, string(dump))
}
