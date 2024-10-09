package pkg

import (
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"github.com/gin-gonic/gin"
)

type NumberProducerHandler struct {
	NumbersInterval int
	Producers       int
	Consumers       int
    requestCounter  int 
	ResultStore     ResultStoreInterface 
}

func NewNumberProducerHandler(store *ResultStore) *NumberProducerHandler {
	numbersInterval, _ := strconv.Atoi(os.Getenv("NUMBERS_INTERVAL"))
	producers, _ := strconv.Atoi(os.Getenv("PRODUCERS"))
	consumers, _ := strconv.Atoi(os.Getenv("CONSUMERS"))

	// Set default values if not configured
	if numbersInterval == 0 {
		numbersInterval = 100
	}
	if producers == 0 {
		producers = 2
	}
	if consumers == 0 {
		consumers = 1
	}

	return &NumberProducerHandler{
		NumbersInterval: numbersInterval,
		Producers:       producers,
		Consumers:       consumers,
		requestCounter: 0,
		ResultStore: store,
	}
}

func (h *NumberProducerHandler) Handle(c *gin.Context) {
	var request struct {
		Amount int `json:"amount"`
	}

	err := c.BindJSON(&request)
	if err != nil || request.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	id := h.requestCounter
	h.requestCounter++

	ch := make(chan int, request.Amount)

	numbersPerProducer := request.Amount / h.Producers

	for i := 0; i < h.Producers; i++ {
		// Add remainder to the last producer
		if i == h.Producers-1 {
			go h.ProduceNumbers(ch, numbersPerProducer+(request.Amount%h.Producers))
		} else {
			go h.ProduceNumbers(ch, numbersPerProducer)
		}
	}

	numbersPerConsumer := request.Amount / h.Consumers
	for i := 0; i < h.Consumers; i++ {
		// Add remainder to the last consumer
		if i == h.Consumers-1 {
			go h.ConsumeNumbers(ch, numbersPerConsumer+(request.Amount%h.Consumers), id)
		} else {
			go h.ConsumeNumbers(ch, numbersPerConsumer, id)
		}
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (h *NumberProducerHandler) ProduceNumbers(ch chan int, amount int) {
	for i := 0; i < amount; i++ {
		randNum := rand.Intn(h.NumbersInterval)
		ch <- randNum
	}
}

func (h *NumberProducerHandler) ConsumeNumbers(ch chan int, amount int, id int) {
	sum := 0
	for i := 0; i < amount; i++ {
		num := <-ch
		sum += num
	}

	h.ResultStore.StoreResult(id,sum)
}