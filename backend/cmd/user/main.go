package user

import (
	"backend/internal/db"
	"backend/internal/handler"
	"backend/internal/middleware"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
)

func HttpPool(pool *pgxpool.Pool) {
	router := http.NewServeMux()

	h := handler.NewHandler(pool)

	router.HandleFunc("POST /api/v1/user/login", h.UserLogin)
	router.HandleFunc("POST /api/v1/user/signup", h.UserSignUp)

	server := http.Server{
		Addr: ":8000",
		Handler: middleware.Logging(router),
	}

	fmt.Println("Server running on http://localhost:8000")
	err := server.ListenAndServe()
	if err != nil {
		log.Fatalf("HTTP Server failed to start with error: %v", err)
	}
}

func main() {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	pool := db.ConnectionPool()
	defer pool.Close()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	go HttpPool(pool)

	<-sigChan
	log.Println("Shutting down gracefully...")

	pool.Close()

	// time buffer to give time for cleanup
	time.Sleep(time.Second)
	log.Println("Shutdown complete")
}