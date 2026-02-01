package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"backend/internal/db"
	"backend/internal/handler"
	"backend/pkg"
)

// Handles setting up the routes and starting the http server for the api 
func HttpServer(pool *pgxpool.Pool) {
	router := http.NewServeMux()

	h := handler.NewHandler(pool)

	router.HandleFunc("POST /api/v1/products/add/name", h.AddProductName)
	router.HandleFunc("GET /api/v1/products/get/{id...}", h.GetUserTrackedProducts)
	router.HandleFunc("DELETE /api/v1/products/delete/{id}", h.DeleteProduct)

	server := http.Server{
		Addr: ":8000",
		Handler: pkg.Logging(router),
	}

	fmt.Println("Server running on http://localhost:8000")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("HTTP Server failed to start with error: %v", err)
	}
}

func main() {
	// Load environment variables from .env file (2 directories up from cmd/product/)
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	pool := db.ConnectionPool()
	defer pool.Close() // cleanup if main exits normally

	// go channel for listening to sigint/sigterm signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	// trigger sigChan channel if app recieves either SIGTERM or SIGINT indicating it should shutdown
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// http server started as a goroutine so sigint/sigterm can shut it down
	go HttpServer(pool)

	// Wait for shutdown signal, <- blocks pool.close until notify runs
	<-sigChan
	log.Println("Shutting down gracefully...")

	pool.Close()

	// time buffer to give time for cleanup
	time.Sleep(time.Second)
	log.Println("Shutdown complete")
}