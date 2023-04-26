package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/keploy/go-sdk/integrations/kchi"
	"github.com/keploy/go-sdk/integrations/kmongo"
	"github.com/keploy/go-sdk/keploy"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var client *mongo.Client
var col *kmongo.Collection
func init() {
	// Set up MongoDB connection
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017/test-db?retryWrites=true&w=majority")
	var err error

	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check if MongoDB is running
	err = client.Ping(context.Background(), readpref.Primary())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
}

func main() {
	// Create a new router
		db:= client.Database("test-db")
	col = kmongo.NewCollection(db.Collection("people"));
	r := chi.NewRouter()
	k := keploy.New(keploy.Config{
				App: keploy.AppConfig{
					Name: "my_app",
					Port: "8080",
				},
				Server: keploy.ServerConfig{
					URL: "http://localhost:6789/api",
				},
				})
		r.Use(kchi.ChiMiddlewareV5(k))
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

	_, err = col.InsertOne(r.Context(), person)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
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