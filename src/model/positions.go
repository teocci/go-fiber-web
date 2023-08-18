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
	Keywords   []WordsData `json:"keywords,omitempty"`
}

type ProductPositionListResponse struct {
	Products []ProductPositionResponse `json:"products"`
}

type infoData struct {
	Product ProductPositionResponse
	Stats   MPStatsKeywords
}

func (ppl *ProductPositionListResponse) GetJSON(req ProductListRequest) (err error) {
	if req.ID == "" {
		return fmt.Errorf("invalid id: null")
	}

	if req.Mode == "" {
		req.Mode = ModeSeller
	}

	if req.Limit < 2 {
		req.Limit = 10
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
		return errors.New(fmt.Sprintf("error getting product topList: %s", err))
	}

	topList := plResponse.Data.Products

	keywords := FindCommonKeywords(topList)
	if len(keywords) == 0 {
		return errors.New(fmt.Sprintf("error getting common keywords: %s", err))
	}

	fmt.Printf("Common keywords: [%s]\n", strings.Join(keywords, ", "))

	ppl.Products = []ProductPositionResponse{}
	for j, prod := range topList {
		p := ProductPositionResponse{
			Id:         prod.Id,
			Name:       prod.Name,
			Brand:      prod.Brand,
			BrandId:    prod.BrandId,
			SupplierId: prod.SupplierId,
		}

		kw := MPStatsKeywords{}
		err = kw.GetJSON(prod.Id)
		if err != nil {
			return errors.New(fmt.Sprintf("error getting product keywords: %s", err))
		}
		fmt.Printf("Prod: [%d][%s]\n", j, prod.Name)

		p.Keywords = []WordsData{}
		for k, v := range kw.Words {
			v.Word = k
			p.Keywords = append(p.Keywords, v)
		}

		ppl.Products = append(ppl.Products, p)

		pLen := len(ppl.Products)
		fmt.Printf("Len: [%d]\n", pLen)
	}

	for i, product := range ppl.Products {
		fmt.Printf("Product: [%d][%s]\n", i, product.Name)
		ppl.Products[i].Keywords = FilterWordsByCommonKeywords(product.Keywords, keywords)
		fmt.Println("--------------------------------------------------")
	}

	return nil
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

	counter := 0
	for i := range occurrences {
		fmt.Printf("occurrences - [keyword: %s][occurrences: %d]\n", i, occurrences[i])
		if counter >= 10 {
			break
		}
		counter++
	}

	counter = 0
	for i := range occurrences {
		fmt.Printf("positions - [keyword: %s][positions: %d]\n", i, positions[i])
		if counter >= 10 {
			break
		}
		counter++
	}
	// Find keywords that occur in all products
	//pLen := 10
	pLen := len(list)
	fmt.Printf("Product len: %d | maxOccurrences count: %d\n", pLen, maxOccurrences)
	j := 0
	for k, count := range occurrences {
		added := false
		if count == pLen {
			added = true
		}

		if count == maxOccurrences && j < 10 {
			j++
			added = true
		}

		if !added {
			delete(positions, k)
		}
	}

	sorted := mapslice.New(positions)
	sorted.SortBy(func(a, b mapslice.Entry[int]) bool {
		return a.Value > b.Value
	})

	return sorted.Keys()
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
