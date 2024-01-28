package handler

import (
	"encoding/json"
	"net/http"

	"github.com/SirNoob97/server-monitor/pkg/load"
	"github.com/SirNoob97/server-monitor/server/utils"
)

func GetLoadAVG(w http.ResponseWriter, r *http.Request) {
	avg, err := load.Status()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(avg)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
