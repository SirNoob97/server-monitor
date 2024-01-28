package middleware

import (
	"log"
	"net/http"
)

func Logging(n http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    log.Printf("Incomming request from: %s", r.Host)
    n.ServeHTTP(w, r)
  })
}
