package main

import (
	"context"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"
)

var (
	loadMu            sync.Mutex
	currentLoadCancel context.CancelFunc
)

// simulateCPULoad performs a busy loop for the specified duration.
func simulateCPULoad(ctx context.Context, duration time.Duration) {
	end := time.Now().Add(duration)
	for time.Now().Before(end) {
		select {
		case <-ctx.Done():
			return
		default:
			_ = math.Sqrt(12345.6789)
		}
	}
}

// simulateMemoryLoad allocates memory of size memMB megabytes and waits for the specified duration.
func simulateMemoryLoad(ctx context.Context, memMB int, duration time.Duration) {
	size := memMB * 1024 * 1024
	data := make([]byte, size)
	for i := range data {
		data[i] = 1
	}
	select {
	case <-ctx.Done():
		return
	case <-time.After(duration):
		return
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("pong"))
}

func initLoadHandler(w http.ResponseWriter, r *http.Request) {
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

	duration := 10 * time.Second

	loadMu.Lock()
	if currentLoadCancel != nil {
		currentLoadCancel()
		currentLoadCancel = nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	currentLoadCancel = cancel
	loadMu.Unlock()

	if cpuGoroutines > 0 {
		for i := 0; i < cpuGoroutines; i++ {
			go simulateCPULoad(ctx, duration)
		}
	}

	if memMB > 0 {
		go simulateMemoryLoad(ctx, memMB, duration)
	}

	responseMsg := fmt.Sprintf("CPU load: %d goroutines, memory load: %d MB, duration: %v", cpuGoroutines, memMB, duration)
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(responseMsg))
}

func abortLoadHandler(w http.ResponseWriter, r *http.Request) {
	loadMu.Lock()
	if currentLoadCancel != nil {
		currentLoadCancel()
		currentLoadCancel = nil
	}
	loadMu.Unlock()

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("load aborted"))
}

func main() {
	http.HandleFunc("/ping", pingHandler)
	http.HandleFunc("/init", initLoadHandler)
	http.HandleFunc("/abort", abortLoadHandler)

	port := "8080"
	log.Printf("Server started: http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server execution error: %v", err)
	}
}
