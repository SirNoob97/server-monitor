package middleware

import (
	"log"
	"net/http"
)

func DenyGetRequestsBody(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.Body == http.NoBody {
			n.ServeHTTP(w, r)
			return
		}
		log.Println("Incomming GET request body is not nil")
		http.Error(w, http.ErrBodyNotAllowed.Error(), http.StatusBadRequest)
	})
}
