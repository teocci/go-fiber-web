// Package model
// Created by RTT.
// Author: teocci@yandex.com on 2023-2월-28
package model

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type ProductFilterResponse struct {
	State int `json:"state"`
	Data  struct {
		Filters []struct {
			Name      string `json:"name"`
			Key       string `json:"key"`
			Maxselect int    `json:"maxselect,omitempty"`
			Items     []struct {
				Id    int    `json:"id"`
				Name  string `json:"name"`
				Count int    `json:"count,omitempty"`
			} `json:"items,omitempty"`
			IsFull      bool `json:"isFull,omitempty"`
			MinPriceU   int  `json:"minPriceU,omitempty"`
			MaxPriceU   int  `json:"maxPriceU,omitempty"`
			Multiselect int  `json:"multiselect,omitempty"`
		} `json:"filters"`
		Previews map[int64]int64 `json:"previews"`
		Total    int             `json:"total"`
	} `json:"data"`
}

func (pfr *ProductFilterResponse) GetJSON(supplierID string) (err error) {
	baseURL := &url.URL{
		Scheme: "https",
		Host:   "catalog.wb.ru",
		Path:   "/sellers/v4/filters",
	}

	params := baseURL.Query()
	params.Set("appType", "1")
	params.Set("couponsGeo", "12,3,18,15,21")
	params.Set("curr", "rub")
	params.Set("dest", "-1257786")
	params.Set("emp", "0")
	params.Set("lang", "ru")
	params.Set("locale", "ru")
	params.Set("pricemarginCoeff", "1.0")
	params.Set("reg", "0")
	params.Set("regions", "80,64,38,4,83,33,68,70,69,30,86,75,40,1,22,66,31,48,110,71")
	params.Set("spp", "0")
	params.Set("supplier", supplierID)
	baseURL.RawQuery = params.Encode()

	apiURL := baseURL.String()
	fmt.Printf("WB_API_URL: %#v\n", apiURL)

	r, err := httpClient.Get(apiURL)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&pfr)

	return err
}
