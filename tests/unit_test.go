package test

import (
	"sync"
	"testing"

	"github.com/random-number-api/pkg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockResultStore struct {
	mock.Mock
}

func (m *MockResultStore) StoreResult(id int, sum int) {
	m.Called(id, sum)
}
func (m *MockResultStore) GetResult(id int) (int, bool) {
	m.Called(id)
	return 0, true
}

func TestProduceNumbers(t *testing.T) {
	handler := &pkg.NumberProducerHandler{
		NumbersInterval: 100, 
	}


	ch := make(chan int, 5)

	go handler.ProduceNumbers(ch, 5)

	var producedNums []int
	for i := 0; i < 5; i++ {
		num := <-ch
		producedNums = append(producedNums, num)
	}

	assert.Equal(t, 5, len(producedNums))

	for _, num := range producedNums {
		assert.GreaterOrEqual(t, num, 0)
		assert.Less(t, num, 100)
	}
}

func TestConsumeNumbers(t *testing.T) {
	mockResultStore := new(MockResultStore)
	handler := &pkg.NumberProducerHandler{
		ResultStore: mockResultStore,
	}

	ch := make(chan int, 5)
	expectedSum := 0
	for i := 1; i <= 5; i++ {
		ch <- i
		expectedSum += i
	}

	mockResultStore.On("StoreResult", 1, expectedSum).Return()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		handler.ConsumeNumbers(ch, 5, 1)
	}()

	// Wait for the goroutine to finish
	wg.Wait()

	mockResultStore.AssertCalled(t, "StoreResult", 1, expectedSum)
}

func TestConsumeNumbers_ZeroNumbers(t *testing.T) {
	mockResultStore := new(MockResultStore)
	handler := &pkg.NumberProducerHandler{
		ResultStore: mockResultStore,
	}

	ch := make(chan int, 5)

	mockResultStore.On("StoreResult", 1, 0).Return()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		handler.ConsumeNumbers(ch, 0, 1)
	}()

	wg.Wait()

	mockResultStore.AssertCalled(t, "StoreResult", 1, 0)
}

func TestConsumeNumbers_ClosedChannel(t *testing.T) {
	mockResultStore := new(MockResultStore)
	handler := &pkg.NumberProducerHandler{
		ResultStore: mockResultStore,
	}

	ch := make(chan int, 5)
	expectedSum := 0
	for i := 1; i <= 5; i++ {
		ch <- i
		expectedSum += i
	}
	close(ch) 

	mockResultStore.On("StoreResult", 1, expectedSum).Return()

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		handler.ConsumeNumbers(ch, 5, 1)
	}()

	wg.Wait()

	mockResultStore.AssertCalled(t, "StoreResult", 1, expectedSum)
}

func TestMultipleConcurrentConsumers(t *testing.T) {
	mockResultStore := new(MockResultStore)
	handler := &pkg.NumberProducerHandler{
		ResultStore: mockResultStore,
	}

	ch := make(chan int, 10)
	expectedSum1 := 0
	expectedSum2 := 0
	for i := 1; i <= 10; i++ {
		ch <- i
		if i <= 5 {
			expectedSum1 += i
		} else {
			expectedSum2 += i
		}
	}
	close(ch)

	// Expect StoreResult to be called for each consumer
	mockResultStore.On("StoreResult", 1, expectedSum1).Return()
	mockResultStore.On("StoreResult", 1, expectedSum2).Return()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		handler.ConsumeNumbers(ch, 5, 1)
	}()
	go func() {
		defer wg.Done()
		handler.ConsumeNumbers(ch, 5, 1)
	}()

	wg.Wait()

	mockResultStore.AssertCalled(t, "StoreResult", 1, expectedSum1)
	mockResultStore.AssertCalled(t, "StoreResult", 1, expectedSum2)
}