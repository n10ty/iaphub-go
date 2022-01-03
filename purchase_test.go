package iaphub_test

import (
	"bytes"
	"fmt"
	"github.com/rentaapp/iaphub-go"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
	"time"
)

var (
	purchaseId     = "purchase-1"
	internalUser   = "user-internal-1"
	origPurchaseId = "orig-purchase-1"
)

func TestClient_GetPurchase(t *testing.T) {
	httpClient := newClient(
		func(req *http.Request) (*http.Response, error) {
			expectedUrl := fmt.Sprintf("https://api.iaphub.com/v1/app/%s/purchase/%s?environment=sandbox", appId1, purchaseId)
			if req.URL.String() != expectedUrl {
				return nil, fmt.Errorf("wrong URL; expected: %s, got: %s", expectedUrl, req.URL.String())
			}

			body := `{"id":"purchase-1","purchaseDate":"2019-10-12T17:34:33.256Z","quantity":1,"platform":"ios","country":"US","tags":{},"orderId":"9873637705964380","app":"5d86507259e828b8fe321f7e","purchase":"5d865c10c41280ba7f0ce9c2","userId":"62785074-8f32-42a5-b86b-90dbd79ce212","product":"5d86507259e828b8fe321f8a","listing":"5d86507259e828b8fe321f32","store":"5d86507259e828b8fe321f85","receipt":"5d86507259e828b8fe321f34","currency":"USD","price":19.99,"convertedCurrency":"USD","convertedPrice":19.99,"isSandbox":false,"isRefunded":false,"isSubscription":true,"isSubscriptionActive":true,"isSubscriptionRenewable":true,"isSubscriptionRetryPeriod":false,"isTrialConversion":false,"subscriptionPeriodType":"normal","expirationDate":"2019-11-12T17:34:33.256Z","linkedPurchase":"2d865c10c41280ba7f0ce9c4","originalPurchase":"2d865c10c41280ba7f0ce9c4","productSku":"membership_pricing1","productType":"renewable_subscription"}`

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

	getPurchaseRequest := iaphub.GetPurchaseRequest{
		PurchaseId: purchaseId,
	}

	purchase, err := client.GetPurchase(getPurchaseRequest)

	if err != nil {
		t.Errorf("GetPurchase failed: %s", err)
	}

	if !reflect.DeepEqual(purchase, dummyPurchase()) {
		t.Errorf("wrong purchase; expected:\n%#v\ngot:\n%#v\n", dummyPurchase(), purchase)
	}
}

func TestClient_GetPurchases(t *testing.T) {
	httpClient := newClient(
		func(req *http.Request) (*http.Response, error) {
			expectedUrl := "https://api.iaphub.com/v1/app/app-id-1/purchases?environment=sandbox&fromDate=2019-10-11T17%3A34%3A33Z&limit=40&order=ask&originalPurchase=orig-purchase-1&page=3&toDate=2019-10-13T17%3A34%3A33Z&user=user-internal-1&userId=user-id-1"
			if req.URL.String() != expectedUrl {
				return nil, fmt.Errorf("wrong URL; expected:\n%s\ngot:\n%s", expectedUrl, req.URL.String())
			}

			body := `{"hasNextPage":true,"list":[{"id":"purchase-1","purchaseDate":"2019-10-12T17:34:33.256Z","quantity":1,"platform":"ios","country":"US","tags":{},"orderId":"9873637705964380","app":"5d86507259e828b8fe321f7e","userId":"62785074-8f32-42a5-b86b-90dbd79ce212","product":"5d86507259e828b8fe321f8a","listing":"5d86507259e828b8fe321f32","receipt":"5d86507259e828b8fe321f34","store":"5d86507259e828b8fe321f85","currency":"USD","price":19.99,"convertedCurrency":"USD","convertedPrice":19.99,"isSandbox":false,"isRefunded":false,"isSubscription":true,"isSubscriptionActive":true,"isSubscriptionRenewable":true,"isSubscriptionRetryPeriod":false,"isTrialConversion":false,"subscriptionPeriodType":"normal","expirationDate":"2019-11-12T17:34:33.256Z","linkedPurchase":"2d865c10c41280ba7f0ce9c4","originalPurchase":"2d865c10c41280ba7f0ce9c4","productSku":"membership_pricing1","productType":"renewable_subscription"}]}`

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

	from, _ := time.Parse(time.RFC3339, "2019-10-11T17:34:33.256Z")
	to, _ := time.Parse(time.RFC3339, "2019-10-13T17:34:33.256Z")
	getPurchasesRequest := iaphub.GetPurchasesRequest{
		Page:             3,
		Limit:            40,
		Order:            iaphub.Ask,
		FromDate:         from,
		ToDate:           to,
		User:             internalUser,
		UserId:           userId1,
		OriginalPurchase: origPurchaseId,
	}

	purchaseListActual, err := client.GetPurchases(getPurchasesRequest)

	if err != nil {
		t.Errorf("GetPurchases failed: %s", err)
	}

	purchaseListExpected := iaphub.PurchaseList{
		HasNextPage: true,
		List:        []iaphub.Purchase{dummyPurchase()},
	}
	if !reflect.DeepEqual(purchaseListActual, purchaseListExpected) {
		t.Errorf("wrong purchase list; expected:\n%#v\ngot:\n%#v\n", purchaseListExpected, purchaseListActual)
	}
}

func dummyPurchase() iaphub.Purchase {
	purchaseDate, _ := time.Parse(time.RFC3339, "2019-10-12T17:34:33.256Z")
	expirationDate, _ := time.Parse(time.RFC3339, "2019-11-12T17:34:33.256Z")
	return iaphub.Purchase{
		Id:                        purchaseId,
		PurchaseDate:              purchaseDate,
		Quantity:                  1,
		Platform:                  iaphub.PlatformIOS,
		Country:                   "US",
		OrderId:                   "9873637705964380",
		App:                       "5d86507259e828b8fe321f7e",
		UserId:                    "62785074-8f32-42a5-b86b-90dbd79ce212",
		Product:                   "5d86507259e828b8fe321f8a",
		Listing:                   "5d86507259e828b8fe321f32",
		Store:                     "5d86507259e828b8fe321f85",
		Receipt:                   "5d86507259e828b8fe321f34",
		Currency:                  "USD",
		Tags:                      map[string]string{},
		Price:                     19.99,
		ConvertedCurrency:         "USD",
		ConvertedPrice:            19.99,
		IsSandbox:                 false,
		IsRefunded:                false,
		IsSubscription:            true,
		IsSubscriptionActive:      true,
		IsSubscriptionRenewable:   true,
		IsSubscriptionRetryPeriod: false,
		IsTrialConversion:         false,
		SubscriptionPeriodType:    iaphub.SubscriptionPeriodTypeNormal,
		ExpirationDate:            expirationDate,
		LinkedPurchase:            "2d865c10c41280ba7f0ce9c4",
		OriginalPurchase:          "2d865c10c41280ba7f0ce9c4",
		ProductSku:                "membership_pricing1",
		ProductType:               iaphub.ProductTypeRenewableSubscription,
	}
}
