package handler

import (
	"encoding/json"
	"net/http"

	"github.com/SirNoob97/server-monitor/pkg/cpu"
)

func GetCPUInfo(w http.ResponseWriter, r *http.Request) {
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
