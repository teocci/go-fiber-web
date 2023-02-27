// Package model
// Created by RTT.
// Author: teocci@yandex.com on 2023-2월-27
package model

import (
	"encoding/json"
	"fmt"
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

func (sr *SellerResponse) GetJSON(url string) (err error) {
	fmt.Printf("JAVA_API_URL: %#v\n", url)

	r, err := httpClient.Get(url)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&sr)

	return err
}
