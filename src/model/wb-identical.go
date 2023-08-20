// Package model
// Created by RTT.
// Author: teocci@yandex.com on 2023-2월-28
package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/teocci/go-fiber-web/src/scache"
	"net/url"
)

type IdenticalProductsResponse []int

var (
	wbIdenticalCache = scache.GetCacheInstance[IdenticalProductsResponse]("WBIdentical")
)

// TODO: replace nm with cacheKey
func (ipr *IdenticalProductsResponse) checkCache(nm string) (ok bool) {
	ok = false
	if nm == "" {
		return
	}

	cacheKey := fmt.Sprintf("wb-identical-%s", nm)
	var cacheData *IdenticalProductsResponse
	if cacheData, ok = wbIdenticalCache.Get(cacheKey); ok {
		fmt.Println("wb-identical cache hit")
		*ipr = *cacheData
		return true
	}

	return false
}

func (ipr *IdenticalProductsResponse) updateCache(nm string) {
	if nm == "" {
		return
	}

	cacheKey := fmt.Sprintf("wb-identical-%s", nm)
	cloned := ipr.Clone()
	wbIdenticalCache.Set(cacheKey, cloned, 0)
}

func (ipr *IdenticalProductsResponse) Clone() IdenticalProductsResponse {
	clone := make(IdenticalProductsResponse, len(*ipr))
	copy(clone, *ipr)
	return clone
}

func (ipr *IdenticalProductsResponse) GetJSON(nm string) (err error) {
	if nm == "" {
		return errors.New("invalid nm: null")
	}

	found := ipr.checkCache(nm)
	if found {
		return nil
	}

	baseURL := &url.URL{
		Scheme: "https",
		Host:   "identical-products.wildberries.ru",
		Path:   "/api/v1/identical",
	}

	params := baseURL.Query()
	params.Set("nmID", nm)
	baseURL.RawQuery = params.Encode()

	apiURL := baseURL.String()
	fmt.Printf("WB_API_URL: %#v\n", apiURL)

	r, err := httpClient.Get(apiURL)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&ipr)

	ipr.updateCache(nm)

	return err
}
