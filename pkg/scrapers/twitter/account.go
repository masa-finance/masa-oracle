package twitter

import (
	"sync"
	"time"
)

type TwitterAccount struct {
	Username         string
	Password         string
	TwoFACode        string
	RateLimitedUntil time.Time
	LastScraped      time.Time
	LoginStatus      string
}

type TwitterAccountManager struct {
	accounts []*TwitterAccount
	index    int
	mutex    sync.Mutex
}

func NewTwitterAccountManager(accounts []*TwitterAccount) *TwitterAccountManager {
	return &TwitterAccountManager{
		accounts: accounts,
		index:    0,
	}
}

func (manager *TwitterAccountManager) GetNextAccount() *TwitterAccount {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	for i := 0; i < len(manager.accounts); i++ {
		account := manager.accounts[manager.index]
		manager.index = (manager.index + 1) % len(manager.accounts)
		if time.Now().After(account.RateLimitedUntil) {
			return account
		}
	}
	return nil
}

func (manager *TwitterAccountManager) MarkAccountRateLimited(account *TwitterAccount) {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	account.RateLimitedUntil = time.Now().Add(GetRateLimitDuration())
}

// AccountState holds the state of a Twitter account
type AccountState struct {
	Username         string
	IsRateLimited    bool
	RateLimitedUntil time.Time
	LastScraped      time.Time
	LoginStatus      string // e.g., "Successful", "Please verify", "Failed - [Reason]"
}

// GetAccountStates returns the state of all Twitter accounts
func (manager *TwitterAccountManager) GetAccountStates() []AccountState {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()

	states := make([]AccountState, len(manager.accounts))
	for i, account := range manager.accounts {
		states[i] = AccountState{
			Username:         account.Username,
			IsRateLimited:    time.Now().Before(account.RateLimitedUntil),
			RateLimitedUntil: account.RateLimitedUntil,
			LastScraped:      account.LastScraped,
			LoginStatus:      account.LoginStatus,
		}
	}
	return states
}
