package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/preetsinghmakkar/crud-in-go/internal/config"
	"github.com/preetsinghmakkar/crud-in-go/internal/http/handlers/student"
	"github.com/preetsinghmakkar/crud-in-go/internal/storage/sqlite"
)

func main() {
	// Tasks to do
	// 1. Load the configuration
	// 2. Database Setup
	// 3. Setup Router
	// 4. Setup Server

	// load config
	cfg := config.MustLoad()
	// database setup

	storage, err := sqlite.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("storage initialized", slog.String("env", cfg.Env), slog.String("version", "1.0.0"))

	// Setup Router
	router := http.NewServeMux()

	router.HandleFunc("POST /api/students", student.New(storage))
	router.HandleFunc("GET /api/students/{id}", student.GetById(storage))
	router.HandleFunc("GET /api/students", student.GetList(storage))

	// setup server
	server := http.Server{
		Addr:    cfg.HTTPServer.Addr,
		Handler: router,
	}

	slog.Info("server started", slog.String("address", cfg.HTTPServer.Addr))

	done := make(chan os.Signal, 1)

	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM) // This will notify the channel when an interrupt signal is received

	go func() { // Starting a goroutine(Thread) to handle server start
		// Start the server
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("Failed to start server: %s", err)
		}
	}()

	<-done // furthur execution will be blocked until a signal is received. It means logic written after this line will be executed only when a signal is received.

	slog.Info("shutting down the server")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel() // defer executes this cancel at the end of the function.
	// Flow of defer :
	// ctx is created with 5s timeout
	// ↓
	// server.Shutdown(ctx) starts (it blocks)
	// ↓
	// If it finishes earlier → great
	// If not, after 5s → context expires, and Shutdown exits
	// ↓
	// After Shutdown finishes → defer cancel() runs

	if err := server.Shutdown(ctx); err != nil {
		slog.Error("failed to shutdown server", slog.String("error", err.Error()))
	}

	slog.Info("server shutdown successfully")

}
