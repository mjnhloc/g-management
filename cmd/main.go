package main

import (
	"log"
	"net/http"
	"os"

	"g-management/internal/services/pkg/container"
	"g-management/internal/services/pkg/mount"
	"g-management/pkg/infrastructure"
	"g-management/pkg/services/elasticsearch/client"

	"github.com/elastic/go-elasticsearch/v9"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, master, err := infrastructure.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	client, err := client.NewClient(elasticsearch.Config{
		Addresses: []string{
			os.Getenv("ES_URL"),
		},
	})
	if err != nil {
		log.Fatalf("Failed to create Elasticsearch client: %v", err)
	}

	repositories := container.NewRepositoryContainers(db, client)
	services := container.NewServiceContainers(client)

	router := infrastructure.NewServer(db)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	err = mount.MountAll(repositories, services, router, db)
	if err != nil {
		log.Fatal("Error happened while mounting all routes: ", "err", err)
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Error happened while starting the server: ", "err", err)
	}

	defer infrastructure.CloseDB(master)
}
