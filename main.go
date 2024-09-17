package main

import (
	"log"
	"net/http"

	"github.com/ashusatp/todo/config"
	"github.com/gorilla/mux"
)

func main() {
	config.ConnectDB()

	router := mux.NewRouter()

	// Register the routes
	// routes

	log.Fatal(http.ListenAndServe(":8080", router))
}
