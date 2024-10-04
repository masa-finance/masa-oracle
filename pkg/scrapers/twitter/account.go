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
