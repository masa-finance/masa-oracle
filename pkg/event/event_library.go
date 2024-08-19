package event

func (a *EventTracker) TrackUserLogin(userID string, client *EventClient) error {
	return a.TrackAndSendEvent("user_login", map[string]interface{}{"user_id": userID}, client)
}

func (a *EventTracker) TrackPageView(pageURL string, client *EventClient) error {
	return a.TrackAndSendEvent("page_view", map[string]interface{}{"url": pageURL}, client)
}

func (a *EventTracker) TrackPurchase(productID string, amount float64, client *EventClient) error {
	return a.TrackAndSendEvent("purchase", map[string]interface{}{
		"product_id": productID,
		"amount":     amount,
	}, client)
}
