package main

import (
	"log"
	"net/http"

	"g-management/internal/services/pkg/container"
	"g-management/internal/services/pkg/mount"
	"g-management/pkg/infrastructure"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, master, err := infrastructure.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	repositories := container.NewRepositoryContainers(db)

	router := infrastructure.NewServer(db)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	err = mount.MountAll(repositories, router, db)
	if err != nil {
		log.Fatal("Error happened while mounting all routes: ", "err", err)
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Error happened while starting the server: ", "err", err)
	}

	defer infrastructure.CloseDB(master)
	// es, err := elasticsearch.NewDefaultClient()
	// if err != nil {
	// 	log.Fatalf("Error creating ES client: %v", err)
	// }

	// _, err = es.Indices.Create("classes")
	// if err != nil {
	// 	log.Fatalf("Cannot create index: %s", err)
	// }

	// fmt.Println(es.Info())
}
