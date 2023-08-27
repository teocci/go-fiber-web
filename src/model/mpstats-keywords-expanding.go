// Package model
// Created by Teocci.
// Author: teocci@yandex.com on 2023-Aug-27
package model

import (
	"encoding/json"
	"os"
	"strings"

	"github.com/teocci/go-fiber-web/src/utils"
)

type MPSKeywordsExpandingResponse struct {
	Result []struct {
		Count int      `json:"count"`
		Word  string   `json:"word"`
		Words []string `json:"words"`
		Keys  []struct {
			Word    string `json:"word"`
			Count   int    `json:"count"`
			Wbcount int    `json:"wbcount"`
			Total   int    `json:"total"`
			Mp      int    `json:"mp"`
		} `json:"keys"`
		KeysCountSum   int `json:"keys_count_sum"`
		KeysWbCountSum int `json:"keys_wb_count_sum"`
		Mp             int `json:"mp"`
	} `json:"result"`
	Words []struct {
		Word            string `json:"word"`
		Count           int    `json:"count"`
		Wbcount         int    `json:"wbcount"`
		Total           int    `json:"total"`
		Mp              int    `json:"mp"`
		PrioritySubject *struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"prioritySubject"`
	} `json:"words"`
	QueryData []struct {
		Sku        string `json:"sku"`
		Brand      string `json:"brand"`
		Seller     string `json:"seller"`
		SupplierId int    `json:"supplierId"`
		Name       string `json:"name"`
		Mp         int    `json:"mp"`
	} `json:"queryData"`
}

type MPSKeywordsExpandingRequest struct {
	SKUList        []string      `json:"-"`
	QueryData      string        `json:"queryData"`
	Type           string        `json:"type"`
	Mp             int           `json:"mp,omitempty"`
	StopWords      []interface{} `json:"stopWords,omitempty"`
	Similar        bool          `json:"similar"`
	SearchFullWord bool          `json:"searchFullWord"`
	D1             string        `json:"d1"`
	D2             string        `json:"d2"`
}

func (keReq *MPSKeywordsExpandingRequest) GenerateQueryData() (queryData string) {
	if keReq.SKUList == nil || len(keReq.SKUList) == 0 {
		return ""
	}

	keReq.QueryData = strings.Join(keReq.SKUList, ",")

	return keReq.QueryData
}

func (ke *MPSKeywordsExpandingResponse) GetJSON(req MPSKeywordsExpandingRequest) (err error) {
	url := "https://mpstats.io/api/keywords/expanding"

	headers := map[string]string{
		"X-Mpstats-TOKEN": os.Getenv("MPS_API_SECRET"),
		"Content-Type":    "application/json",
	}
	r, err := utils.PostWithHeaders(url, req, headers)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&ke)
	if err != nil {
		return err
	}

	return nil
}
