package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/kube-pilot-labs/resource-simulator/internal/service"
	"github.com/kube-pilot-labs/resource-simulator/internal/util"
)

func InitLoadHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	cpuParam := query.Get("cpu")
	memParam := query.Get("mem")

	cpuLoad := 0
	memLoadMB := 0
	var err error

	if cpuParam != "" {
		cpuLoad, err = strconv.Atoi(cpuParam)
		if err != nil || cpuLoad < 0 {
			http.Error(w, "Invalid cpu parameter value.", http.StatusBadRequest)
			return
		}
	}

	if memParam != "" {
		memLoadMB, err = strconv.Atoi(memParam)
		if err != nil || memLoadMB < 0 {
			http.Error(w, "Invalid mem parameter value.", http.StatusBadRequest)
			return
		}
	}

	duration := 60 * time.Second

	service.StartLoad(cpuLoad, memLoadMB, duration)

	response := map[string]interface{}{
		"cpuLoad":   cpuLoad,
		"memLoadMB": memLoadMB,
		"duration":  duration.String(),
	}

	err = util.WriteJSONResponse(w, http.StatusOK, response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func AbortLoadHandler(w http.ResponseWriter, r *http.Request) {
	service.AbortLoad()

	response := map[string]interface{}{
		"result": "true",
	}
	err := util.WriteJSONResponse(w, http.StatusOK, response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
