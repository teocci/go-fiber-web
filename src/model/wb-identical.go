// Package model
// Created by RTT.
// Author: teocci@yandex.com on 2023-2월-28
package model

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type IdenticalProductsResponse []int

func (ipr *IdenticalProductsResponse) GetJSON(id string) (err error) {
	baseURL := &url.URL{
		Scheme: "https",
		Host:   "identical-products.wildberries.ru",
		Path:   "/api/v1/identical",
	}

	params := baseURL.Query()
	params.Set("nmID", id)
	baseURL.RawQuery = params.Encode()

	apiURL := baseURL.String()
	fmt.Printf("WB_API_URL: %#v\n", apiURL)

	r, err := httpClient.Get(apiURL)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&ipr)

	return err
}
