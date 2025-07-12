package main

import (
	"fmt"
	"log"
	"net/http"

	config "github.com/preetsinghmakkar/crud-in-go/pkg"
)

func main() {
	// Tasks to do
	// 1. Load the configuration
	// 2. Database Setup
	// 3. Setup Router
	// 4. Setup Server

	// Load Configuration
	cfg := config.MustLoad()

	// Setup Router
	router := http.NewServeMux()

	// Let's call the router to handle the Get request.
	router.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})

	// setup server
	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}

	fmt.Printf("Server is running on %s", cfg.HTTPServer.Addr)

	// Start the server
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Failed to start server: %s", err)
	}
}
