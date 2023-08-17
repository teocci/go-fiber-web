// Package model
// Created by RTT.
// Author: teocci@yandex.com on 2023-2월-27
package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

const totalPerPage = 100

type ProductLR struct {
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
	DiffPrice    bool            `json:"diffPrice"`
	PanelPromoId int             `json:"panelPromoId,omitempty"`
	PromoTextCat string          `json:"promoTextCat,omitempty"`
	Identical    []ProductDetail `json:"identical,omitempty"`
}

type ProductListResponse struct {
	State   int `json:"state"`
	Version int `json:"version"`
	Params  struct {
		Curr    string `json:"curr"`
		Spp     int    `json:"spp"`
		Version int    `json:"version"`
	} `json:"params"`
	Data struct {
		Products []ProductLR `json:"products"`
	} `json:"data"`
}

type ProductListRequest struct {
	ID       string `json:"id"`
	Mode     string `json:"mode"`
	Xsubject int    `json:"xsubject"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

func (pr *ProductListResponse) GetJSON(req ProductListRequest) (err error) {
	if req.ID == "" {
		return fmt.Errorf("invalid id: null")
	}

	if req.Mode == "" {
		req.Mode = ModeSeller
	}

	var baseURL *url.URL
	switch req.Mode {
	case ModeSeller:
		baseURL = &url.URL{
			Scheme: "https",
			Host:   "catalog.wb.ru",
			Path:   "/sellers/catalog",
		}
	case ModeCategory:
		shard := "beauty3"
		if req.ID == "9000" {
			shard = "beauty4"
		}

		baseURL = &url.URL{
			Scheme: "https",
			Host:   "catalog.wb.ru",
			Path:   fmt.Sprintf("/catalog/%s/catalog", shard),
		}
	}

	params := baseURL.Query()
	params.Set("appType", "1")
	params.Set("curr", "rub")
	params.Set("dest", "-1257786")
	params.Set("regions", "80,64,38,4,83,33,68,70,69,30,86,75,40,1,22,66,31,48,110,71")
	params.Set("sort", "popular")
	params.Set("spp", "0")
	if req.Limit > 1 {
		params.Set("limit", strconv.Itoa(req.Limit))
	}
	if req.Page > 1 {
		params.Set("page", strconv.Itoa(req.Page))
	}
	//params.Set("couponsGeo", "12,3,18,15,21")
	//params.Set("emp", "0")
	//params.Set("lang", "ru")
	//params.Set("locale", "ru")
	//params.Set("pricemarginCoeff", "1.0")
	//params.Set("reg", "0")

	switch req.Mode {
	case ModeSeller:
		params.Set("supplier", req.ID)
	case ModeCategory:
		params.Set("cat", req.ID)
	}

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

func (pr *ProductListResponse) GetFirstPage(req ProductListRequest) (err error) {
	req.Page = 1
	return pr.GetJSON(req)
}

func (pr *ProductListResponse) GetAll(req ProductListRequest) (err error) {
	if req.ID == "" {
		return fmt.Errorf("invalid id: null")
	}

	if req.Mode == "" {
		req.Mode = ModeSeller
	}

	fr := FilterRequest{
		ID:   req.ID,
		Mode: ModeSeller,
	}
	filter := FilterResponse{}
	err = filter.GetJSON(fr)
	if err != nil {
		return err
	}

	useLimit := req.Limit > 1

	req.Page = 1
	totalPages := int(math.Ceil(float64(filter.Data.Total) / float64(totalPerPage)))
	fmt.Printf("%+v, %+v, %+v\n", filter.Data.Total, totalPerPage, totalPages)

	pr.Data.Products = make([]ProductLR, 0)
	raw := make([]ProductLR, 0)

	wg := &sync.WaitGroup{}
	for req.Page <= totalPages {
		wg.Add(1)
		go func(request ProductListRequest) {
			var tmp []ProductLR
			tmp, err = fetchProductListByPage(request)
			raw = append(raw, tmp...)
			wg.Done()
		}(req)
		req.Page++
	}
	wg.Wait()

	for i, p := range raw {
		if p.Id == 120793222 {
			continue
		}

		if useLimit && i > req.Limit {
			break
		}

		identical := IdenticalProductsResponse{}
		err = identical.GetJSON(strconv.Itoa(p.Id))
		if err != nil {
			return err
		}

		if len(identical) == 0 {
			continue
		}

		p.Identical = make([]ProductDetail, 0)

		ids := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(identical)), ";"), "[]")
		fmt.Printf("%d -> [%v]\n", p.Id, identical)

		identicalProducts := ProductDetailResponse{}
		err = identicalProducts.GetJSON(ids)

		products := identicalProducts.Data.Products
		for j, prod := range products {
			seller := SellerResponse{}
			err = seller.GetJSON(strconv.Itoa(prod.SupplierId))
			if err != nil {
				return err
			}

			products[j].SupplierInfo = seller
		}

		p.Identical = append(p.Identical, products...)

		pr.Data.Products = append(pr.Data.Products, p)
	}

	return err
}

func fetchProductListByPage(req ProductListRequest) (raw []ProductLR, err error) {
	tmp := ProductListResponse{}
	err = tmp.GetJSON(req)
	if err != nil {
		return nil, err
	}

	raw = tmp.Data.Products
	length := len(tmp.Data.Products)
	if length == 0 {
		return nil, errors.New("no products")
	}

	return raw, nil
}
