package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

const (
	staticDir        = "./static"
	responseTemplate = `<div class='square lined %s' style='background-color: %s'>&nbsp;</div>`
)

var (
	version   string = "orange"
	errorRate string = "0"
	rate      int
)

type response struct {
	Version   string `json:"version"`
	ErrorRate int    `json:"error_rate"`
}

func main() {
	var err error
	rate, err = strconv.Atoi(errorRate)
	if err != nil {
		log.Fatalf("Error parsing error rate: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	handler := http.NewServeMux()
	handler.HandleFunc("/ready", readyHandler)
	handler.HandleFunc("/version", versionHandler)
	handler.HandleFunc("/hit", hitHandler)
	handler.Handle("/", http.FileServer(http.Dir(staticDir)))

	fmt.Printf("Starting server on port %s...\n", port)
	http.ListenAndServe(":"+port, handler)
}

func hitHandler(w http.ResponseWriter, r *http.Request) {
	lineType := "ok"
	status := http.StatusOK
	if shouldError() {
		lineType = "error"
		status = http.StatusInternalServerError
	}

	html := fmt.Sprintf(responseTemplate, lineType, version)
	w.WriteHeader(status)
	fmt.Fprint(w, html)
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	resp := response{Version: version, ErrorRate: rate}
	json.NewEncoder(w).Encode(resp)
}

// readyHandler is a simple handler for K8s readiness probe
func readyHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Println("Ready!")
	fmt.Fprint(w, "Ready!")
}

func shouldError() bool {
	r := rand.Intn(100) + 1
	return rate >= r
}
