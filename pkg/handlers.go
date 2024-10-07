package pkg

import (
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/gin-gonic/gin"
)

var ResultsMap = make(map[int]int)
var mutex sync.Mutex

var requestCounter int 
var numbersInterval int
var producers int
var consumers int

func init() {
    // Initialize configuration from environment variables
    numbersInterval, _ = strconv.Atoi(os.Getenv("NUMBERS_INTERVAL"))
    producers, _ = strconv.Atoi(os.Getenv("PRODUCERS"))
    consumers, _ = strconv.Atoi(os.Getenv("CONSUMERS"))

    // Default values in case the environment variables are not set
    if numbersInterval == 0  {
        numbersInterval = 100
    }
    if producers == 0 {
        producers = 2
    }
    if consumers == 0 {
        consumers = 1
    }
	requestCounter = 0
}

func Generate(c *gin.Context) {
	var request struct {
		Amount int `json:"amount"`
	}

	err := c.BindJSON(&request)
	if err != nil || request.Amount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	id := requestCounter
	requestCounter++

	ch := make(chan int, request.Amount)

	numbersPerProducer := request.Amount / producers

	for i := 0; i < producers; i++ {
		// Add remainder to the last producer
		if i == producers-1 {
			go ProduceNumbers(ch, numbersPerProducer+(request.Amount%producers))
		} else {
			go ProduceNumbers(ch, numbersPerProducer)
		}
	}

	numbersPerConsumer := request.Amount / consumers
	for i := 0; i < consumers; i++ {
		// Add remainder to the last consumer
		if i == consumers-1 {
			go ConsumeNumbers(ch, numbersPerConsumer+(request.Amount%consumers), id)
		} else {
			go ConsumeNumbers(ch, numbersPerConsumer, id)
		}
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func ConsumeNumbers(ch chan int, amount int, id int) {
	sum := 0
	for i := 0; i < amount; i++ {
		num := <-ch
		sum += num
	}

	mutex.Lock()
	ResultsMap[id] = sum
	mutex.Unlock()
}

func ProduceNumbers(ch chan int, amount int) []int {
	numbers := make([]int, amount)
	for i := 0; i < amount; i++ {
		randNum := rand.Intn(numbersInterval) // Generate random number wasn't specified
		ch <- randNum
	}

	return numbers
}

func GetResults(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	mutex.Lock()
	result, exists := ResultsMap[id]
	mutex.Unlock()

	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Result not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": result})
}
