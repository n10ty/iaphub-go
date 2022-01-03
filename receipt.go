package iaphub

import (
	"encoding/json"
	"fmt"
	"time"
)

// Context of the receipt post
type ReceiptContext string

// Status of the receipt
type ReceiptStatus string

// Webhook status of the receipt
type WebhookStatus string

const (
	ReceiptContextRefresh  ReceiptContext = "refresh"
	ReceiptContextPurchase ReceiptContext = "purchase"
	ReceiptContextRestore  ReceiptContext = "restore"

	ReceiptStatusProcessed  ReceiptStatus = "processed"
	ReceiptStatusProcessing ReceiptStatus = "processing"
	ReceiptStatusDeferred   ReceiptStatus = "deferred"
	ReceiptStatusFailed     ReceiptStatus = "failed"
	ReceiptStatusInvalid    ReceiptStatus = "invalid"
	ReceiptStatusStale      ReceiptStatus = "stale"
	ReceiptStatusSuccess    ReceiptStatus = "success"

	WebhookStatusSuccess WebhookStatus = "success"
	WebhookStatusFailed  WebhookStatus = "failed"
)

type GetReceiptRequest struct {
	ReceiptId string `json:"receiptId"`
}

func (c *Client) GetReceipt(request GetReceiptRequest) (Receipt, error) {
	var receipt Receipt

	if request.ReceiptId == "" {
		return receipt, fmt.Errorf("required parameter \"receiptId\" is missing")
	}

	path := fmt.Sprintf(pathGetReceipt, c.appId, request.ReceiptId)

	response, err := c.requestGet(path, map[string]string{})
	if err != nil {
		return receipt, err
	}

	err = json.Unmarshal(response, &receipt)

	return receipt, err
}

type UpdateReceiptRequest struct {
	UserId        string
	Env           Env
	Platform      Platform
	Token         string
	Sku           string
	Context       ReceiptContext
	ProrationMode ProrationMode
	Upsert        bool
}

func (c *Client) UpdateReceipt(request UpdateReceiptRequest) (ReceiptUpdate, error) {
	var receiptUpdate ReceiptUpdate

	if request.UserId == "" || request.Platform == "" || request.Token == "" || request.Context == "" {
		return receiptUpdate, fmt.Errorf("one of the required parameters(\"userId\", \"platform\", \"token\", \"context\") is missing")
	} else if request.Platform == PlatformAndroid && (request.ProrationMode == "" || request.Sku == "") {
		return receiptUpdate, fmt.Errorf("one of the required parameters(\"sku\", \"prorationMode\") for Android platform are missing")
	}

	if request.Env == "" {
		request.Env = c.env
	}

	path := fmt.Sprintf(pathUpdateReceipt, c.appId, request.UserId)

	body, err := json.Marshal(request)
	if err != nil {
		return receiptUpdate, err
	}

	response, err := c.requestPost(path, map[string]string{}, body)
	if err != nil {
		return receiptUpdate, err
	}

	err = json.Unmarshal(response, &receiptUpdate)

	return receiptUpdate, err
}

type Receipt struct {
	Id           string        `json:"id"`
	CreatedDate  time.Time     `json:"createdDate"`
	ProcessCount int           `json:"processCount"`
	ProcessDate  time.Time     `json:"processedDate"`
	RefreshDate  time.Time     `json:"refreshDate"`
	UserId       string        `json:"user"`
	Platform     Platform      `json:"platform"`
	Status       ReceiptStatus `json:"status"`
	Token        string        `json:"token"`
	Sku          string        `json:"sku"`
}

type ReceiptUpdate struct {
	Status          ReceiptStatus `json:"status"`
	NewTransactions []Transaction `json:"newTransactions"`
	OldTransactions []Transaction `json:"oldTransactions"`
}

type Transaction struct {
	Id                        string                 `json:"id"`
	Sku                       string                 `json:"sku"`
	Purchase                  string                 `json:"purchase"`
	PurchaseDate              time.Time              `json:"purchaseDate"`
	Group                     string                 `json:"group"`
	GroupName                 string                 `json:"groupName"`
	ExpirationDate            time.Time              `json:"expirationDate"`
	AutoResumeDate            time.Time              `json:"autoResumeDate"`
	IsSubscriptionRenewable   bool                   `json:"isSubscriptionRenewable"`
	IsSubscriptionRetryPeriod bool                   `json:"IsSubscriptionRetryPeriod"`
	SubscriptionPeriodType    SubscriptionPeriodType `json:"subscriptionPeriodType"`
}
