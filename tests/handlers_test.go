package test

import (
    "bytes"
    "net/http"
    "net/http/httptest"
    "testing"

	"github.com/random-number-api/pkg"
    "github.com/gin-gonic/gin"
    "github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.POST("/generate", pkg.Generate)

	w := httptest.NewRecorder()
	reqBody := []byte(`{"amount": 100}`)
	req, _ := http.NewRequest("POST", "/generate", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	assert.Contains(t, w.Body.String(), `"id":`)
}

func TestGenerate_InvalidRequest(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.POST("/generate", pkg.Generate)

	w := httptest.NewRecorder()
	reqBody := []byte(`{"amount": -1}`) // Invalid amount
	req, _ := http.NewRequest("POST", "/generate", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	assert.Contains(t, w.Body.String(), `"Invalid request"`)
}

func TestGetResults_NotFound(t *testing.T) {
	gin.SetMode(gin.TestMode)

	router := gin.Default()
	router.GET("/result/:id", pkg.GetResults)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/result/999", nil) 

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	assert.Contains(t, w.Body.String(), `"Result not found"`)
}

func TestGetResults_Success(t *testing.T) {
	gin.SetMode(gin.TestMode)


	pkg.ResultsMap[0] = 1234 // Simulate a result for ID 0

	router := gin.Default()
	router.GET("/result/:id", pkg.GetResults)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/result/0", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	assert.Contains(t, w.Body.String(), `"result":1234`)
}