package app

import (
	"fmt"
	"sync"
	"time"
)

// Item struct to store value and lifetime.
type Item struct {
	Value    interface{}
	Lifetime time.Time
}

// Stream struct to hold entries.
type Stream struct {
	Entries map[string]interface{}
}

// MemoryStorage struct to handle storage of items and streams.
type MemoryStorage struct {
	storage map[string]Item
	mutex   sync.Mutex
}

// NewMemoryStorage creates a new instance of MemoryStorage.
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		storage: make(map[string]Item),
	}
}

// Save stores a value with an optional lifetime.
func (ms *MemoryStorage) Save(key string, value interface{}, lifetime *time.Time) {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	item := Item{Value: value}
	if lifetime != nil {
		item.Lifetime = *lifetime
	}
	ms.storage[key] = item
}

// Get retrieves a value by key.
func (ms *MemoryStorage) Get(key string) interface{} {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	item, exists := ms.storage[key]
	fmt.Print("INFO || MEMORY || GET key=", key, " value=", item.Value, "\n")
	if !exists {
		return nil
	}

	return item.Value
}

// GetType retrieves the type of the value stored at key.
func (ms *MemoryStorage) GetType(key string) string {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	item, exists := ms.storage[key]
	if !exists || (!item.Lifetime.IsZero() && ms.Expired(item.Lifetime)) {
		return "none"
	}

	switch item.Value.(type) {
	case string:
		return "string"
	case int:
		return "integer"
	case float64:
		return "float"
	default:
		return "unknown"
	}
}

// Delete removes an item by key.
func (ms *MemoryStorage) Delete(key string) int {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	if _, exists := ms.storage[key]; exists {
		delete(ms.storage, key)
		return 1
	}
	return 0
}

// Expired checks if the item is expired.
func (ms *MemoryStorage) Expired(lifetime time.Time) bool {
	return time.Now().After(lifetime)
}

// PrintAll prints all items in storage.
func (ms *MemoryStorage) PrintAll() {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	if len(ms.storage) == 0 {
		fmt.Println("Memory is empty")
		return
	}
	for key, item := range ms.storage {
		fmt.Printf("%s: %v\n", key, item.Value)
	}
}

// Exists checks if a key exists in storage.
func (ms *MemoryStorage) Exists(key string) bool {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	_, exists := ms.storage[key]
	return exists
}

// Keys retrieves all keys in storage.
func (ms *MemoryStorage) Keys() []string {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	keys := make([]string, 0, len(ms.storage))
	for key := range ms.storage {
		keys = append(keys, key)
	}
	return keys
}

// Flush clears all items in storage.
func (ms *MemoryStorage) Flush() {
	ms.mutex.Lock()
	defer ms.mutex.Unlock()
	ms.storage = make(map[string]Item)
}
