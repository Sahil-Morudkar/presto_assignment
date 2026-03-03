package main

import (
	"log"
	"net/http"

	"github.com/Sahil-Morudkar/presto_assignment/db"
	"github.com/Sahil-Morudkar/presto_assignment/internal/routes"
)

func main() {

	// Initialize the database connection
	database, err := db.NewDB()
	if err != nil {
		log.Fatal(err)
	}

	// Create the router and start the server
	r := routes.NewRouter(database)

	log.Println("Server running on :8080")
	http.ListenAndServe(":8080", r)
}