// Package model
// Created by RTT.
// Author: teocci@yandex.com on 2023-2월-27
package model

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type ProductDetail struct {
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
	IsNew           bool          `json:"isNew,omitempty"`
	Colors          []interface{} `json:"colors"`
	Promotions      []int         `json:"promotions"`
	Sizes           []struct {
		Name     string `json:"name"`
		OrigName string `json:"origName"`
		Rank     int    `json:"rank"`
		OptionId int    `json:"optionId"`
		Stocks   []struct {
			Wh    int `json:"wh"`
			Qty   int `json:"qty"`
			Time1 int `json:"time1"`
			Time2 int `json:"time2"`
		} `json:"stocks"`
		Time1 int    `json:"time1"`
		Time2 int    `json:"time2"`
		Wh    int    `json:"wh"`
		Sign  string `json:"sign"`
	} `json:"sizes"`
	DiffPrice bool `json:"diffPrice"`
	Time1     int  `json:"time1"`
	Time2     int  `json:"time2"`
	Wh        int  `json:"wh"`
	Extended  struct {
		BasicSale   int `json:"basicSale"`
		BasicPriceU int `json:"basicPriceU"`
	} `json:"extended,omitempty"`
}

type ProductDetailResponse struct {
	State  int `json:"state"`
	Params struct {
		Curr    string `json:"curr"`
		Spp     int    `json:"spp"`
		Version int    `json:"version"`
	} `json:"params"`
	Data struct {
		Products []ProductDetail `json:"products"`
	} `json:"data"`
}

func (pr *ProductDetailResponse) GetJSON(nm string) (err error) {
	baseURL := &url.URL{
		Scheme: "https",
		Host:   "catalog.wb.ru",
		Path:   "/sellers/catalog",
	}

	params := baseURL.Query()
	params.Set("appType", "1")
	params.Set("couponsGeo", "12,3,18,15,21")
	params.Set("curr", "rub")
	params.Set("dest", "-1257786")
	params.Set("emp", "0")
	params.Set("lang", "ru")
	params.Set("locale", "ru")
	params.Set("pricemarginCoeff", "1.0")
	params.Set("reg", "0")
	params.Set("regions", "80,64,38,4,83,33,68,70,69,30,86,75,40,1,22,66,31,48,110,71")
	params.Set("spp", "0")
	params.Set("nm", nm)
	baseURL.RawQuery = params.Encode()

	apiURL := baseURL.String()
	fmt.Printf("WB_API_URL: %#v\n", apiURL)

	r, err := httpClient.Get(apiURL)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&pr)

	return err
}
