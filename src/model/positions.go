// Package model
// Created by Teocci.
// Author: teocci@yandex.com on 2023-Aug-16
package model

import (
	"errors"
	"fmt"
)

type ProductPositionResponse struct {
	Id         int             `json:"id"`
	Name       string          `json:"name"`
	Brand      string          `json:"brand"`
	BrandId    int             `json:"brandId"`
	SupplierId int             `json:"supplierId"`
	Keywords   MPStatsKeywords `json:"keywords,omitempty"`
}

type ProductPositionListResponse struct {
	Products []ProductPositionResponse `json:"products"`
}

func (ppl *ProductPositionListResponse) GetJSON(req ProductListRequest) (err error) {
	if req.ID == "" {
		return fmt.Errorf("invalid id: null")
	}

	if req.Mode == "" {
		req.Mode = ModeSeller
	}

	req.Limit = 10

	plResponse := ProductListResponse{}
	err = plResponse.GetJSON(req)
	if err != nil {
		return errors.New(fmt.Sprintf("error getting product topList: %s", err))
	}

	ppl.Products = []ProductPositionResponse{}

	topList := plResponse.Data.Products

	fmt.Printf("ProductListResponse: %+v\n", plResponse)
	for j, prod := range topList {
		p := ProductPositionResponse{
			Id:         prod.Id,
			Name:       prod.Name,
			Brand:      prod.Brand,
			BrandId:    prod.BrandId,
			SupplierId: prod.SupplierId,
			Keywords:   MPStatsKeywords{},
		}
		err = p.Keywords.GetJSON(prod.Id)
		if err != nil {
			return errors.New(fmt.Sprintf("error getting product keywords: %s", err))
		}
		fmt.Printf("Prod: [%d][%s]\n", j, prod.Name)
		ppl.Products = append(ppl.Products, p)

		pLen := len(ppl.Products)
		fmt.Printf("Len: [%d]\n", pLen)
	}

	keywords := FindCommonKeywords(ppl.Products)
	fmt.Printf("Common keywords: %+v\n", keywords)
	for i, product := range ppl.Products {
		ppl.Products[i].Keywords = FilterWordsByCommonKeywords(product.Keywords, keywords)
	}

	return nil
}

func FilterWordsByCommonKeywords(stat MPStatsKeywords, keywords []string) MPStatsKeywords {
	words := make(map[string]WordsData)
	sorted := stat.SortByWbCount()

	for _, k := range sorted {
		if contains(keywords, k) {
			words[k] = stat.Words[k]
			//fmt.Printf("Keyword: %s[wb-count: %d]\n", k, data.WbCount)
		}
		if len(words) == 15 {
			break
		}
	}

	stat.Words = words

	return stat
}

func contains(slice []string, element string) bool {
	for _, item := range slice {
		if item == element {
			return true
		}
	}
	return false
}

func FindCommonKeywords(products []ProductPositionResponse) []string {
	// Create a map to keep track of keyword occurrences
	occurrences := make(map[string]int)

	max := 0
	// Count keyword occurrences across products
	for _, product := range products {
		for keyword := range product.Keywords.Words {
			occurrences[keyword]++
			if occurrences[keyword] > max {
				max = occurrences[keyword]
			}
		}
	}

	// Find keywords that occur in all products
	//pLen := 10
	pLen := len(products)
	fmt.Printf("Product len: %d | max count: %d\n", pLen, max)
	keywords := make([]string, 0)
	for keyword, count := range occurrences {
		if count == pLen || count == max {
			keywords = append(keywords, keyword)
		}
	}

	return keywords
}
