package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bigsky-park/go-http-example/api/v1/handler"
)

func main() {
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	hh := handler.NewHello(l)

	sm := http.NewServeMux()
	sm.Handle("/hello", hh)

	s := &http.Server{
		Addr:         ":18080",
		Handler:      sm,
		IdleTimeout:  60 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	log.Println("Got signal:", sig)

	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err := s.Shutdown(tc)
	if err != nil {
		l.Fatal(err)
	}
}
