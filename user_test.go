package iaphub_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rentaapp/iaphub-go"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

var (
	userId1 = "user-id-1"
	appId1  = "app-id-1"
	apiKey1 = "api-key-1"
	env     = "sandbox"
	tagKey1 = "tag-key-1"
	tagVal1 = "tag-val-1"
)

func TestClient_GetUser(t *testing.T) {
	httpClient := newClient(
		func(req *http.Request) (*http.Response, error) {
			expectedUrl := fmt.Sprintf("https://api.iaphub.com/v1/app/%s/user/%s?environment=sandbox&platform=android&upsert=true", appId1, userId1)
			if req.URL.String() != expectedUrl {
				return nil, fmt.Errorf("wrong URL; expected: %s, got: %s", expectedUrl, req.URL.String())
			}

			authHeaderVal := "ApiKey " + apiKey1
			if req.Header.Get("Authorization") != authHeaderVal {
				return nil, fmt.Errorf("wrong auth header; expected: %s, got: %s", authHeaderVal, req.Header.Get("Authorization"))
			}

			body, _ := json.Marshal(dummyUser())
			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
			}, nil
		},
	)

	client, _ := iaphub.NewClient(
		apiKey1,
		appId1,
		iaphub.UseClient(httpClient),
		iaphub.UseEnv(iaphub.Env(env)),
	)

	getUserRequest := iaphub.GetUserRequest{
		UserId:   userId1,
		Platform: iaphub.PlatformAndroid,
		Upsert:   true,
	}

	user, err := client.GetUser(getUserRequest)

	if err != nil {
		t.Errorf("GetUser failed: %s", err)
	}

	if !reflect.DeepEqual(user, dummyUser()) {
		t.Errorf("wrong user; expected: %#v, got: %#v", dummyUser(), user)
	}
}

func TestClient_GetUserUrl(t *testing.T) {
	type fields struct {
		apiKey   string
		appId    string
		env      iaphub.Env
		userId   string
		platform iaphub.Platform
		upsert   bool
	}

	tests := []struct {
		name        string
		fields      fields
		expectedUrl string
	}{
		{
			"Test case 1",
			fields{
				apiKey1, appId1, iaphub.Env(env), userId1, iaphub.PlatformIOS, false,
			},
			fmt.Sprintf("https://api.iaphub.com/v1/app/%s/user/%s?environment=sandbox&platform=ios", appId1, userId1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			httpClient := newClient(
				func(req *http.Request) (*http.Response, error) {
					if req.URL.String() != tt.expectedUrl {
						return nil, fmt.Errorf("wrong URL; expected: %s, got: %s", tt.expectedUrl, req.URL.String())
					} else {
						body, _ := json.Marshal(dummyUser())
						return &http.Response{StatusCode: http.StatusOK, Body: ioutil.NopCloser(bytes.NewBuffer(body))}, nil
					}
				},
			)

			client, _ := iaphub.NewClient(
				apiKey1,
				appId1,
				iaphub.UseClient(httpClient),
				iaphub.UseEnv(tt.fields.env),
			)

			getUserRequest := iaphub.GetUserRequest{
				UserId:   tt.fields.userId,
				Platform: tt.fields.platform,
				Upsert:   tt.fields.upsert,
			}

			_, err := client.GetUser(getUserRequest)
			if err != nil {
				t.Errorf(err.Error())
			}
		})
	}
}

func TestClient_GetUserMigrate(t *testing.T) {

	httpClient := newClient(
		func(req *http.Request) (*http.Response, error) {
			expectedUrl := fmt.Sprintf("https://api.iaphub.com/v1/app/%s/user/%s/migrate?environment=sandbox", appId1, userId1)
			if req.URL.String() != expectedUrl {
				return nil, fmt.Errorf("wrong URL; expected: %s, got: %s", expectedUrl, req.URL.String())
			}

			authHeaderVal := "ApiKey " + apiKey1
			if req.Header.Get("Authorization") != authHeaderVal {
				return nil, fmt.Errorf("wrong auth header; expected: %s, got: %s", authHeaderVal, req.Header.Get("Authorization"))
			}

			body, _ := json.Marshal(dummyLatestUser())

			return &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(bytes.NewBuffer(body)),
			}, nil
		},
	)

	client, _ := iaphub.NewClient(
		apiKey1,
		appId1,
		iaphub.UseClient(httpClient),
		iaphub.UseEnv(iaphub.Env(env)),
	)

	getUserRequest := iaphub.GetUserMigrateRequest{
		UserId: userId1,
	}

	user, err := client.GetUserMigrate(getUserRequest)

	if err != nil {
		t.Errorf("GetUserMigrate failed: %s", err)
		return
	}

	if !reflect.DeepEqual(user, dummyLatestUser()) {
		t.Errorf("wrong user migration; expected: %#v, got: %#v", dummyLatestUser(), user)
	}
}

func TestClient_UpdateUser(t *testing.T) {
	httpClient := newClient(
		func(req *http.Request) (*http.Response, error) {
			if req.Method != http.MethodPost {
				return nil, fmt.Errorf("method not supported\n")
			}
			expectedUrl := fmt.Sprintf("https://api.iaphub.com/v1/app/%s/user/%s?environment=sandbox", appId1, userId1)
			if req.URL.String() != expectedUrl {
				return nil, fmt.Errorf("wrong URL; expected: %s, got: %s", expectedUrl, req.URL.String())
			}

			expectedBody := fmt.Sprintf(`{"userId":"%s","country":"Ukraine","upsert":false,"environment":"sandbox","tags":{"tag-key-1":"tag-val-1"}}`, userId1)

			actualBody, _ := ioutil.ReadAll(req.Body)

			if expectedBody != string(actualBody) {
				return nil, fmt.Errorf("wrong request body; expected: %s, got: %s", expectedBody, actualBody)
			}

			return &http.Response{
				StatusCode: 200,
			}, nil
		},
	)

	client, _ := iaphub.NewClient(
		apiKey1,
		appId1,
		iaphub.UseClient(httpClient),
		iaphub.UseEnv(iaphub.Env(env)),
	)

	updateUserRequest := iaphub.UpdateUserRequest{
		UserId:  userId1,
		Upsert:  false,
		Country: "Ukraine",
		Tags: map[string]string{
			tagKey1: tagVal1,
		},
	}

	err := client.UpdateUser(updateUserRequest)

	if err != nil {
		t.Errorf("UpdateUser failed: %s", err)
	}
}

func newClient(roundTripper roundTripper) *http.Client {
	return &http.Client{
		Transport: roundTripper,
	}
}

type roundTripper func(req *http.Request) (*http.Response, error)

func (r roundTripper) RoundTrip(request *http.Request) (*http.Response, error) {
	return r(request)
}

func dummyUser() iaphub.User {
	return iaphub.User{
		ProductForSale: []iaphub.Product{
			{
				Id:       "1",
				Type:     "non_consumable",
				Sku:      "sku1",
				Purchase: "id1",
			},
		},
		ActiveProducts: []iaphub.Product{
			{
				Id:       "2",
				Type:     "non_consumable",
				Sku:      "sku2",
				Purchase: "id2",
			},
		},
	}
}

func dummyLatestUser() iaphub.LatestUser {
	return iaphub.LatestUser{
		UserId: userId1,
	}
}
