package data_types

// LoginEvent represents a user login event.
type LoginEvent struct {
	PeerID      string `json:"peer_id"`
	UserID      string `json:"user_id"`
	Success     bool   `json:"success"`
	Error       string `json:"error,omitempty"`
	ServiceType string `json:"service_type"`
}

// NewLoginEvent creates a new LoginEvent instance.
func NewLoginEvent(peerID, userID, serviceType string, success bool, error string) *LoginEvent {
	return &LoginEvent{
		PeerID:      peerID,
		UserID:      userID,
		Success:     success,
		Error:       error,
		ServiceType: serviceType,
	}
}
