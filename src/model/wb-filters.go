// Package model
// Created by RTT.
// Author: teocci@yandex.com on 2023-2월-28
package model

import (
	"encoding/json"
	"fmt"
	"github.com/teocci/go-fiber-web/src/scache"
	"net/url"
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

var (
	wbFiltersCache = scache.GetCacheInstance[FilterResponse]("WBFilters")
)

func (pfr *FilterResponse) checkCache(cacheKey string) (ok bool) {
	ok = false
	if cacheKey == "" {
		return
	}

	var cacheData *FilterResponse
	if cacheData, ok = wbFiltersCache.Get(cacheKey); ok {
		fmt.Println("wb-filters cache hit")
		*pfr = *cacheData
		return true
	}

	return false
}

func (pfr *FilterResponse) updateCache(cacheKey string) {
	if cacheKey == "" {
		return
	}

	cloned := pfr.Clone()
	wbFiltersCache.Set(cacheKey, cloned, 0)
}

func (pfr *FilterResponse) Clone() FilterResponse {
	clone := FilterResponse{
		State: pfr.State,
		Data:  pfr.Data,
	}

	return clone
}

func (pfr *FilterResponse) GetJSON(req FilterRequest) (err error) {
	if req.ID == "" {
		return fmt.Errorf("invalid id: null")
	}

	if req.Mode == "" {
		req.Mode = ModeSeller
	}

	var cacheKey string
	var baseURL *url.URL
	cacheKey, baseURL = req.generateCacheKeyAndURL()

	found := pfr.checkCache(cacheKey)
	if found {
		return nil
	}

	apiURL := req.generateAPIURL(baseURL)
	fmt.Printf("WB_API_URL: %#v\n", apiURL)

	r, err := httpClient.Get(apiURL)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&pfr)

	pfr.updateCache(cacheKey)

	return err
}

func (fReq *FilterRequest) generateCacheKeyAndURL() (cacheKey string, baseURL *url.URL) {
	switch fReq.Mode {
	case ModeSeller:
		baseURL = &url.URL{
			Scheme: "https",
			Host:   "catalog.wb.ru",
			Path:   "/sellers/v4/filters",
		}
		cacheKey = fmt.Sprintf("wb-filters-%s-%s", fReq.Mode, fReq.ID)
	case ModeCategory:
		shard := "beauty3"
		if fReq.ID == "9000" {
			shard = "beauty4"
		}

		baseURL = &url.URL{
			Scheme: "https",
			Host:   "catalog.wb.ru",
			Path:   fmt.Sprintf("/catalog/%s/v4/filters", shard),
		}
		cacheKey = fmt.Sprintf("wb-filters-%s-%s-%s", fReq.Mode, shard, fReq.ID)
	}

	return
}

func (fReq *FilterRequest) generateAPIURL(baseURL *url.URL) string {
	params := baseURL.Query()
	params.Set("appType", "1")
	params.Set("curr", "rub")
	params.Set("dest", "-1257786")
	params.Set("reg", "0")
	params.Set("regions", "80,38,83,4,64,33,68,70,30,40,86,75,69,22,1,31,66,110,48,71,114")
	params.Set("spp", "0")

	switch fReq.Mode {
	case ModeSeller:
		params.Set("supplier", fReq.ID)
	case ModeCategory:
		params.Set("cat", fReq.ID)
	}

	baseURL.RawQuery = params.Encode()
	return baseURL.String()
}
