package routes

import (
	"github.com/ashusatp/todo/controllers"
	"github.com/gorilla/mux"
)

func AuthRouters(router *mux.Router) {
	router.HandleFunc("/login", controllers.Login).Methods("POST")
	router.HandleFunc("/register", controllers.Register).Methods("POST")
}
