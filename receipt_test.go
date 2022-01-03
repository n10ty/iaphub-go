package iaphub_test

import (
	"bytes"
	"fmt"
	"github.com/n10ty/iaphub-go"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"
)

var (
	receiptId = "receipt-1"
	token     = "token-1"
	sku       = "sku-1"
)

func TestClient_GetReceipt(t *testing.T) {
	httpClient := newClient(
		func(req *http.Request) (*http.Response, error) {
			expectedUrl := fmt.Sprintf("https://api.iaphub.com/v1/app/%s/receipt/%s?environment=sandbox", appId1, receiptId)
			if req.URL.String() != expectedUrl {
				return nil, fmt.Errorf("wrong URL; expected: %s, got: %s", expectedUrl, req.URL.String())
			}

			authHeaderVal := "ApiKey " + apiKey1
			if req.Header.Get("Authorization") != authHeaderVal {
				return nil, fmt.Errorf("wrong auth header; expected: %s, got: %s", authHeaderVal, req.Header.Get("Authorization"))
			}

			body := `{"id":"receipt-1","createdDate":"2019-10-12T17:34:33.256Z","processCount":1,"processedDate":"2019-10-12T17:34:34.256Z","refreshDate":"2019-10-13T17:34:34.256Z","user":"user-id-1","receipt":"receipt-1","platform":"android","status":"processed","token":"token-1","sku":"subscription_1"}`
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBuffer([]byte(body))),
			}, nil
		},
	)

	client, _ := iaphub.NewClient(
		apiKey1,
		appId1,
		iaphub.UseClient(httpClient),
		iaphub.UseEnv(iaphub.Env(env)),
	)

	getReceiptRequest := iaphub.GetReceiptRequest{
		ReceiptId: receiptId,
	}

	receipt, err := client.GetReceipt(getReceiptRequest)

	if err != nil {
		t.Errorf("GetReceipt failed: %s", err)
	}

	if !reflect.DeepEqual(receipt, dummyReceipt()) {
		t.Errorf("wrong receipt; expected:\n%#v\ngot:\n%#v\n", dummyReceipt(), receipt)
	}
}

func TestClient_UpdateReceipt(t *testing.T) {
	httpClient := newClient(
		func(req *http.Request) (*http.Response, error) {
			expectedUrl := fmt.Sprintf("https://api.iaphub.com/v1/app/%s/user/%s/receipt?environment=sandbox", appId1, userId1)
			if req.URL.String() != expectedUrl {
				return nil, fmt.Errorf("wrong URL; expected: %s, got: %s", expectedUrl, req.URL.String())
			}

			authHeaderVal := "ApiKey " + apiKey1
			if req.Header.Get("Authorization") != authHeaderVal {
				return nil, fmt.Errorf("wrong auth header; expected: %s, got: %s", authHeaderVal, req.Header.Get("Authorization"))
			}

			body := `{"status":"success","newTransactions":[{"id":"5e517bdd0613c16f11e7fae0","sku":"pack30_tier20","purchase":"2e517bdd0613c16f11e7faz2","purchaseDate":"2019-10-12T17:34:33.256Z","group": "3e517bdd0613c16f41e7fae2","groupName":"pack","webhookStatus":"success"}],"oldTransactions":[]}`
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBuffer([]byte(body))),
			}, nil
		},
	)

	client, _ := iaphub.NewClient(
		apiKey1,
		appId1,
		iaphub.UseClient(httpClient),
		iaphub.UseEnv(iaphub.Env(env)),
	)

	updateReceiptRequest := iaphub.UpdateReceiptRequest{
		UserId:        userId1,
		Env:           iaphub.Env(env),
		Platform:      iaphub.PlatformAndroid,
		Token:         token,
		Sku:           sku,
		Context:       iaphub.ReceiptContextRefresh,
		ProrationMode: iaphub.ProrationModeImmediateAndChargeProratedPrice,
		Upsert:        true,
	}

	actualReceiptUpdate, err := client.UpdateReceipt(updateReceiptRequest)

	loc, _ := time.LoadLocation("UTC")
	purchaseDate := time.Unix(1570901673, 256*1000*1000).In(loc)
	expectedReceiptUpdate := iaphub.ReceiptUpdate{
		Status: iaphub.ReceiptStatusSuccess,
		NewTransactions: []iaphub.Transaction{
			{
				Id:                        "5e517bdd0613c16f11e7fae0",
				Sku:                       "pack30_tier20",
				Purchase:                  "2e517bdd0613c16f11e7faz2",
				PurchaseDate:              purchaseDate,
				Group:                     "3e517bdd0613c16f41e7fae2",
				GroupName:                 "pack",
				ExpirationDate:            time.Time{},
				AutoResumeDate:            time.Time{},
				IsSubscriptionRenewable:   false,
				IsSubscriptionRetryPeriod: false,
				SubscriptionPeriodType:    "",
			},
		},
		OldTransactions: []iaphub.Transaction{},
	}
	if err != nil {
		t.Errorf("UpdateReceipt failed: %s", err)
		return
	}

	if !reflect.DeepEqual(expectedReceiptUpdate, actualReceiptUpdate) {
		t.Errorf("wrong receipt update; expected:\n%#v\ngot:\n%#v\n", expectedReceiptUpdate, actualReceiptUpdate)
	}
}

func dummyReceipt() iaphub.Receipt {
	loc, _ := time.LoadLocation("UTC")
	date := time.Unix(1570901673, 256*1000*1000).In(loc)
	return iaphub.Receipt{
		Id:           receiptId,
		CreatedDate:  date,
		ProcessCount: 1,
		ProcessDate:  date.Add(time.Second),
		RefreshDate:  date.Add(time.Second).Add(24 * time.Hour),
		UserId:       userId1,
		Platform:     iaphub.PlatformAndroid,
		Status:       "processed",
		Token:        token,
		Sku:          "subscription_1",
	}
}
