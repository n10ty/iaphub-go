package iaphub

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
)

type GetUserRequest struct {
	UserId   string
	Platform Platform
	Upsert   bool
}

func (c *Client) GetUser(request GetUserRequest) (User, error) {
	var user User

	if request.UserId == "" || request.Platform == "" {
		return user, fmt.Errorf("required parameter \"userId\" or \"platform\" is missing")
	}

	path := fmt.Sprintf(pathGetUser, c.appId, request.UserId)

	params := map[string]string{
		"platform": string(request.Platform),
	}

	if request.Upsert {
		params["upsert"] = strconv.FormatBool(true)
	}

	response, err := c.requestGet(path, params)
	if err != nil {
		return user, err
	}

	err = json.Unmarshal(response, &user)
	return user, err
}

type GetUserMigrateRequest struct {
	UserId string
}

func (c *Client) GetUserMigrate(request GetUserMigrateRequest) (LatestUser, error) {
	var latestUser LatestUser

	if request.UserId == "" {
		return latestUser, errors.New("required parameter")
	}

	var params map[string]string
	path := fmt.Sprintf(pathMigrateUser, c.appId, request.UserId)
	response, err := c.requestGet(path, params)
	if err != nil {
		return latestUser, err
	}

	err = json.Unmarshal(response, &latestUser)

	return latestUser, err
}

type UpdateUserRequest struct {
	UserId  string            `json:"userId"`
	Country string            `json:"country"`
	Upsert  bool              `json:"upsert"`
	Env     Env               `json:"environment,omitempty"`
	Tags    map[string]string `json:"tags"`
}

func (c *Client) UpdateUser(request UpdateUserRequest) error {
	if request.UserId == "" || request.Country == "" {
		return fmt.Errorf("required parameter \"userId\" or \"country\" is missing")
	}

	if request.Env == "" {
		request.Env = c.env
	}
	path := fmt.Sprintf(pathUpdateUser, c.appId, request.UserId)

	_, err := c.requestPost(path, map[string]string{}, request)

	return err
}

type User struct {
	ProductForSale []Product `json:"productsForSale"`
	ActiveProducts []Product `json:"activeProducts"`
}

type Product struct {
	Id       string `json:"id"`
	Type     string `json:"type"`
	Sku      string `json:"sku"`
	Purchase string `json:"purchase"`
}

type LatestUser struct {
	UserId string `json:"userId"`
}
