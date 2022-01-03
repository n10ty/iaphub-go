package iaphub

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
)

// Product type of the purchase
type ProductType string

// Reason of the refund
type RefundReason string

// Current state of the subscription
type SubscriptionState string

// Period type of the subscription
type SubscriptionPeriodType string

// Reason of the renewable subscription cancellation
type SubscriptionCancelReason string

// Proration mode when replacing a subscription (Only required for Android)
type ProrationMode string

const (
	ProductTypeConsumable            ProductType = "consumable"
	ProductTypeNonConsumable         ProductType = "non_consumable"
	ProductTypeRenewableSubscription ProductType = "renewable_subscription"
	ProductTypeSubscription          ProductType = "subscription"

	RefundReasonSubscriptionReplaced RefundReason = "subscription_replaced"
	RefundReasonOther                RefundReason = "other"
	RefundReasonIssue                RefundReason = "issue"
	RefundReasonRemorse              RefundReason = "remorse"
	RefundReasonNotReceived          RefundReason = "not_received"
	RefundReasonDefective            RefundReason = "defective"
	RefundReasonAccidentalPurchase   RefundReason = "accidental_purchase"
	RefundReasonFraud                RefundReason = "fraud"
	RefundReasonFriendlyFraud        RefundReason = "friendly_fraud"
	RefundReasonChargeback           RefundReason = "chargeback"

	SubscriptionStateActive      SubscriptionState = "active"
	SubscriptionStateGracePeriod SubscriptionState = "grace_period"
	SubscriptionStateRetryPeriod SubscriptionState = "retry_period"
	SubscriptionStatePaused      SubscriptionState = "paused"
	SubscriptionStateExpired     SubscriptionState = "expired"

	SubscriptionPeriodTypeNormal SubscriptionPeriodType = "normal"
	SubscriptionPeriodTypeIntro  SubscriptionPeriodType = "intro"
	SubscriptionPeriodTypeTrial  SubscriptionPeriodType = "trial"

	SubscriptionCancelReasonRefunded             SubscriptionCancelReason = "refunded"
	SubscriptionCancelReasonCustomerCancelled    SubscriptionCancelReason = "customer_canceled"
	SubscriptionCancelReasonDeveloperCanceled    SubscriptionCancelReason = "developer_canceled"
	SubscriptionCancelReasonSubscriptionReplaced SubscriptionCancelReason = "subscription_replaced"
	SubscriptionCancelReasonRejectPriceIncrease  SubscriptionCancelReason = "reject_price_increase"
	SubscriptionCancelReasonBillingError         SubscriptionCancelReason = "billing_error"
	SubscriptionCancelReasonProductNotAvailable  SubscriptionCancelReason = "product_not_available"
	SubscriptionCancelReasonUnknown              SubscriptionCancelReason = "unknown"

	ProrationModeImmediateWithTimeProration      ProrationMode = "immediate_with_time_proration"
	ProrationModeImmediateAndChargeProratedPrice ProrationMode = "immediate_and_charge_prorated_price"
	ProrationModeImmediateWithoutProration       ProrationMode = "immediate_without_proration"
)

type GetPurchaseRequest struct {
	PurchaseId string
}

func (c *Client) GetPurchase(request GetPurchaseRequest) (Purchase, error) {
	var purchase Purchase

	if request.PurchaseId == "" {
		return purchase, fmt.Errorf("required parameter \"purchaseId\" is missing")
	}

	path := fmt.Sprintf(pathGetPurchase, c.appId, request.PurchaseId)
	response, err := c.requestGet(path, map[string]string{})

	if err != nil {
		return purchase, err
	}

	err = json.Unmarshal(response, &purchase)

	return purchase, err
}

type GetPurchasesRequest struct {
	Page             int
	Limit            int
	Order            Order
	FromDate         time.Time
	ToDate           time.Time
	User             string
	UserId           string
	OriginalPurchase string
}

func (c *Client) GetPurchases(request GetPurchasesRequest) (PurchaseList, error) {
	var purchaseList PurchaseList

	path := fmt.Sprintf(pathGetPurchases, c.appId)
	params := map[string]string{
		"environment": string(c.env),
	}
	if request.Page != 0 {
		params["page"] = strconv.Itoa(request.Page)
	}
	if request.Limit > 0 && request.Limit <= 100 {
		params["limit"] = strconv.Itoa(request.Limit)
	}
	if request.Order != "" {
		params["order"] = string(request.Order)
	}
	if !request.FromDate.IsZero() {
		params["fromDate"] = request.FromDate.Format(time.RFC3339)
	}
	if !request.ToDate.IsZero() {
		params["toDate"] = request.ToDate.Format(time.RFC3339)
	}
	if request.User != "" {
		params["user"] = request.User
	}
	if request.UserId != "" {
		params["userId"] = request.UserId
	}
	if request.OriginalPurchase != "" {
		params["originalPurchase"] = request.OriginalPurchase
	}

	response, err := c.requestGet(path, params)
	if err != nil {
		return purchaseList, err
	}

	err = json.Unmarshal(response, &purchaseList)

	return purchaseList, err
}

type Subscription = Purchase

type Purchase struct {
	Id                            string                   `json:"id"`
	PurchaseDate                  time.Time                `json:"purchaseDate"`
	Quantity                      int                      `json:"quantity"`
	Platform                      Platform                 `json:"platform"`
	Country                       string                   `json:"country"`
	Tags                          map[string]string        `json:"tags"`
	OrderId                       string                   `json:"orderId"`
	App                           string                   `json:"app"`
	User                          string                   `json:"user"`
	UserId                        string                   `json:"userId"`
	UserIds                       []string                 `json:"userIds"`
	Receipt                       string                   `json:"receipt"`
	AndroidToken                  string                   `json:"androidToken"`
	Product                       string                   `json:"product"`
	ProductSku                    string                   `json:"productSku"`
	ProductType                   ProductType              `json:"productType"`
	ProductGroupName              string                   `json:"productGroupName"`
	Listing                       string                   `json:"listing"`
	Store                         string                   `json:"store"`
	StoreSegmentIndex             int                      `json:"storeSegmentIndex"`
	Currency                      string                   `json:"currency"`
	Price                         float64                  `json:"price"`
	ConvertedCurrency             string                   `json:"convertedCurrency"`
	ConvertedPrice                float64                  `json:"convertedPrice"`
	IsSandbox                     bool                     `json:"isSandbox"`
	IsFamilyShare                 bool                     `json:"isFamilyShare"`
	IsPromo                       bool                     `json:"isPromo"`
	IsRefunded                    bool                     `json:"isRefunded"`
	RefundDate                    time.Time                `json:"refundDate"`
	RefundReason                  RefundReason             `json:"refundReason"`
	RefundAmount                  float64                  `json:"refundAmount"`
	ConvertedRefundAmount         float64                  `json:"convertedRefundAmount"`
	IsSubscription                bool                     `json:"isSubscription"`
	IsSubscriptionActive          bool                     `json:"isSubscriptionActive"`
	IsSubscriptionRenewable       bool                     `json:"isSubscriptionRenewable"`
	IsSubscriptionRetryPeriod     bool                     `json:"isSubscriptionRetryPeriod"`
	IsSubscriptionGracePeriod     bool                     `json:"isSubscriptionGracePeriod"`
	IsTrialConversion             bool                     `json:"isTrialConversion"`
	SubscriptionState             SubscriptionState        `json:"subscriptionState"`
	SubscriptionPeriodType        SubscriptionPeriodType   `json:"SubscriptionPeriodType"`
	SubscriptionCancelReason      SubscriptionCancelReason `json:"subscriptionCancelReason"`
	SubscriptionProrationMode     ProrationMode            `json:"subscriptionProrationMode"`
	SubscriptionRenewalProduct    string                   `json:"subscriptionRenewalProduct"`
	SubscriptionRenewalProductSku string                   `json:"subscriptionRenewalProductSku"`
	ExpirationDate                time.Time                `json:"expirationDate"`
	AutoResumeDate                time.Time                `json:"autoResumeDate"`
	NextPurchase                  string                   `json:"nextPurchase"`
	LinkedPurchase                string                   `json:"linkedPurchase"`
	OriginalPurchase              string                   `json:"originalPurchase"`
}

type PurchaseList struct {
	HasNextPage bool       `json:"hasNextPage"`
	List        []Purchase `json:"list"`
}
