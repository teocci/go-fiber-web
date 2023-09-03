// Package model
// Created by Teocci.
// Author: teocci@yandex.com on 2023-Aug-16
package model

import (
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/teocci/go-fiber-web/src/utils"
)

type CommonWord struct {
	Word    string `json:"words"`
	WbCount int    `json:"wb_count"`
	Count   int    `json:"count"`
	Pos     int    `json:"pos"`
}

type ProductPositionResponse struct {
	Id         int         `json:"id"`
	Name       string      `json:"name"`
	Brand      string      `json:"brand"`
	BrandId    int         `json:"brandId"`
	SupplierId int         `json:"supplierId"`
	Pos        int         `json:"pos"`
	WordsData  []WordsData `json:"words,omitempty"`
}

type ProductPositionListResponse struct {
	Keywords []string                  `json:"keywords"`
	Products []ProductPositionResponse `json:"products"`
}

type infoData struct {
	Product ProductPositionResponse
	Stats   MPStatsKeywords
}

const (
	maxKeywords     = 15
	wordsPerProduct = 30
)

var (
	brandBlackList = []string{
		"tom ford",
		"chanel",
		"Рив гош",
		"antonio banderas",
		"Avon",
		"Trusardi",
		"Love is",
		"lacoste",
		"christine lavoisier",
		"эйвон",
		"лакост",
		"антонио бандерас",
		"наркотик",
		"anua",
		"шаман",
		"lalique",
	}

	blackListedWords = []string{
		"для волос",
		"tom ford",
		"база",
		"7days",
		"белорусская косметика",
		"chanel",
		"молекула 02",
		"Мужские ароматы - туалетная вода",
		"Lacost",
		"Рив гош",
		"Antonio banderas",
		"Виски",
		"Lacost духи",
		"Lacost мужской",
		"Антонио Бандерас для мужчин",
		"Женские ароматы -туалетная вода",
		"Avon",
		"Hallow kitty",
		"Trusardi",
		"Avon духи женские",
		"Love is",
		"Духи женские avon",
		"шиммер для тела",
		"косметика для подростков",
		"trussardi",
		"lacoste",
		"lacoste духи",
		"lacoste женские",
		"lacoste мужской",
		"эйвон духи",
		"табак",
		"silver",
		"aqua",
		"духи наркотик женские",
	}
)

func (p *ProductPositionResponse) fetchKeywords() (err error) {
	kw := MPStatsKeywords{}
	err = kw.GetJSON(p.Id)
	if err != nil {
		return errors.New(fmt.Sprintf("error getting product keywords: %s", err))
	}

	wordsData := make([]WordsData, 0, len(kw.Words))
	for k, v := range kw.Words {
		v.Word = k
		wordsData = append(wordsData, v)
	}

	p.WordsData = wordsData

	return nil
}

func (ppl *ProductPositionListResponse) GetJSON(req ProductListRequest) (err error) {
	if req.SellerID == "" {
		return fmt.Errorf("invalid id: null")
	}

	if req.Mode == "" {
		req.Mode = ModeSeller
	}

	kReq := ProductListRequest{
		SellerID: req.SellerID,
		Mode:     req.Mode,
		Xsubject: req.Xsubject,
		Page:     1,
		Limit:    10,
	}

	if kReq.Mode == ModeCategory {
		if req.CategoryID == "" {
			return errors.New(fmt.Sprintf("invalid category id: null"))
		}

		kReq.CategoryID = req.CategoryID
		kReq.Limit = 25
	}

	fmt.Printf("Request: [%v]\n", kReq)
	plResponse := ProductListResponse{}
	err = plResponse.GetJSON(kReq)
	if err != nil {
		return errors.New(fmt.Sprintf("error getting product top list: %s", err))
	}

	topList := plResponse.Data.Products

	fmt.Printf("GetJSON: [1.s]\n")
	commonWords := FindCommonKeywords(topList)
	fmt.Printf("GetJSON: [1.e]\n")

	fmt.Printf("GetJSON: [2.s]\n")
	ppl.Keywords = []string{}
	for _, word := range commonWords {
		ppl.Keywords = append(ppl.Keywords, word.Word)
	}
	fmt.Printf("GetJSON: [2.e] kw: %s\n", ppl.Keywords)

	if len(ppl.Keywords) == 0 {
		return errors.New(fmt.Sprintf("error getting common keywords: %s", err))
	}

	fmt.Printf("Common keywords: [%s]\n", strings.Join(ppl.Keywords, ", "))

	pReq := ProductListRequest{
		SellerID: req.SellerID,
		Mode:     ModeSeller,
		Xsubject: req.Xsubject,
		Limit:    req.Limit,
	}

	plResponse = ProductListResponse{}
	fmt.Printf("Request: [%v]\n", pReq)
	err = plResponse.GetAll(pReq)
	if err != nil {
		return errors.New(fmt.Sprintf("error getting product list: %s", err))
	}

	ppl.Products = []ProductPositionResponse{}
	for j, prod := range plResponse.Data.Products {
		p := ProductPositionResponse{
			Id:         prod.Id,
			Name:       prod.Name,
			Brand:      prod.Brand,
			BrandId:    prod.BrandId,
			SupplierId: prod.SupplierId,
			Pos:        j + 1,
		}

		//workerPool <- struct{}{}
		//go func(j int, prod ProductPositionResponse) {
		//	defer func() {
		//		// Signal the worker is done
		//		<-workerPool
		//	}()
		//
		//	fmt.Printf("Prod: [%d][%s]\n", j, p.Name)
		//
		//	// Fetch keywords and process
		//	err := p.fetchKeywords()
		//	if err != nil {
		//		fmt.Printf("Error: %s\n", err)
		//		return
		//	}
		//
		//	// Append the processed product to the results
		//	ppl.Products = append(ppl.Products, p)
		//
		//	pLen := len(ppl.Products)
		//	fmt.Printf("Len: [%d]\n", pLen)
		//}(j, p)
		//
		//// Introduce a delay between requests
		//time.Sleep(200 * time.Millisecond)

		kw := MPStatsKeywords{}
		err = kw.GetJSON(prod.Id)
		if err != nil {
			return errors.New(fmt.Sprintf("error getting product keywords: %s", err))
		}
		fmt.Printf("Prod: [%d][%s]\n", j, prod.Name)

		p.WordsData = []WordsData{}
		for k, v := range kw.Words {
			v.Word = k
			p.WordsData = append(p.WordsData, v)
		}

		ppl.Products = append(ppl.Products, p)

		pLen := len(ppl.Products)
		fmt.Printf("Len: [%d]\n", pLen)
	}

	ppl.FilterProductsByCommonKeywords()

	return nil
}

func (ppl *ProductPositionListResponse) FilterProductsByCommonKeywords() {
	for i, product := range ppl.Products {
		fmt.Printf("Product: [%d][%s]\n", i, product.Name)
		ppl.Products[i].WordsData = FilterWordsByCommonKeywords(product.WordsData, ppl.Keywords)
		fmt.Println("--------------------------------------------------")
	}
}

func FindCommonKeywords(products []ProductLR) []CommonWord {
	wordFreq := make(map[string]map[string]WordsData)

	for j, prod := range products {
		kw := MPStatsKeywords{}
		err := kw.GetJSON(prod.Id)
		if err != nil {
			continue
		}

		fmt.Printf("Prod: [%d][%s]\n", j, prod.Name)
		key := fmt.Sprintf("key-%d", prod.Id)

		// Convert map values to a slice
		var wordsSlice []WordsData
		for _, wordsData := range kw.Words {
			wordsSlice = append(wordsSlice, wordsData)
		}

		// Sort each product's WordsData by WbCount descending
		sort.SliceStable(wordsSlice, func(i, j int) bool {
			return wordsSlice[i].WbCount > wordsSlice[j].WbCount
		})

		for count, data := range wordsSlice {
			if _, exists := wordFreq[data.Word]; !exists {
				if utils.ContainsString(brandBlackList, data.Word) {
					continue
				}
				if utils.Contains(blackListedWords, data.Word) {
					continue
				}
				wordFreq[data.Word] = make(map[string]WordsData)
			}
			wordFreq[data.Word][key] = data
			if count == wordsPerProduct {
				break
			}
		}
	}

	var commonWords = make(map[string]*CommonWord)
	maxOccurrences := 0
	// Populate the wordFreq map
	for word, productsData := range wordFreq {
		for _, data := range productsData {
			if _, exists := commonWords[word]; !exists {
				commonWords[word] = &CommonWord{
					Word:    word,
					WbCount: data.WbCount,
					Count:   1,
				}
			} else {
				commonWords[word].WbCount = utils.MaxInt(commonWords[word].WbCount, data.WbCount)
				commonWords[word].Count++
			}

			if commonWords[word].Count > maxOccurrences {
				maxOccurrences = commonWords[word].Count
			}
		}
	}

	pLen := len(commonWords)

	fmt.Printf("FindCommonKeywords: [4.s] pLen: %d\n", pLen)
	// Sort common words by WbCount descending
	var tempWordList []CommonWord
	for _, word := range commonWords {
		tempWordList = append(tempWordList, *word)
	}

	sort.SliceStable(tempWordList, func(i, j int) bool {
		return tempWordList[i].WbCount > tempWordList[j].WbCount
	})

	var commonWordList []CommonWord
	var count int
	for _, word := range tempWordList {
		isCommon := false
		if word.Count == pLen || word.Count == pLen-1 {
			isCommon = true
		}

		if maxOccurrences < pLen-1 && count < maxKeywords {
			isCommon = true
			count++
		}

		if isCommon {
			commonWordList = append(commonWordList, word)
		}
	}

	fmt.Printf("FindCommonKeywords: [4.s] commonWordList: %d\n", len(commonWordList))
	// Assign positions to the common words
	for i, _ := range commonWordList {
		commonWordList[i].Pos = i + 1
	}
	fmt.Printf("FindCommonKeywords: [5.e] %#v\n", commonWordList[:5])

	size := utils.MinInt(maxKeywords, len(commonWordList))

	return commonWordList[:size]
}

//func FindCommonKeywords(list []ProductLR) []string {
//	var kwList []MPStatsKeywords
//
//	for j, prod := range list {
//		kw := MPStatsKeywords{}
//		err := kw.GetJSON(prod.Id)
//		if err != nil {
//			continue
//		}
//
//		fmt.Printf("Prod: [%d][%s]\n", j, prod.Name)
//		kwList = append(kwList, kw)
//
//		pLen := len(kwList)
//		fmt.Printf("Len: [%d]\n", pLen)
//	}
//
//	// Create a map to keep track of k occurrences
//	occurrences := make(map[string]int)
//	positions := make(map[string]int)
//
//	maxOccurrences := 0
//	// Count k occurrences across products
//	for _, item := range kwList {
//		keys := item.SortByWbCount()
//		for _, k := range keys {
//			occurrences[k]++
//
//			positions[k] = utils.MaxInt(positions[k], item.Words[k].WbCount)
//
//			if occurrences[k] > maxOccurrences {
//				maxOccurrences = occurrences[k]
//			}
//		}
//	}
//
//	// Find keywords that occur in all products
//	//pLen := 10
//	pLen := len(list)
//	fmt.Printf("Product len: %d | maxOccurrences count: %d\n", pLen, maxOccurrences)
//	j := 0
//	for k, count := range occurrences {
//		added := false
//		if count == pLen || count == pLen-1 {
//			added = true
//		}
//
//		if maxOccurrences < pLen-1 && j < 10 {
//			added = true
//			j++
//		}
//
//		if !added {
//			delete(positions, k)
//		}
//	}
//
//	sorted := mapslice.New(positions)
//	sorted.SortBy(func(a, b mapslice.Entry[int]) bool {
//		return a.Value > b.Value
//	})
//
//	fmt.Printf("Sorted: %v\n", sorted.Entries())
//
//	size := utils.MinInt(maxKeywords, sorted.Len())
//
//	return sorted.Keys()[:size]
//}

func FilterWordsByCommonKeywords(list []WordsData, keywords []string) (words []WordsData) {
	sort.Slice(list, func(i, j int) bool {
		return list[i].WbCount > list[j].WbCount
	})

	for _, item := range list {
		if utils.Contains(keywords, item.Word) {
			words = append(words, item)
			fmt.Printf("Keyword: %s[wb-count: %d]\n", item.Word, item.WbCount)
		}
		if len(words) == maxKeywords {
			break
		}
	}

	return words
}
