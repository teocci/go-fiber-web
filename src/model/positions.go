// Package model
// Created by Teocci.
// Author: teocci@yandex.com on 2023-Aug-16
package model

import (
	"errors"
	"fmt"
	"github.com/teocci/go-fiber-web/src/utils"
	"github.com/teocci/go-fiber-web/src/utils/mapslice"
	"sort"
	"strings"
)

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
	if req.ID == "" {
		return fmt.Errorf("invalid id: null")
	}

	if req.Mode == "" {
		req.Mode = ModeSeller
	}

	kReq := ProductListRequest{
		ID:       req.ID,
		Mode:     req.Mode,
		Xsubject: req.Xsubject,
		Page:     1,
		Limit:    10,
	}

	plResponse := ProductListResponse{}
	err = plResponse.GetJSON(kReq)
	if err != nil {
		return errors.New(fmt.Sprintf("error getting product top list: %s", err))
	}

	topList := plResponse.Data.Products

	ppl.Keywords = FindCommonKeywords(topList)
	if len(ppl.Keywords) == 0 {
		return errors.New(fmt.Sprintf("error getting common keywords: %s", err))
	}

	fmt.Printf("Common keywords: [%s]\n", strings.Join(ppl.Keywords, ", "))

	plResponse = ProductListResponse{}
	fmt.Printf("Request: [%v]\n", req)
	err = plResponse.GetAll(req)
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

func FindCommonKeywords(list []ProductLR) []string {
	var kwList []MPStatsKeywords

	for j, prod := range list {
		kw := MPStatsKeywords{}
		err := kw.GetJSON(prod.Id)
		if err != nil {
			continue
		}

		fmt.Printf("Prod: [%d][%s]\n", j, prod.Name)
		kwList = append(kwList, kw)

		pLen := len(kwList)
		fmt.Printf("Len: [%d]\n", pLen)
	}

	// Create a map to keep track of k occurrences
	occurrences := make(map[string]int)
	positions := make(map[string]int)

	maxOccurrences := 0
	// Count k occurrences across products
	for _, item := range kwList {
		keys := item.SortByWbCount()
		for _, k := range keys {
			occurrences[k]++

			positions[k] = utils.MaxInt(positions[k], item.Words[k].WbCount)

			if occurrences[k] > maxOccurrences {
				maxOccurrences = occurrences[k]
			}
		}
	}

	// Find keywords that occur in all products
	//pLen := 10
	pLen := len(list)
	fmt.Printf("Product len: %d | maxOccurrences count: %d\n", pLen, maxOccurrences)
	j := 0
	for k, count := range occurrences {
		added := false
		if count == pLen || count == pLen-1 {
			added = true
		}

		if maxOccurrences < pLen-1 && j < 10 {
			added = true
			j++
		}

		if !added {
			delete(positions, k)
		}
	}

	sorted := mapslice.New(positions)
	sorted.SortBy(func(a, b mapslice.Entry[int]) bool {
		return a.Value > b.Value
	})

	return sorted.Keys()[:10]
}

func FilterWordsByCommonKeywords(list []WordsData, keywords []string) (words []WordsData) {
	sort.Slice(list, func(i, j int) bool {
		return list[i].WbCount > list[j].WbCount
	})

	for _, item := range list {
		if utils.Contains(keywords, item.Word) {
			words = append(words, item)
			fmt.Printf("Keyword: %s[wb-count: %d]\n", item.Word, item.WbCount)
		}
		if len(words) == 15 {
			break
		}
	}

	return words
}
