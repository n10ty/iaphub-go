package iaphub

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

// Get User
type Platform string

// App environment (production by default)
type Env string

// Direction
type Order string

const (
	PlatformIOS     Platform = "ios"
	PlatformAndroid Platform = "android"

	EnvProduction Env = "production"

	Ask  Order = "ask"
	Desc Order = "desc"
)

const ApiUrl = "https://api.iaphub.com/v1"

const (
	pathGetUser         = "/app/%s/user/%s"
	pathMigrateUser     = "/app/%s/user/%s/migrate"
	pathUpdateUser      = "/app/%s/user/%s"
	pathGetReceipt      = "/app/%s/receipt/%s"
	pathUpdateReceipt   = "/app/%s/user/%s/receipt"
	pathGetPurchase     = "/app/%s/purchase/%s"
	pathGetPurchases    = "/app/%s/purchases"
	pathGetSubscription = "/app/%s/subscription/%s"
)

type Client struct {
	apiKey string
	appId  string
	client *http.Client
	env    Env
}

func NewClient(apiKey string, appId string, options ...Option) (*Client, error) {
	config := &config{
		requestTimeout: 3 * time.Second,
		client:         http.DefaultClient,
		env:            EnvProduction,
	}

	for _, o := range options {
		err := o(config)
		if err != nil {
			return nil, err
		}
	}

	return &Client{
		apiKey: apiKey,
		appId:  appId,
		client: config.client,
		env:    config.env,
	}, nil
}

func (c *Client) requestGet(path string, queryParams map[string]string) ([]byte, error) {
	fpath := ApiUrl + path

	request, err := c.newRequest(http.MethodGet, fpath, queryParams, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(request)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		if err = resp.Body.Close(); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%d: %s", resp.StatusCode, string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()

	return body, err
}

func (c *Client) requestPost(path string, queryParams map[string]string, data interface{}) ([]byte, error) {
	fpath := ApiUrl + path

	request, err := c.newRequest(http.MethodPost, fpath, queryParams, data)
	if err != nil {
		return nil, err
	}
	resp, err := c.client.Do(request)
	if err != nil {
		return nil, err
	} else if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		} else if err = resp.Body.Close(); err != nil {
			return nil, err
		}

		return nil, fmt.Errorf("%d: %s", resp.StatusCode, string(body))
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	err = resp.Body.Close()

	return body, err
}

func (c *Client) newRequest(method string, path string, params map[string]string, data interface{}) (*http.Request, error) {
	baseURL, err := url.Parse(path)
	if err != nil {
		return nil, err
	}
	ps := url.Values{}
	if params != nil {
		for k, v := range params {
			ps.Set(k, v)
		}
	}
	ps.Set("environment", string(c.env))
	baseURL.RawQuery = ps.Encode()

	req, err := http.NewRequest(method, baseURL.String(), nil)
	if err != nil {
		return nil, err
	}

	if data != nil {
		jsonString, err := json.Marshal(data)
		if err != nil {
			return nil, err
		}
		req, err = http.NewRequest(method, baseURL.String(), bytes.NewBuffer(jsonString))
		if err != nil {
			return nil, err
		}
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Set("Authorization", "ApiKey "+c.apiKey)

	return req, err
}

// UseClient sets custom HTTP client.
func UseClient(httpClient *http.Client) Option {
	return func(c *config) error {
		if httpClient == nil {
			return errors.New("HTTP client is not specified")
		}
		c.client = httpClient

		return nil
	}
}

func UseEnv(env Env) Option {
	return func(c *config) error {
		if env == "" {
			return errors.New("environment is not specified")
		}
		c.env = env

		return nil
	}
}

type config struct {
	requestTimeout time.Duration
	client         *http.Client
	env            Env
}

type Option func(*config) error
