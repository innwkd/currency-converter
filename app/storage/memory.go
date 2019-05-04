package storage

import (
	"sync"
	"time"

	"github.com/yddmat/currency-converter/app/types"
)

type value struct {
	rate types.CurrencyRate
	ttl  time.Duration
}

type MemoryStorage struct {
	storage map[types.CurrencyPair]value
	mu      sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{storage: make(map[types.CurrencyPair]value)}
}

func (re *MemoryStorage) Set(pair types.CurrencyPair, rate types.CurrencyRate, duration time.Duration) (types.CurrencyRate, error) {
	re.mu.Lock()
	defer re.mu.Unlock()

	old, exists := re.storage[pair]
	if exists && !re.expired(old) {
		return old.rate, nil
	}

	re.storage[pair] = value{rate: rate, ttl: duration}
	return rate, nil
}

func (re *MemoryStorage) Get(pair types.CurrencyPair) (types.CurrencyRate, error) {
	re.mu.RLock()
	defer re.mu.RUnlock()

	value, exists := re.storage[pair]
	if !exists || (exists && re.expired(value)) {
		return value.rate, ErrNotExists
	}

	return value.rate, nil
}

func (re *MemoryStorage) GetAll() ([]types.CurrencyRate, error) {
	re.mu.RLock()
	defer re.mu.RUnlock()

	rates := make([]types.CurrencyRate, 0)
	for _, value := range re.storage {
		if !re.expired(value) {
			rates = append(rates, value.rate)
		}
	}

	return rates, nil
}

func (re *MemoryStorage) expired(value value) bool {
	return time.Now().After(value.rate.UpdatedAt.Add(value.ttl))
}
