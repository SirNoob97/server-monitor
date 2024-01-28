package handler

import (
	"encoding/json"
	"net/http"

	"github.com/SirNoob97/server-monitor/pkg/memory"
	"github.com/SirNoob97/server-monitor/server/utils"
)

func GetMemoryStatus(w http.ResponseWriter, r *http.Request) {
	memory, err := memory.Status()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(memory)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
