// Package model
// Created by Teocci.
// Author: teocci@yandex.com on 2023-Aug-16
package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"sort"

	"github.com/teocci/go-fiber-web/src/utils"
)

type WordsData struct {
	Pos     []int `json:"pos"`
	Count   int   `json:"count"`
	WbCount int   `json:"wb_count"`
	Total   int   `json:"total"`
	AvgPos  int   `json:"avgPos"`
}

type MPStatsKeywords struct {
	Words      map[string]WordsData `json:"words"`
	Days       []string             `json:"days"`
	Sales      []int                `json:"sales"`
	Balance    []int                `json:"balance"`
	FinalPrice []int                `json:"final_price"`
	Comments   []int                `json:"comments"`
	Rating     []int                `json:"rating"`
}

func (k *MPStatsKeywords) GetJSON(productId int) (err error) {
	if productId == 0 {
		return errors.New("invalid id: null")
	}

	baseURL := &url.URL{
		Scheme: "https",
		Host:   "mpstats.io",
		Path:   fmt.Sprintf("/api/wb/get/item/%d/by_keywords", productId),
	}
	params := baseURL.Query()
	params.Set("d1", utils.GetLastWeekDate())
	params.Set("d2", utils.GetTodayDate())
	params.Set("full", "true")

	baseURL.RawQuery = params.Encode()

	apiURL := baseURL.String()
	fmt.Printf("MPS_API_URL: %#v\n", apiURL)

	headers := map[string]string{
		"X-Mpstats-TOKEN": os.Getenv("MPS_API_SECRET"),
	}
	r, err := utils.GetWithHeaders(apiURL, headers)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&k)

	return err
}

func (k *MPStatsKeywords) SortByWbCount() []string {
	wordEntries := make([]struct {
		Keyword string
		Data    WordsData
	}, 0)

	for keyword, data := range k.Words {
		wordEntries = append(wordEntries, struct {
			Keyword string
			Data    WordsData
		}{Keyword: keyword, Data: data})
	}

	sort.Slice(wordEntries, func(i, j int) bool {
		return wordEntries[i].Data.WbCount > wordEntries[j].Data.WbCount
	})

	keys := make([]string, 0, len(wordEntries))

	sortedWords := make(map[string]WordsData)
	for _, entry := range wordEntries {
		keys = append(keys, entry.Keyword)
		sortedWords[entry.Keyword] = entry.Data
		if len(sortedWords) < 11 {
			fmt.Printf("Keyword: %s[wb-count: %d]\n", entry.Keyword, entry.Data.WbCount)
		}
	}

	k.Words = sortedWords

	return keys
}
