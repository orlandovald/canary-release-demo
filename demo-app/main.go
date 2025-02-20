package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

var version string

func main() {
	// Get version from environment variable
	version = os.Getenv("APP_VERSION")
	if version == "" {
		version = "dev"
	}

	handler := http.NewServeMux()
	handler.HandleFunc("/ready", readyHandler)
	handler.HandleFunc("/version", versionHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Starting server on port %s...\n", port)
	http.ListenAndServe(":"+port, handler)
}

func versionHandler(w http.ResponseWriter, r *http.Request) {
	accept := r.Header.Get("Accept")
	if accept == "text/html" {
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, "<html><body><h1>Version: %s</h1></body></html>", version)
		return
	}

	// Default to JSON response
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"version": version}
	json.NewEncoder(w).Encode(response)
}

func readyHandler(w http.ResponseWriter, _ *http.Request) {
	fmt.Println("Ready!")
	fmt.Fprint(w, "Ready!")
}

// const staticDir = "./static"
//
// var emoticons = []string{
// 	"emoticon1.png",
// 	"emoticon2.png",
// 	"emoticon3.png",
// 	"emoticon4.png",
// }
//
// type EmoticonResponse struct {
// 	URL string `json:"url"`
// }
//
// func getEmoticonHandler(w http.ResponseWriter, r *http.Request) {
// 	randomIndex := rand.Intn(len(emoticons))
// 	response := EmoticonResponse{URL: emoticons[randomIndex]}
// 	w.Header().Set("Content-Type", "application/json")
// 	json.NewEncoder(w).Encode(response)
// }
//
// func main() {
// 	http.HandleFunc("/get_emoticon", getEmoticonHandler)
//
// 	// Serve static files
// 	fs := http.FileServer(http.Dir(staticDir))
// 	http.Handle("/static/", http.StripPrefix("/static/", fs))
//
// 	port := ":8080"
// 	println("Server is running on port", port)
// 	if err := http.ListenAndServe(port, nil); err != nil {
// 		println("Failed to start server:", err)
// 	}
// }
