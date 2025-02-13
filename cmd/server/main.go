package main

import (
	"log"
	"net/http"

	"github.com/kube-pilot-labs/resource-simulator/internal/handler"
)

func main() {
	http.HandleFunc("/ping", handler.PingHandler)
	http.HandleFunc("/init", handler.InitLoadHandler)
	http.HandleFunc("/abort", handler.AbortLoadHandler)

	port := "8080"
	log.Printf("Server started: http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Server execution error: %v", err)
	}
}
