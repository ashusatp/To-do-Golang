package controllers

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/ashusatp/todo/config"
	"github.com/ashusatp/todo/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = config.GetCollection("users")

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func Register(w http.ResponseWriter, r *http.Request) {

	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	hashedPassword, _ := HashPassword(user.Password)
	user.Password = hashedPassword

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	user.ID = primitive.NewObjectID()
	_, err := userCollection.InsertOne(ctx, user)
	if err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	token, _ := models.GenerateJWT(user.Username)
	w.Header().Set("Context-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}

func Login(w http.ResponseWriter, r *http.Request) {

	var user models.User
	json.NewDecoder(r.Body).Decode(&user)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var foundUser models.User

	err := userCollection.FindOne(ctx, bson.M{"username": user.Username}).Decode(&foundUser)

	if err != nil || !CheckPasswordHash(user.Password, foundUser.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token, _ := models.GenerateJWT(foundUser.Username)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
