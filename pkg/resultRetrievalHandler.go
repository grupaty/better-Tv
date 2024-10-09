package pkg

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ResultRetrievalHandler struct{
	resultStore ResultStoreInterface
}

func NewResultRetrievalHandler(store *ResultStore) *ResultRetrievalHandler {
	return &ResultRetrievalHandler{resultStore: store}
}

func (h *ResultRetrievalHandler) Handle(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	ooutput, exists := h.resultStore.GetResult(id)

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Result not found"})
		return
	}

	// Return the result
	c.JSON(http.StatusOK, gin.H{"result": ooutput})
}