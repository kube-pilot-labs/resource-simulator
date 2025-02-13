package service

import (
	"context"
	"math"
	"sync"
	"time"
)

var (
	loadMu     sync.Mutex
	loadCancel context.CancelFunc
)

// simulateCPULoad simulates CPU load by running a busy loop for the specified duration.
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

// simulateMemoryLoad allocates memory of size memMB megabytes and holds it for the specified duration.
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

// StartLoad starts CPU and memory load.
func StartLoad(cpuGoroutines int, memMB int, duration time.Duration) {
	loadMu.Lock()
	if loadCancel != nil {
		loadCancel()
		loadCancel = nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	loadCancel = cancel
	loadMu.Unlock()

	if cpuGoroutines > 0 {
		for i := 0; i < cpuGoroutines; i++ {
			go simulateCPULoad(ctx, duration)
		}
	}
	if memMB > 0 {
		go simulateMemoryLoad(ctx, memMB, duration)
	}
}

// AbortLoad cancels the ongoing load.
func AbortLoad() {
	loadMu.Lock()
	if loadCancel != nil {
		loadCancel()
		loadCancel = nil
	}
	loadMu.Unlock()
}
