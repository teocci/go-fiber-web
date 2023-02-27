// Package model
// Created by RTT.
// Author: teocci@yandex.com on 2023-2월-27
package model

type ProductResponse struct {
	State   int `json:"state"`
	Version int `json:"version"`
	Params  struct {
		Curr    string `json:"curr"`
		Spp     int    `json:"spp"`
		Version int    `json:"version"`
	} `json:"params"`
	Data struct {
		Products []struct {
			Sort            int           `json:"__sort"`
			Ksort           int           `json:"ksort"`
			Time1           int           `json:"time1"`
			Time2           int           `json:"time2"`
			Dist            int           `json:"dist"`
			Id              int           `json:"id"`
			Root            int           `json:"root"`
			KindId          int           `json:"kindId"`
			SubjectId       int           `json:"subjectId"`
			SubjectParentId int           `json:"subjectParentId"`
			Name            string        `json:"name"`
			Brand           string        `json:"brand"`
			BrandId         int           `json:"brandId"`
			SiteBrandId     int           `json:"siteBrandId"`
			SupplierId      int           `json:"supplierId"`
			Sale            int           `json:"sale"`
			PriceU          int           `json:"priceU"`
			SalePriceU      int           `json:"salePriceU"`
			LogisticsCost   int           `json:"logisticsCost"`
			SaleConditions  int           `json:"saleConditions"`
			Pics            int           `json:"pics"`
			Rating          int           `json:"rating"`
			Feedbacks       int           `json:"feedbacks"`
			Volume          int           `json:"volume"`
			Colors          []interface{} `json:"colors"`
			Sizes           []struct {
				Name     string `json:"name"`
				OrigName string `json:"origName"`
				Rank     int    `json:"rank"`
				OptionId int    `json:"optionId"`
				Wh       int    `json:"wh"`
				Sign     string `json:"sign"`
			} `json:"sizes"`
			DiffPrice    bool   `json:"diffPrice"`
			PanelPromoId int    `json:"panelPromoId,omitempty"`
			PromoTextCat string `json:"promoTextCat,omitempty"`
		} `json:"products"`
	} `json:"data"`
}
