package main

import (
	"fmt"
	"log"
	"net/http"
	"github.com/joho/godotenv"
	"backend/internal/handler"
	"backend/internal/middleware"
	"backend/internal/service"
)

func main() {
	// Load environment variables from .env file (2 directories up from cmd/product/)
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := http.NewServeMux()

	ebayClient, err := service.NewEbayClient()
	if err != nil {
		log.Fatalf("Failed to initialize eBay client: %v", err)
	}

	h := handler.NewSearchHandler(ebayClient)

	router.HandleFunc("GET /api/v1/search/products?q={name...}", h.SearchProductByName)

	server := http.Server{
		Addr: ":8000",
		Handler: middleware.Logging(router),
	}

	fmt.Println("Server running on http://localhost:8000")
	serverErr := server.ListenAndServe()
	if serverErr != nil {
		log.Fatalf("HTTP Server failed to start with error: %v", serverErr)
	}
}