// Package model
// Created by RTT.
// Author: teocci@yandex.com on 2023-2월-27
package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/teocci/go-fiber-web/src/scache"
	"io"
	"net/url"
)

type SellerResponse struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	FineName     string `json:"fineName"`
	OGRN         string `json:"ogrn"`
	Trademark    string `json:"trademark"`
	LegalAddress string `json:"legalAddress"`
	IsUnknown    bool   `json:"isUnknown"`
}

const (
	urlFormatSeller = "%s/%s"
)

var (
	wbSellerCache = scache.GetCacheInstance[SellerResponse]("WBSeller")
)

func (sr *SellerResponse) checkCache(cacheKey string) (ok bool) {
	ok = false
	if cacheKey == "" {
		return
	}

	var cacheData *SellerResponse
	if cacheData, ok = wbSellerCache.Get(cacheKey); ok {
		fmt.Println("wb-seller cache hit")
		*sr = *cacheData
		return true
	}

	return false
}

func (sr *SellerResponse) updateCache(cacheKey string) {
	if cacheKey == "" {
		return
	}

	cloned := sr.Clone()
	wbSellerCache.Set(cacheKey, cloned, 0)
}

func (sr *SellerResponse) Clone() SellerResponse {
	return SellerResponse{
		ID:           sr.ID,
		Name:         sr.Name,
		FineName:     sr.FineName,
		OGRN:         sr.OGRN,
		Trademark:    sr.Trademark,
		LegalAddress: sr.LegalAddress,
		IsUnknown:    sr.IsUnknown,
	}
}

func (sr *SellerResponse) GetJSON(sellerID string) (err error) {
	if sellerID == "" {
		return errors.New("invalid sellerID: null")
	}

	cacheKey := fmt.Sprintf("wb-seller-data-%s", sellerID)
	found := sr.checkCache(cacheKey)
	if found {
		return nil
	}

	baseURL := &url.URL{
		Scheme: "https",
		Host:   "www.wildberries.ru",
		Path:   "/webapi/seller/data/short",
	}

	apiURL := fmt.Sprintf(urlFormatSeller, baseURL.String(), sellerID)
	fmt.Printf("WB_API_URL: %#v\n", apiURL)

	r, err := httpClient.Get(apiURL)
	if err != nil {
		return err
	}
	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(r.Body)

	err = json.NewDecoder(r.Body).Decode(&sr)

	sr.updateCache(cacheKey)

	return err
}
