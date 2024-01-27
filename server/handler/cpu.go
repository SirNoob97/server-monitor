package handler

import (
	"encoding/json"
	"net/http"

	"github.com/SirNoob97/server-monitor/pkg/cpu"
)

type Cpu struct{}

func (h *Cpu) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch {
	case r.Method == http.MethodGet:
		h.getCPUInfo(w, r)
		return
	case r.Method == http.MethodPost:
	case r.Method == http.MethodPut:
	case r.Method == http.MethodDelete:
	case r.Method == http.MethodHead:
	case r.Method == http.MethodPatch:
	case r.Method == http.MethodTrace:
	case r.Method == http.MethodConnect:
	case r.Method == http.MethodOptions:
		methodNotAllowed(w, r)
		return
	default:
		notFoundHandler(w, r)
		return
	}
}

func (h *Cpu) getCPUInfo(w http.ResponseWriter, r *http.Request) {
	cpuInfo, err := cpu.Status()
	if err != nil {
		writeError(w, http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(cpuInfo)
	if err != nil {
		writeError(w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
