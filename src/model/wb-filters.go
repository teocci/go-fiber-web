// Package model
// Created by RTT.
// Author: teocci@yandex.com on 2023-2월-28
package model

import (
	"encoding/json"
	"fmt"
	"net/url"
)

const (
	ModeSeller   = "seller"
	ModeCategory = "category"
)

type FilterRequest struct {
	ID   string `json:"id"`
	Mode string `json:"mode"`
}

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

func (pfr *FilterResponse) GetJSON(req FilterRequest) (err error) {
	if req.ID == "" {
		return fmt.Errorf("invalid id: null")
	}

	if req.Mode == "" {
		req.Mode = ModeSeller
	}

	var baseURL *url.URL
	switch req.Mode {
	case ModeSeller:
		baseURL = &url.URL{
			Scheme: "https",
			Host:   "catalog.wb.ru",
			Path:   "/sellers/v4/filters",
		}
	case ModeCategory:
		shard := "beauty3"
		if req.ID == "9000" {
			shard = "beauty4"
		}

		baseURL = &url.URL{
			Scheme: "https",
			Host:   "catalog.wb.ru",
			Path:   fmt.Sprintf("/catalog/%s/v4/filters", shard),
		}
	}

	params := baseURL.Query()
	params.Set("appType", "1")
	params.Set("curr", "rub")
	params.Set("dest", "-1257786")
	params.Set("reg", "0")
	params.Set("regions", "80,38,83,4,64,33,68,70,30,40,86,75,69,22,1,31,66,110,48,71,114")
	params.Set("spp", "0")

	switch req.Mode {
	case ModeSeller:
		params.Set("supplier", req.ID)
	case ModeCategory:
		params.Set("cat", req.ID)
	}

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
