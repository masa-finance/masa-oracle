package api

import (
	"encoding/json"
	"time"
)

type AccountState struct {
	Username         string
	IsRateLimited    bool
	RateLimitedUntil time.Time
	LastScraped      time.Time
	LoginStatus      string // e.g., "Successful", "Please verify", "Failed - [Reason]"
}

func GetAccountStates(value string) ([]AccountState, error) {
	var accountStates []AccountState
	err := json.Unmarshal([]byte(value), &accountStates)
	if err != nil {
		return nil, err
	}
	return accountStates, nil
}
