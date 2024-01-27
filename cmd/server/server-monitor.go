package main

import (
	"net/http"

	"github.com/SirNoob97/server-monitor/server/handler"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/cpu", &handler.Cpu{})

	http.ListenAndServe(":8080", mux)
}
