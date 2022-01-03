package iaphub_test

import (
	"bytes"
	"fmt"
	"github.com/n10ty/iaphub-go"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

func TestClient_GetSubscription(t *testing.T) {
	httpClient := newClient(
		func(req *http.Request) (*http.Response, error) {
			expectedUrl := fmt.Sprintf("https://api.iaphub.com/v1/app/%s/subscription/%s?environment=sandbox", appId1, origPurchaseId)
			if req.URL.String() != expectedUrl {
				return nil, fmt.Errorf("wrong URL; expected: %s, got: %s", expectedUrl, req.URL.String())
			}

			body := `{"id":"purchase-1","purchaseDate":"2019-10-12T17:34:33.256Z","quantity":1,"platform":"ios","country":"US","tags":{},"orderId":"9873637705964380","app":"5d86507259e828b8fe321f7e","subscription":"5d865c10c41280ba7f0ce9c2","userId":"62785074-8f32-42a5-b86b-90dbd79ce212","product":"5d86507259e828b8fe321f8a","listing":"5d86507259e828b8fe321f32","store":"5d86507259e828b8fe321f85","receipt":"5d86507259e828b8fe321f34","currency":"USD","price":19.99,"convertedCurrency":"USD","convertedPrice":19.99,"isSandbox":false,"isRefunded":false,"isSubscription":true,"isSubscriptionActive":true,"isSubscriptionRenewable":true,"isSubscriptionRetryPeriod":false,"isTrialConversion":false,"subscriptionPeriodType":"normal","expirationDate":"2019-11-12T17:34:33.256Z","linkedPurchase":"2d865c10c41280ba7f0ce9c4","originalPurchase":"2d865c10c41280ba7f0ce9c4","productSku":"membership_pricing1","productType":"renewable_subscription"}`

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

	getSubscriptionRequest := iaphub.GetSubscriptionRequest{
		OriginalPurchaseId: origPurchaseId,
	}

	subscription, err := client.GetSubscription(getSubscriptionRequest)

	if err != nil {
		t.Errorf("GetSuscription failed: %s", err)
	}

	if !reflect.DeepEqual(dummyPurchase(), subscription) {
		t.Errorf("wrong subscription; expected:\n%#v\ngot:\n%#v\n", dummyPurchase(), subscription)
	}
}
