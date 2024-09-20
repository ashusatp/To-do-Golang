package routes

import (
	"github.com/ashusatp/todo/controllers"
	"github.com/ashusatp/todo/middlewares"
	"github.com/gorilla/mux"
)

func TodoRoutes(router *mux.Router) {

	todoRouter := router.PathPrefix("/todos").Subrouter()

	// Apply JWT Middleware
	todoRouter.Use(middlewares.JWTMiddleware)

	todoRouter.HandleFunc("", controllers.CreateTodo).Methods("POST")
	todoRouter.HandleFunc("", controllers.GetTodos).Methods("GET")
	todoRouter.HandleFunc("", controllers.UpdateTodo).Methods("PUT")
	todoRouter.HandleFunc("/status", controllers.UpdateTodoStatus).Methods("PUT")
	todoRouter.HandleFunc("", controllers.DeleteTodo).Methods("DELETE")

}
