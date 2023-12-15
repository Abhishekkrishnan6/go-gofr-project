// main_test.go
package main_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "gofr.dev/pkg/gofr"
    "github.com/stretchr/testify/assert"
)

func setupTestApp() *gofr.App {
	app := gofr.New()
	app.POST("/addEntry", addEntryHandler)
	app.GET("/listCars", listCarsHandler)
	app.POST("/updateEntry", updateEntryHandler)
	app.POST("/deleteEntry", deleteEntryHandler)
	app.POST("/register", registerHandler)
	app.POST("/login", loginHandler)
	return app
}

func performRequest(t *testing.T, app *gofr.App, method, path string, data interface{}) *httptest.ResponseRecorder {
	var body []byte
	if data != nil {
		var err error
		body, err = json.Marshal(data)
		assert.Nil(t, err)
	}

	req, err := http.NewRequest(method, path, bytes.NewBuffer(body))
	assert.Nil(t, err)

	rr := httptest.NewRecorder()
	app.ServeHTTP(rr, req)

	return rr
}

func TestRegisterHandler(t *testing.T) {
	app := setupTestApp()

	// Test Case 1: Successful registration
	user := User{
		Username: "testuser",
		Password: "testpassword",
	}

	response := performRequest(t, app, "POST", "/register", user)
	assert.Equal(t, http.StatusOK, response.Code)

	var responseBody map[string]string
	err := json.Unmarshal(response.Body.Bytes(), &responseBody)
	assert.Nil(t, err)
	assert.NotEmpty(t, responseBody["id"])

	// Test Case 2: Attempting to register with an existing username
	response = performRequest(t, app, "POST", "/register", user)
	assert.Equal(t, http.StatusBadRequest, response.Code)
}

// Similarly, write tests for other handler functions.
// Ensure you cover both success and error cases for each handler.
