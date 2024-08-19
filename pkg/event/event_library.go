package event

import "fmt"

// TODO: Define our events properly. These are placeholders and should be replaced with actual event definitions.

func (a *EventTracker) TrackUserLogin(userID string, client *EventClient) error {
	if a == nil {
		return fmt.Errorf("EventTracker is nil")
	}
	return a.TrackAndSendEvent("user_login", map[string]interface{}{"user_id": userID}, client)
}

func (a *EventTracker) TrackPageView(pageURL string, client *EventClient) error {
	if a == nil {
		return fmt.Errorf("EventTracker is nil")
	}
	return a.TrackAndSendEvent("page_view", map[string]interface{}{"url": pageURL}, client)
}

func (a *EventTracker) TrackPurchase(productID string, amount float64, client *EventClient) error {
	if a == nil {
		return fmt.Errorf("EventTracker is nil")
	}
	return a.TrackAndSendEvent("purchase", map[string]interface{}{
		"product_id": productID,
		"amount":     amount,
	}, client)
}
