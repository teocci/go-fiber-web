// Package model
// Created by Teocci.
// Author: teocci@yandex.com on 2023-Aug-16
package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/teocci/go-fiber-web/src/utils"
	"github.com/teocci/go-fiber-web/src/utils/mapslice"
	"net/url"
	"os"
)

type WordsData struct {
	Word    string `json:"word,omitempty"`
	Type    string `json:"type,omitempty"`
	Pos     []int  `json:"pos"`
	Count   int    `json:"count"`
	WbCount int    `json:"wb_count"`
	Total   int    `json:"total"`
	AvgPos  int    `json:"avgPos"`
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

func (s *MPStatsKeywords) GetJSON(productId int) (err error) {
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

	err = json.NewDecoder(r.Body).Decode(&s)

	s.ProcessData()

	return err
}

func (s *MPStatsKeywords) ProcessData() {
	for k, v := range s.Words {
		v.Word = k
		s.Words[k] = v
	}
}

func (s *MPStatsKeywords) SortByWbCount() []string {
	sortedWords := mapslice.New(s.Words)
	sortedWords.SortBy(func(a, b mapslice.Entry[WordsData]) bool {
		return a.Value.Count > b.Value.Count
	})

	keys := sortedWords.Keys()

	s.Words = make(map[string]WordsData)
	for _, entries := range sortedWords.Entries() {
		entries.Value.Word = entries.Key
		s.Words[entries.Key] = entries.Value
		if len(s.Words) < 11 {
			fmt.Printf("SortByWbCount - Keyword: %s[wb-count: %d]\n", entries.Key, entries.Value.WbCount)
		}
	}

	return keys
}

func (s *MPStatsKeywords) Clone() MPStatsKeywords {
	clonedWords := make(map[string]WordsData)
	for key, value := range s.Words {
		clonedWords[key] = value
	}

	return MPStatsKeywords{
		Words:      clonedWords,
		Days:       append([]string(nil), s.Days...),
		Sales:      append([]int(nil), s.Sales...),
		Balance:    append([]int(nil), s.Balance...),
		FinalPrice: append([]int(nil), s.FinalPrice...),
		Comments:   append([]int(nil), s.Comments...),
		Rating:     append([]int(nil), s.Rating...),
	}
}
