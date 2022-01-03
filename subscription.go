package iaphub

import (
	"encoding/json"
	"fmt"
)

type GetSubscriptionRequest struct {
	OriginalPurchaseId string
}

func (c *Client) GetSubscription(request GetSubscriptionRequest) (Subscription, error) {
	var subscription Subscription

	if request.OriginalPurchaseId == "" {
		return subscription, fmt.Errorf("required parameter \"originalPurchaseId\" is missing")
	}

	path := fmt.Sprintf(pathGetSubscription, c.appId, request.OriginalPurchaseId)

	response, err := c.requestGet(path, map[string]string{})
	if err != nil {
		return subscription, err
	}

	err = json.Unmarshal(response, &subscription)

	return subscription, err
}
