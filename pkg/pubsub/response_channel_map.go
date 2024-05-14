package pubsub

import "sync"

type ResponseChannelMap struct {
	mu    sync.RWMutex
	items map[string]chan []byte
}

var (
	instance *ResponseChannelMap
	once     sync.Once
)

// GetResponseChannelMap returns the singleton instance of ResponseChannelMap.
func GetResponseChannelMap() *ResponseChannelMap {
	once.Do(func() {
		instance = &ResponseChannelMap{
			items: make(map[string]chan []byte),
		}
	})
	return instance
}

// Set associates the specified value with the specified key in the ResponseChannelMap.
// It acquires a write lock to ensure thread-safety while setting the value.
func (drm *ResponseChannelMap) Set(key string, value chan []byte) {
	drm.mu.Lock()
	defer drm.mu.Unlock()
	drm.items[key] = value
}

// Get retrieves the value associated with the specified key from the ResponseChannelMap.
// It acquires a read lock to ensure thread-safety while reading the value.
// If the key exists in the ResponseChannelMap, it returns the corresponding value and true.
// If the key does not exist, it returns nil and false.
func (drm *ResponseChannelMap) Get(key string) (chan []byte, bool) {
	drm.mu.RLock()
	defer drm.mu.RUnlock()
	value, ok := drm.items[key]
	return value, ok
}

// Delete removes the item with the specified key from the ResponseChannelMap.
// It acquires a write lock to ensure thread-safety while deleting the item.
func (drm *ResponseChannelMap) Delete(key string) {
	drm.mu.Lock()
	defer drm.mu.Unlock()
	delete(drm.items, key)
}

// Len returns the number of items in the ResponseChannelMap.
// It acquires a read lock to ensure thread-safety while reading the length.
func (drm *ResponseChannelMap) Len() int {
	drm.mu.RLock()
	defer drm.mu.RUnlock()
	return len(drm.items)
}

func (drm *ResponseChannelMap) CreateChannel(key string) chan []byte {
	ch := make(chan []byte)
	drm.Set(key, ch)
	return ch
}
