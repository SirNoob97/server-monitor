package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/SirNoob97/server-monitor/server/handler"
	"github.com/SirNoob97/server-monitor/server/middleware"
	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()
	sr := router.PathPrefix("/server-monitor").Methods(http.MethodGet).Subrouter()
	sr.HandleFunc("/cpu", handler.GetCPUInfo)
	sr.Use(middleware.Logging)
	sr.Use(middleware.DenyGetRequestsBody)

	server := http.Server{
		Addr:         "0.0.0.0:8080",
		WriteTimeout: time.Second * 2,
		ReadTimeout:  time.Second * 5,
		IdleTimeout:  time.Second * 10,
		Handler:      sr,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	signalChann := make(chan os.Signal, 1)
	signal.Notify(signalChann, os.Interrupt)
	<-signalChann

	timeout := time.Second * 15
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	server.Shutdown(ctx)
	log.Println("Server shutting down")
	os.Exit(0)
}
