package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/kube-pilot-labs/resource-simulator/internal/service"
)

func InitLoadHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	cpuParam := query.Get("cpu")
	memParam := query.Get("mem")

	cpuGoroutines := 0
	memMB := 0
	var err error

	if cpuParam != "" {
		cpuGoroutines, err = strconv.Atoi(cpuParam)
		if err != nil || cpuGoroutines < 0 {
			http.Error(w, "Invalid cpu parameter value.", http.StatusBadRequest)
			return
		}
	}

	if memParam != "" {
		memMB, err = strconv.Atoi(memParam)
		if err != nil || memMB < 0 {
			http.Error(w, "Invalid mem parameter value.", http.StatusBadRequest)
			return
		}
	}

	duration := 60 * time.Second

	// 기존 로직 대신 service의 StartLoad 호출
	service.StartLoad(cpuGoroutines, memMB, duration)

	response := map[string]interface{}{
		"cpuLoad":  cpuGoroutines,
		"memLoad":  memMB,
		"duration": duration.String(),
	}
	jsonData, err := json.Marshal(response)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(jsonData)
}

func AbortLoadHandler(w http.ResponseWriter, r *http.Request) {
	service.AbortLoad()

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("load aborted"))
}
