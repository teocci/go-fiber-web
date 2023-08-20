// Package model
// Created by Teocci.
// Author: teocci@yandex.com on 2023-Aug-16
package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"time"

	"github.com/teocci/go-fiber-web/src/scache"
	"github.com/teocci/go-fiber-web/src/utils"
	"github.com/teocci/go-fiber-web/src/utils/mapslice"
)

type WordsData struct {
	Count   int    `json:"count"`
	WbCount int    `json:"wb_count"`
	Total   int    `json:"total"`
	AvgPos  int    `json:"avgPos"`
	Word    string `json:"word,omitempty"`
	Type    string `json:"type,omitempty"`
	Pos     []int  `json:"pos,omitempty"`
}

type MPStatsKeywords struct {
	Days       []string             `json:"days"`
	Sales      []int                `json:"sales"`
	Balance    []int                `json:"balance"`
	FinalPrice []int                `json:"final_price"`
	Comments   []int                `json:"comments"`
	Rating     []int                `json:"rating"`
	Words      map[string]WordsData `json:"words,omitempty"`
}

type RequestRate struct {
	Consecutive int
	LastTime    time.Time
}

const mpstMaxConsecutiveRequests = 200
const mpstMaxDuration = 1 * time.Minute
const mpstPauseDuration = 20 * time.Second

var (
	mpstCache     = scache.GetCacheInstance[MPStatsKeywords]("MPStats")
	mpstRateLimit = RequestRate{
		Consecutive: 0,
		LastTime:    time.Now(),
	}
)

// TODO: replace nm with cacheKey
func (s *MPStatsKeywords) checkCache(nm int) (ok bool) {
	ok = false
	if nm == 0 {
		return false
	}

	cacheKey := fmt.Sprintf("mpstats-keywords-%d", nm)
	fmt.Printf("MPStatsKeywords checkCache cacheKey: %s\n", cacheKey)

	var sCache *MPStatsKeywords
	if sCache, ok = mpstCache.Get(cacheKey); ok {
		fmt.Println("mpstats-keywords - Cache hit")
		*s = *sCache
		return true
	}

	return
}

func (s *MPStatsKeywords) updateCache(nm int) {
	cacheKey := fmt.Sprintf("mpstats-keywords-%d", nm)
	fmt.Printf("MPStatsKeywords updateCache cacheKey: %s\n", cacheKey)
	mpstCache.Set(cacheKey, s.Clone(), 0)
}

func (s *MPStatsKeywords) GetJSON(nm int) (err error) {
	if nm == 0 {
		return errors.New("invalid id: null")
	}

	found := s.checkCache(nm)
	if found {
		return nil
	}

	checkConsecutiveRequests()

	baseURL := &url.URL{
		Scheme: "https",
		Host:   "mpstats.io",
		Path:   fmt.Sprintf("/api/wb/get/item/%d/by_keywords", nm),
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
		"Content-Type":    "application/json",
	}
	r, err := utils.GetWithHeaders(apiURL, headers)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	code := r.StatusCode
	fmt.Printf("Response code: %d\n", code)

	var body []byte
	body, err = io.ReadAll(r.Body)
	if err != nil {
		return err
	}

	tmp := interface{}(nil)

	err = json.Unmarshal(body, &tmp)
	if err != nil {
		return err
	}

	fmt.Printf("Response status: %s\n", r.Status)
	//fmt.Printf("Response tmp: %#v\n", tmp)

	fmt.Println("Request convertInterfaceToMPStatsKeywords")
	err = s.convertInterfaceToMPStatsKeywords(tmp)
	if err != nil {
		return err
	}

	fmt.Println("Request convertInterfaceToMPStatsKeywords ended")

	//err = json.NewDecoder(r.Body).Decode(&s)
	//if err != nil {
	//	return err
	//}

	updateConsecutiveRequests()

	s.ProcessData()
	s.updateCache(nm)

	return nil
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

func (s *MPStatsKeywords) CopyFrom(d MPStatsKeywords) {
	s.Words = d.Words
	s.Days = d.Days
	s.Sales = d.Sales
	s.Balance = d.Balance
	s.FinalPrice = d.FinalPrice
	s.Comments = d.Comments
	s.Rating = d.Rating
}

func (s *MPStatsKeywords) Clone() MPStatsKeywords {
	clonedWords := make(map[string]WordsData)
	for key, value := range s.Words {
		clonedWords[key] = value
	}

	return MPStatsKeywords{
		Days:       append([]string(nil), s.Days...),
		Sales:      append([]int(nil), s.Sales...),
		Balance:    append([]int(nil), s.Balance...),
		FinalPrice: append([]int(nil), s.FinalPrice...),
		Comments:   append([]int(nil), s.Comments...),
		Rating:     append([]int(nil), s.Rating...),
		Words:      clonedWords,
	}
}

func (w *WordsData) convertInterfaceToWordsData(data interface{}) (err error) {
	wordDataMap, ok := data.(map[string]interface{})
	if !ok {
		return errors.New("type assertion failed for WordsData")
	}

	w.Pos, _ = toIntSlice(wordDataMap["pos"])
	w.Count = int(wordDataMap["count"].(float64))
	w.WbCount = int(wordDataMap["wb_count"].(float64))
	w.Total = int(wordDataMap["total"].(float64))
	w.AvgPos = int(wordDataMap["avgPos"].(float64))

	return nil
}

func (s *MPStatsKeywords) convertInterfaceToMPStatsKeywords(data interface{}) (err error) {
	fmt.Println("convertInterfaceToMPStatsKeywords started")

	fmt.Println("data assertion started")
	dataMap, ok := data.(map[string]interface{})
	fmt.Println("data assertion ended")
	if !ok {
		return errors.New("type assertion failed for MPStatsKeywords")
	}
	fmt.Println("dataMap ready")

	// Convert other fields like 'days', 'sales', etc.
	s.Days, _ = toStringSlice(dataMap["days"])
	s.Sales, _ = toIntSlice(dataMap["sales"])
	s.Balance, _ = toIntSlice(dataMap["balance"])
	s.FinalPrice, _ = toIntSlice(dataMap["final_price"])
	s.Comments, _ = toIntSlice(dataMap["comments"])
	s.Rating, _ = toIntSlice(dataMap["rating"])

	// Convert 'words' map
	fmt.Println("wordsDataInterface assertion started")
	//fmt.Printf("dataMap: %#v\n", dataMap["words"])
	wordsDataInterface, ok := dataMap["words"]
	//fmt.Printf("wordsDataInterface: %#v\n", wordsDataInterface)
	if !ok {
		return errors.New("'words' field not found")
	}
	fmt.Println("wordsDataInterface assertion ended")

	if wordsDataInterface == nil {
		s.Words = make(map[string]WordsData)
		return nil
	} else {
		fmt.Println("wordsMap assertion started")
		wordsMap, ok := wordsDataInterface.(map[string]interface{})
		//fmt.Printf("wordsMap: %#v\n", wordsMap)
		if wordsMap == nil {
			s.Words = make(map[string]WordsData)
			return
		}
		fmt.Println("wordsMap assertion ended")
		if !ok {
			return errors.New("'words' field not found or type assertion failed")
		}
		fmt.Println("wordsMap ready")

		fmt.Println("convert words started")
		s.Words = make(map[string]WordsData)
		for key, value := range wordsMap {
			wordData := WordsData{}
			err = wordData.convertInterfaceToWordsData(value)
			if err != nil {
				return err
			}
			s.Words[key] = wordData
		}
		fmt.Println("convert words ended")
	}

	fmt.Println("convertInterfaceToMPStatsKeywords ended")

	return nil
}

func toStringSlice(value interface{}) ([]string, error) {
	if slice, ok := value.([]interface{}); ok {
		result := make([]string, len(slice))
		for i, item := range slice {
			result[i] = item.(string)
		}
		return result, nil
	}
	return nil, fmt.Errorf("type assertion failed for string slice")
}

func toIntSlice(value interface{}) ([]int, error) {
	if slice, ok := value.([]interface{}); ok {
		result := make([]int, len(slice))
		for i, item := range slice {
			result[i] = int(item.(float64))
		}
		return result, nil
	}
	return nil, fmt.Errorf("type assertion failed for int slice")
}

func checkConsecutiveRequests() {
	fmt.Printf("checkConsecutiveRequests Rate limit: %d/%d\n", mpstRateLimit.Consecutive, mpstMaxConsecutiveRequests)
	if time.Since(mpstRateLimit.LastTime) < mpstMaxDuration && mpstRateLimit.Consecutive > mpstMaxConsecutiveRequests {
		fmt.Printf("Rate limit exceeded, pausing for %0.00f\n", mpstPauseDuration.Seconds())
		time.Sleep(mpstPauseDuration)
		mpstRateLimit.Consecutive = 0
		mpstRateLimit.LastTime = time.Now()
	}
}

func updateConsecutiveRequests() {
	mpstRateLimit.Consecutive++
	if time.Since(mpstRateLimit.LastTime) > mpstMaxDuration {
		mpstRateLimit.Consecutive = 0
	}

	mpstRateLimit.LastTime = time.Now()
	fmt.Printf("updateConsecutiveRequests Rate limit: %d/%d\n", mpstRateLimit.Consecutive, mpstMaxConsecutiveRequests)
}
