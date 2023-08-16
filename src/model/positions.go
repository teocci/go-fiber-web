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

	plr := ProductListResponse{}
	err = plr.GetJSON(req)
	if err != nil {
		return errors.New(fmt.Sprintf("error getting product list: %s", err))
	}

	ppl.Products = make([]ProductPositionResponse, len(plr.Data.Products))

	list := plr.Data.Products

	fmt.Printf("ProductListResponse: %+v\n", plr)
	for j, prod := range list {
		ppl.Products[j] = ProductPositionResponse{
			Id:         prod.Id,
			Name:       prod.Name,
			Brand:      prod.Brand,
			BrandId:    prod.BrandId,
			SupplierId: prod.SupplierId,
			Keywords:   MPStatsKeywords{},
		}

		err = ppl.Products[j].Keywords.GetJSON(prod.Id)
		if err != nil {
			return errors.New(fmt.Sprintf("error getting product keywords: %s", err))
		}
		if j > 10 {
			break
		}
	}

	keywords := FindCommonKeywords(ppl.Products)
	for i, product := range ppl.Products {
		ppl.Products[i].Keywords = FilterWordsByCommonKeywords(product.Keywords, keywords)
	}

	return nil
}

func FilterWordsByCommonKeywords(stat MPStatsKeywords, keywords []string) MPStatsKeywords {
	words := make(map[string]WordsData)

	for k, data := range stat.Words {
		if contains(keywords, k) {
			words[k] = data
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

	// Count keyword occurrences across products
	for _, product := range products {
		for keyword := range product.Keywords.Words {
			occurrences[keyword]++
		}
	}

	// Find keywords that occur in all products
	maxCount := 10
	keywords := make([]string, 0)
	for keyword, count := range occurrences {
		if count == maxCount {
			keywords = append(keywords, keyword)
		}
	}

	return keywords
}
