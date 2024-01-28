package handler

import (
	"encoding/json"
	"net/http"

	"github.com/SirNoob97/server-monitor/pkg/cpu"
	"github.com/SirNoob97/server-monitor/server/utils"
)

func GetCPUInfo(w http.ResponseWriter, r *http.Request) {
	cpuInfo, err := cpu.Status()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(cpuInfo)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
