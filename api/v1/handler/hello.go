package handler

import (
	"fmt"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}

func (h *Hello) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	h.l.Println("Hello is called")

	if r.Method == http.MethodGet {
		fmt.Fprintf(rw, "Hello")
	}
	rw.WriteHeader(http.StatusMethodNotAllowed)
}
