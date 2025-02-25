package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
)

const (
	staticDir           = "./static"
	responseTemplate    = `<div class='square lined %s' style='background-color: %s'>&nbsp;</div>`
	defaultSquareLength = 30
	defaultCount        = 4
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

	// If the ERROR_RATE environment variable is set, use it
	if rateEnv, err := strconv.Atoi(os.Getenv("ERROR_RATE")); err == nil {
		rate = rateEnv
	}

	t, err := template.ParseFiles("static/index.html")
	if err != nil {
		log.Fatalf("Error parsing template: %v", err)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	handler := http.NewServeMux()
	handler.HandleFunc("/", indexHandler(t))
	handler.HandleFunc("/ready", readyHandler)
	handler.HandleFunc("/version", versionHandler)
	handler.HandleFunc("/hit", hitHandler)

	fmt.Printf("Starting server on port %s...\n", port)
	http.ListenAndServe(":"+port, handler)
}

func indexHandler(t *template.Template) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		count := defaultCount
		squareLength := defaultSquareLength

		if c, err := strconv.Atoi(r.URL.Query().Get("c")); err == nil && c > 0 && c <= 10 {
			count = c
		}
		data := map[string]int{
			"SquareLength": squareLength,
			"Count":        count,
			"Total":        count * count,
		}
		err := t.Execute(w, data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
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
