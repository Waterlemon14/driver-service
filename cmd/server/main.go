package main

import (
	"driver-service/internal/adapter"
	"driver-service/internal/handler"
	"driver-service/internal/service"
	"log"
	"net/http"
)

func main() {
	// 1. Initialize Adapters (Infrastructure)
	// In production, swap these lines for Postgres/Redis/Kafka adapters
	repo, err := adapter.NewSQLiteDB("./drivers.db")
	if err != nil {
		log.Fatal("Failed to init SQLite:", err)
	}
	log.Println("Connected to SQLite")

	cache := adapter.NewMemCache()
	queue := adapter.NewChanQueue()

	// 2. Initialize Service with the adapters
	svc := service.NewDriverService(repo, cache, queue)

	// 3. Initialize Handler with the service
	h := handler.NewDriverHandler(svc)

	// 4. Setup Router with the handler
	mux := http.NewServeMux()
	h.RegisterRoutes(mux)

	// 5. Start Server
	log.Println("Driver Service running on :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
