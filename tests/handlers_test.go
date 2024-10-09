package test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/random-number-api/pkg"
	"github.com/stretchr/testify/assert"
)
var router *gin.Engine      
var resultStore *pkg.ResultStore

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)

	resultStore = pkg.NewResultStore()
	numberProducerHandler := pkg.NewNumberProducerHandler(resultStore)
	resultRetrievalHandler := pkg.NewResultRetrievalHandler(resultStore)

	router = gin.Default()
	router.POST("/generate", numberProducerHandler.Handle)
	router.GET("/result/:id", resultRetrievalHandler.Handle)

	exitVal := m.Run()
	tearDown()
	os.Exit(exitVal)
}

func tearDown() {
	resultStore = nil 
}


func TestGenerate(t *testing.T) {
	w := httptest.NewRecorder()
	reqBody := []byte(`{"amount": 100}`)
	req, _ := http.NewRequest("POST", "/generate", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	assert.Contains(t, w.Body.String(), `"id":`)
}

func TestGenerate_InvalidRequest(t *testing.T) {
	w := httptest.NewRecorder()
	reqBody := []byte(`{"amount": -1}`)
	req, _ := http.NewRequest("POST", "/generate", bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)

	assert.Contains(t, w.Body.String(), `"Invalid request"`)
}

func TestGetResults_NotFound(t *testing.T) {
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/result/999", nil) 

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)

	assert.Contains(t, w.Body.String(), `"Result not found"`)
}

func TestGetResults_Success(t *testing.T) {
	resultStore.StoreResult(0, 1234) 

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/result/0", nil)

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	assert.Contains(t, w.Body.String(), `"result":1234`)
}