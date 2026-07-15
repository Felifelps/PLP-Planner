package main

import (
	"context"
	"log"
	"net/http"

	"plp-planner/bootstrap"
	"plp-planner/database"
)

func main() {
	ctx := context.Background()

	db, err := database.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if err := database.RunMigrations(ctx); err != nil {
		log.Fatal(err)
	}

	dependencies := bootstrap.InitializeDependencies(db)
	router := bootstrap.InitializeRouter(dependencies)

	log.Println("Servidor rodando em http://localhost:8080")

	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}