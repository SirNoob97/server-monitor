package handler

import "net/http"

func methodNotAllowed(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusMethodNotAllowed)
}

func notFoundHandler(w http.ResponseWriter, r *http.Request) {
	writeError(w, http.StatusNotFound)
}

func writeError(w http.ResponseWriter, errorCode int) {
	msg := http.StatusText(errorCode)
	http.Error(w, msg, errorCode)
}
