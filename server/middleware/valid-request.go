package middleware

import (
	"log"
	"net/http"

	"github.com/SirNoob97/server-monitor/server/utils"
)

func DenyGetRequestsBody(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet && r.Body == nil {
			n.ServeHTTP(w, r)
		}
		log.Println("Incomming GET request body is not nil")
		utils.WriteError(w, http.StatusBadRequest)
	})
}
