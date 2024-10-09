package pkg
import (
	"sync"
)
type ResultStoreInterface interface {
	StoreResult(id int, sum int)
	GetResult(id int) (int, bool)
}

type ResultStore struct {
	ResultsMap map[int]int
	mutex      sync.Mutex
}

func NewResultStore() *ResultStore {
	return &ResultStore{
		ResultsMap: make(map[int]int),
	}
}

func (rs *ResultStore) StoreResult(id, sum int) {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	rs.ResultsMap[id] = sum
}

func (rs *ResultStore) GetResult(id int) (int, bool) {
	rs.mutex.Lock()
	defer rs.mutex.Unlock()
	result, exists := rs.ResultsMap[id]
	return result, exists
}