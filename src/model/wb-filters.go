// Package model
// Created by RTT.
// Author: teocci@yandex.com on 2023-2월-28
package model

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type FilterResponse struct {
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

func (pfr *FilterResponse) GetJSON(supplierID string) (err error) {
	baseURL := &url.URL{
		Scheme: "https",
		Host:   "catalog.wb.ru",
		Path:   "/sellers/v4/filters",
	}

	params := baseURL.Query()
	params.Set("appType", "1")
	params.Set("curr", "rub")
	params.Set("dest", "-1257786")
	params.Set("reg", "0")
	params.Set("regions", "80,38,83,4,64,33,68,70,30,40,86,75,69,22,1,31,66,110,48,71,114")
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
