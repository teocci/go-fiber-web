// Package model
// Created by RTT.
// Author: teocci@yandex.com on 2023-2월-27
package model

import (
	"encoding/json"
	"fmt"
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

func (sr *SellerResponse) GetJSON(sellerID string) (err error) {
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

	return err
}
