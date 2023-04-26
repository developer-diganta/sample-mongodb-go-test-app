package main_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPersonByID(t *testing.T) {
	// Create a test HTTP server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify that the request URL and method are correct
		assert.Equal(t, "/person/123", r.URL.Path)
		assert.Equal(t, http.MethodGet, r.Method)

		// Write a sample JSON response to the response writer
		person := Person{ID: "123", Name: "John"}
		err := json.NewEncoder(w).Encode(person)
		assert.NoError(t, err)
	}))
	defer ts.Close()

	// Make a request to the test server
	resp, err := http.Get(ts.URL + "/person/123")
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Verify the response status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Decode the response body into a Person struct
	var p Person
	err = json.NewDecoder(resp.Body).Decode(&p)
	assert.NoError(t, err)

	// Verify the Person struct fields
	assert.Equal(t, "123", p.ID)
	assert.Equal(t, "John", p.Name)
}

// Define a sample Person struct
type Person struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
