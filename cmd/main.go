package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"g-management/internal/services/pkg/container"
	"g-management/internal/services/pkg/mount"
	"g-management/pkg/infrastructure"
	"g-management/pkg/services"
	"g-management/pkg/services/elasticsearch/client"
	"g-management/pkg/shared/middleware"

	"github.com/elastic/go-elasticsearch/v9"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	// Initialize context with cancellation
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Initialize tracing
	tp, err := infrastructure.InitTracer()
	if err != nil {
		log.Fatalf("Failed to initialize tracer: %v", err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	// Initialize database
	db, master, err := infrastructure.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer infrastructure.CloseDB(master)

	// Initialize Elasticsearch client
	esClient, err := client.NewClient(elasticsearch.Config{
		Addresses: []string{
			os.Getenv("ES_URL"),
		},
	})
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v", err)
	}

	// Initialize Redis client
	redisClient, err := services.NewRedisClient()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	defer redisClient.Client.Close()

	// Initialize WebSocket hub
	hub := services.NewHub(redisClient)
	go hub.Run(ctx)

	// Initialize containers
	repositories := container.NewRepositoryContainers(db, esClient)
	services := container.NewServiceContainers(esClient)

	// Initialize router with tracing middleware
	router := infrastructure.NewServer(db, redisClient)
	router.Use(middleware.TracingMiddleware())

	// Mount WebSocket handler
	router.GET("/ws", hub.HandleWebSocket)

	// Create HTTP server
	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	// Mount all routes
	err = mount.MountAll(repositories, services, router, db)
	if err != nil {
		log.Fatal("Error happened while mounting all routes: ", "err", err)
	}

	// Start server in a goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Give outstanding requests a deadline for completion
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited properly")
}
