package main

import (
	"log"
	"net/http"

	"g-management/pkg/infrastructure"

	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, master, err := infrastructure.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	router := infrastructure.NewServer(db)

	server := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Error happened while starting the server: ", "err", err)
	}

	defer infrastructure.CloseDB(master)
}
