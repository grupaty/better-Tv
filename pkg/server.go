package pkg

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gin-gonic/gin"
)

type Server struct {
	httpServer *http.Server
	resultStore *ResultStore // Dependency for handlers
}

func NewServer(resultStore *ResultStore) *Server {
	router := gin.Default()

	numberGenerationHandler := NewNumberProducerHandler(resultStore)
	resultRetrievalHandler := NewResultRetrievalHandler(resultStore)

	router.POST("/generate", numberGenerationHandler.Handle)
	router.GET("/result/:id", resultRetrievalHandler.Handle)

	return &Server{
		httpServer: &http.Server{
			Addr:    ":8080",
			Handler: router,
		},
		resultStore: resultStore,
	}
}

func (s *Server) Start() {
	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
}

func (s *Server) GracefulShutdown() {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := s.httpServer.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}