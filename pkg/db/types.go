package db

import "time"

type NodeStatus struct {
	PeerID        string        `json:"peerId"`
	IsStaked      bool          `json:"isStaked"`
	TotalUpTime   time.Duration `json:"totalUpTime"`
	FirstLaunched time.Time     `json:"firstLaunched"`
	LastLaunched  time.Time     `json:"lastLaunched"`
}

type SharedData map[string]interface{}
