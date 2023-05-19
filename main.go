package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client
var col *mongo.Collection

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	username := os.Getenv("MONGODB_USERNAME")
	password := os.Getenv("MONGODB_PASSWORD")

	// Construct the connection string with the retrieved credentials
	connectionString := fmt.Sprintf("mongodb+srv://%s:%s@cluster0.wyqvltm.mongodb.net/?retryWrites=true&w=majority", username, password)
	fmt.Println(connectionString)
	// Set up MongoDB connection
	clientOptions := options.Client().ApplyURI(connectionString)

	var mongoErr error
	client, mongoErr = mongo.Connect(context.Background(), clientOptions)
	if mongoErr != nil {
		log.Fatal(mongoErr)
	}

	// Check if MongoDB is running
	pingErr := client.Ping(context.Background(), readpref.Primary())
	if pingErr != nil {
		log.Fatal(pingErr)
	}
	fmt.Println("Connected to MongoDB!")

	// Set the collection
	db := client.Database("test-db")
	col = db.Collection("people")
}

func main() {
	// Create a new router
	r := chi.NewRouter()

	// Define API endpoints
	r.Post("/person", createPerson)
	r.Get("/person/{id}", getPerson)

	// Start the server
	log.Fatal(http.ListenAndServe(":8080", r))
}

type Person struct {
	ID   string `json:"id" bson:"_id"`
	Name string `json:"name" bson:"name"`
}

func createPerson(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var person Person
	err := json.NewDecoder(r.Body).Decode(&person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Generate unique ID
	id := uuid.New().String()

	// Set ID in the Person struct
	person.ID = id

	// Insert person data into MongoDB
	_, insertErr := col.InsertOne(r.Context(), person)
	if insertErr != nil {
		http.Error(w, insertErr.Error(), http.StatusInternalServerError)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)
}

func getPerson(w http.ResponseWriter, r *http.Request) {
	// Get person ID from URL parameter
	id := chi.URLParam(r, "id")

	// Find person data in MongoDB
	var person Person
	err := col.FindOne(r.Context(), bson.M{"_id": id}).Decode(&person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(person)
}
