package main

import (
	"fmt"
	"net/http"
)

const port = "8080"

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func main() {
	mux := http.NewServeMux()

	// Add readiness endpoint
	mux.HandleFunc("/healthz", healthzHandler)

	// Update fileserver to use /app/ path with prefix stripping
	mux.Handle("/app/", http.StripPrefix("/app", http.FileServer(http.Dir("."))))

	server := &http.Server{
		Addr:    ":" + port,
		Handler: mux,
	}

	fmt.Printf("Server starting on %s\n", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}