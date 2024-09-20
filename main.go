package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ashusatp/todo/config"
	"github.com/ashusatp/todo/routes"
	"github.com/gorilla/mux"
)

func main() {
	config.ConnectDB()

	router := mux.NewRouter()

	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"success": true,
			"message": "Hi, I am live",
		})
	}).Methods("GET")

	routes.AuthRouters(router)
	routes.TodoRoutes(router)
	// Register the routes
	// routes
	log.Fatal(http.ListenAndServe(":8080", router))
}
