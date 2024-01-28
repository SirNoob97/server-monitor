package utils

import "net/http"

func WriteError(w http.ResponseWriter, errorCode int) {
	msg := http.StatusText(errorCode)
	http.Error(w, msg, errorCode)
}
