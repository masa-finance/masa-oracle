package pubsub

import (
	"encoding/json"
	"fmt"
	"github.com/masa-finance/masa-oracle/pkg/config"
	"sort"
	"sync"
	"time"
)

type SafeMap struct {
	mu    sync.RWMutex
	items map[string]*NodeData
}

// NewSafeMap creates a new instance of SafeMap.
// It initializes the underlying map that will store the key-value pairs.
func NewSafeMap() *SafeMap {
	return &SafeMap{
		items: make(map[string]*NodeData),
	}
}

// Set associates the specified value with the specified key in the SafeMap.
// It acquires a write lock to ensure thread-safety while setting the value.
func (sm *SafeMap) Set(key string, value *NodeData) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.items[key] = value
}

// Get retrieves the value associated with the specified key from the SafeMap.
// It acquires a read lock to ensure thread-safety while reading the value.
// If the key exists in the SafeMap, it returns the corresponding value and true.
// If the key does not exist, it returns nil and false.
func (sm *SafeMap) Get(key string) (*NodeData, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	value, ok := sm.items[key]
	return value, ok
}

// Delete removes the item with the specified key from the SafeMap.
// It acquires a write lock to ensure thread-safety while deleting the item.
func (sm *SafeMap) Delete(key string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.items, key)
}

// Len returns the number of items in the SafeMap.
// It acquires a read lock to ensure thread-safety while reading the length.
func (sm *SafeMap) Len() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.items)
}

// GetStakedNodesSlice returns a slice of NodeData for all staked nodes in the SafeMap.
// It creates a new slice, copies the NodeData from the SafeMap, and populates additional fields
// such as CurrentUptime, AccumulatedUptime, and their string representations.
// The resulting slice is sorted based on the LastUpdated timestamp of each NodeData.
func (sm *SafeMap) GetStakedNodesSlice() []NodeData {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	result := make([]NodeData, 0)
	cfg := config.GetInstance()
	for _, nodeData := range sm.items {
		nd := *nodeData
		nd.CurrentUptime = nodeData.GetCurrentUptime()
		nd.AccumulatedUptime = nodeData.GetAccumulatedUptime()
		nd.CurrentUptimeStr = PrettyDuration(nd.CurrentUptime)
		nd.AccumulatedUptimeStr = PrettyDuration(nd.AccumulatedUptime)
		nd.IsValidator = cfg.Validator
		nd.IsDiscordScraper = cfg.DiscordScraper
		nd.IsTwitterScraper = cfg.TwitterScraper
		nd.IsWebScraper = cfg.WebScraper
		result = append(result, nd)
	}
	// Sort the slice based on the timestamp
	sort.Slice(result, func(i, j int) bool {
		return result[i].LastUpdatedUnix < result[j].LastUpdatedUnix
	})
	return result
}

// MarshalJSON override json MarshalJSON to just return the map
func (sm *SafeMap) MarshalJSON() ([]byte, error) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return json.Marshal(sm.items)
}

// UnmarshalJSON override json UnmarshalJSON to just set the map
func (sm *SafeMap) UnmarshalJSON(b []byte) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	return json.Unmarshal(b, &sm.items)
}

// PrettyDuration takes a time.Duration and returns a string representation
// rounded to the nearest minute. It will include the number of days, hours,
// and minutes as applicable. For example:
//   - 1 day 2 hours 3 minutes
//   - 2 hours 3 minutes
//   - 3 minutes
func PrettyDuration(d time.Duration) string {
	const (
		MinutesInHour = 60
		HoursInDay    = 24
	)

	d = d.Round(time.Minute)
	totalMinutes := int64(d / time.Minute)
	hours := totalMinutes / MinutesInHour
	minutes := totalMinutes % MinutesInHour
	days := hours / HoursInDay
	hours %= HoursInDay

	if days > 0 {
		return fmt.Sprintf("%d days %d hours %d minutes", days, hours, minutes)
	}
	if hours > 0 {
		return fmt.Sprintf("%d hours %d minutes", hours, minutes)
	}
	return fmt.Sprintf("%d minutes", minutes)
}
