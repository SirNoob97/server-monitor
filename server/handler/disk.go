package handler

import (
	"encoding/json"
	"net/http"

	"github.com/SirNoob97/server-monitor/pkg/disk"
	"github.com/SirNoob97/server-monitor/server/utils"
)

func GetDiskStatus(w http.ResponseWriter, r *http.Request) {
	disk, err := disk.PartitionsInfo()
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(disk)
	if err != nil {
		utils.WriteError(w, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}
