package service

import "sync"

type RateStore interface {
	Add(laptopId string, score float64) (*Rate, error)
}

type Rate struct {
	Count    uint32
	SumScore float64
}

type InMemoryRateStore struct {
	mutex sync.Mutex
	data  map[string]*Rate
}

func NewInMemoryRateStore() RateStore {
	return &InMemoryRateStore{
		data: make(map[string]*Rate),
	}
}

func (store *InMemoryRateStore) Add(laptopId string, score float64) (*Rate, error) {
	store.mutex.Lock()
	defer store.mutex.Unlock()
	rate := store.data[laptopId]
	if rate == nil {
		rate = &Rate{Count: 1, SumScore: score}
	} else {
		rate.Count += 1
		rate.SumScore += score
	}
	store.data[laptopId] = rate
	return rate, nil
}
