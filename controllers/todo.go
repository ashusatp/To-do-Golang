package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/ashusatp/todo/config"
	"github.com/ashusatp/todo/middlewares"
	"github.com/ashusatp/todo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var todoCollection *mongo.Collection = config.GetCollection("todos")

func CreateTodo(w http.ResponseWriter, r *http.Request) {

	var todo models.Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if todo.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	todo.ID = primitive.NewObjectID()
	userID, err := getUserFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	todo.UserID = userID
	todo.Done = false

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = todoCollection.InsertOne(ctx, todo)
	if err != nil {
		http.Error(w, "Could not create todo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Context-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    todo,
	})
}

func GetTodos(w http.ResponseWriter, r *http.Request) {

	var todos []models.Todo

	page, limit := getPaginationParams(r)
	skip := (page - 1) * limit

	userID, err := getUserFromContext(r)
	if err != nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	findOptions := options.Find().SetSkip(int64(skip)).SetLimit(int64(limit))

	curr, err := todoCollection.Find(ctx, filter, findOptions)
	if err != nil {
		http.Error(w, "Could not retrieve todos", http.StatusInternalServerError)
		return
	}

	// Decode each document in the cursor into the todos slice
	for curr.Next(ctx) {
		var todo models.Todo
		if err := curr.Decode(&todo); err != nil {
			log.Printf("Error decoding todo: %v", err)
			continue
		}
		todos = append(todos, todo)
	}

	if err := curr.Err(); err != nil {
		http.Error(w, "Error retrieving todos", http.StatusInternalServerError)
		return
	}

	curr.Close(ctx)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"data":    todos,
		"page":    page,
		"limit":   limit,
	})
}

func UpdateTodo(w http.ResponseWriter, r *http.Request) {
	var todoUpdate models.Todo

	err := json.NewDecoder(r.Body).Decode(&todoUpdate)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if todoUpdate.Title == "" {
		http.Error(w, "Title is required", http.StatusBadRequest)
		return
	}

	todoID := r.URL.Query().Get("id")
	objID, err := primitive.ObjectIDFromHex(todoID)
	if err != nil {
		log.Printf("Invalid todo ID: %v", err)
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	userID, err := getUserFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"title": todoUpdate.Title,
			"done":  todoUpdate.Done,
		},
	}

	filter := bson.M{"_id": objID, "user_id": userID}

	result, err := todoCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		log.Printf("Error updating todo: %v", err)
		http.Error(w, "Could not update todo", http.StatusInternalServerError)
		return
	}

	if result.MatchedCount == 0 {
		http.Error(w, "Todo not found or unauthorized", http.StatusNotFound)
		return
	}

	w.Header().Set("Context-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Todo updated successfully",
	})
}

func DeleteTodo(w http.ResponseWriter, r *http.Request) {
	todoID := r.URL.Query().Get("id")
	objID, err := primitive.ObjectIDFromHex(todoID)
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	userID, err := getUserFromContext(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Set a 10-second timeout for the DB operation
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Filter for the todo with the correct ID and userID (for security reasons)
	filter := bson.M{"_id": objID, "user_id": userID}

	// Delete the todo from the database
	result, err := todoCollection.DeleteOne(ctx, filter)
	if err != nil {
		log.Printf("Error deleting todo: %v", err)
		http.Error(w, "Could not delete todo", http.StatusInternalServerError)
		return
	}
	if result.DeletedCount == 0 {
		http.Error(w, "Todo not found or unauthorized", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success": true,
		"message": "Todo deleted successfully",
	})
}

func getUserFromContext(r *http.Request) (string, error) {
	userID := r.Context().Value(middlewares.UserContextKey)
	if userID == nil {
		return "", errors.New("user name not found in context")
	}
	return userID.(string), nil
}

func getPaginationParams(r *http.Request) (int, int) {
	// Default pagination values
	page := 1
	limit := 10

	if p, err := strconv.Atoi(r.URL.Query().Get("page")); err == nil && p > 0 {
		page = p
	}

	if l, err := strconv.Atoi(r.URL.Query().Get("limit")); err == nil && l > 0 {
		limit = l
	}

	return page, limit
}
