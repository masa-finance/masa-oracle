package pubsub

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"
)

type SafeMap struct {
	mu    sync.RWMutex
	items map[string]*NodeData
}

func NewSafeMap() *SafeMap {
	return &SafeMap{
		items: make(map[string]*NodeData),
	}
}

func (sm *SafeMap) Set(key string, value *NodeData) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.items[key] = value
}

func (sm *SafeMap) Get(key string) (*NodeData, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	value, ok := sm.items[key]
	return value, ok
}

func (sm *SafeMap) Delete(key string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.items, key)
}

func (sm *SafeMap) Len() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.items)
}

func (sm *SafeMap) GetStakedNodesSlice() []NodeData {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	result := make([]NodeData, 0)
	for _, nodeData := range sm.items {
		nd := *nodeData
		nd.CurrentUptime = nodeData.GetCurrentUptime()
		nd.AccumulatedUptime = nodeData.GetAccumulatedUptime()
		nd.CurrentUptimeStr = PrettyDuration(nd.CurrentUptime)
		nd.AccumulatedUptimeStr = PrettyDuration(nd.AccumulatedUptime)
		result = append(result, nd)
	}
	// Sort the slice based on the timestamp
	sort.Slice(result, func(i, j int) bool {
		return result[i].LastUpdated.Before(result[j].LastUpdated)
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

// DumpNodeData writes the entire nodeData map to a file in JSON format.
func (sm *SafeMap) DumpNodeData(filePath string) error {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	data, err := json.Marshal(sm.items)
	if err != nil {
		return fmt.Errorf("could not marshal node data: %w", err)
	}

	err = os.WriteFile(filePath, data, 0644)
	if err != nil {
		return fmt.Errorf("could not write to file: %s, error: %w", filePath, err)
	}

	return nil
}

// LoadNodeData reads nodeData from a file in JSON format and loads it into the map.
func (sm *SafeMap) LoadNodeData(filePath string) error {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("could not read from file: %s, error: %w", filePath, err)
	}

	err = json.Unmarshal(data, &sm.items)
	if err != nil {
		return fmt.Errorf("could not unmarshal JSON data: %w", err)
	}

	return nil
}

func PrettyDuration(d time.Duration) string {
	d = d.Round(time.Minute)
	minute := int64(d / time.Minute)
	h := minute / 60
	minute %= 60
	days := h / 24
	h %= 24

	if days > 0 {
		return fmt.Sprintf("%d days %d hours %d minutes", days, h, minute)
	}
	if h > 0 {
		return fmt.Sprintf("%d hours %d minutes", h, minute)
	}
	return fmt.Sprintf("%d minutes", minute)
}
