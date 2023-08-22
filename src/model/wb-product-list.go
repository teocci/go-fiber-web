// Package model
// Created by RTT.
// Author: teocci@yandex.com on 2023-2월-27
package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/teocci/go-fiber-web/src/scache"
	"math"
	"net/url"
	"strconv"
	"strings"
	"sync"
)

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
	Xsubject string `json:"xsubject"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

var (
	wbProductListCache = scache.GetCacheInstance[ProductListResponse]("WBProductList")
)

func (pr *ProductListResponse) checkCache(cacheKey string) (ok bool) {
	ok = false
	if cacheKey == "" {
		return
	}

	var cacheData *ProductListResponse
	if cacheData, ok = wbProductListCache.Get(cacheKey); ok {
		fmt.Println("wb-product-list cache hit")
		*pr = *cacheData
		return true
	}

	return false
}

func (pr *ProductListResponse) updateCache(cacheKey string) {
	if cacheKey == "" {
		return
	}
	wbProductListCache.Set(cacheKey, *pr, 0)
}

func (pr *ProductListResponse) GetJSON(req ProductListRequest) (err error) {
	if req.ID == "" {
		return fmt.Errorf("invalid id: null")
	}

	if req.Mode == "" {
		req.Mode = ModeSeller
	}

	cacheKey, baseURL := req.generateCacheKeyAndURL()

	found := pr.checkCache(cacheKey)
	if found {
		return nil
	}

	fmt.Printf("limit: %d\n", req.Limit)

	apiURL := req.generateAPIURL(baseURL)
	fmt.Printf("WB_API_URL: %#v\n", apiURL)

	r, err := httpClient.Get(apiURL)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&pr)

	pr.updateCache(cacheKey)

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

	limit := filter.Data.Total
	if req.Limit > 1 {
		limit = req.Limit
	}

	req.Page = 1
	fullPages := limit / totalPerPage
	remaining := limit % totalPerPage
	totalPages := fullPages
	if remaining > 0 {
		totalPages++
	}
	fmt.Printf("%+v, %+v, %+v\n", limit, totalPerPage, totalPages)

	pr.Data.Products = make([]ProductLR, 0)
	var mu sync.Mutex

	wg := &sync.WaitGroup{}
	for req.Page <= totalPages {
		wg.Add(1)
		if req.Page == totalPages && remaining > 0 {
			req.Limit = remaining
		}
		go func(request ProductListRequest) {
			defer wg.Done()

			tmp := ProductListResponse{}
			err = tmp.GetJSON(request)
			if err != nil {
				fmt.Printf("error getting product list: %s\n", err)
				return
			}
			mu.Lock()
			pr.Data.Products = append(pr.Data.Products, tmp.Data.Products...)
			mu.Unlock()
		}(req)
		req.Page++
	}
	wg.Wait()

	return err
}

func (pr *ProductListResponse) GetIdenticalForAll(req ProductListRequest) (err error) {
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

	raw := make([]ProductLR, 0)

	var mu sync.Mutex

	wg := &sync.WaitGroup{}
	for req.Page <= totalPages {
		wg.Add(1)
		go func(request ProductListRequest) {
			defer wg.Done()
			tmp, err := fetchProductListByPage(request)
			if err != nil {
				fmt.Printf("error getting product list: %s\n", err)
				return
			}

			mu.Lock()
			raw = append(raw, tmp...)
			mu.Unlock()
		}(req)
		req.Page++
	}
	wg.Wait()

	pr.Data.Products = []ProductLR{}

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

		p.Identical = []ProductDetail{}

		ids := strings.Trim(strings.Join(strings.Fields(fmt.Sprint(identical)), ";"), "[]")
		fmt.Printf("%d -> [%v]\n", p.Id, identical)

		identicalProducts := ProductDetailResponse{}
		err = identicalProducts.GetJSON(ids)
		if err != nil {
			return err
		}

		for _, prod := range identicalProducts.Data.Products {
			seller := SellerResponse{}
			err = seller.GetJSON(strconv.Itoa(prod.SupplierId))
			if err != nil {
				return err
			}

			prod.SupplierInfo = seller
			p.Identical = append(p.Identical, prod)
		}

		pr.Data.Products = append(pr.Data.Products, p)
	}

	return err
}

func fetchProductListByPage(req ProductListRequest) ([]ProductLR, error) {
	tmp := ProductListResponse{}
	err := tmp.GetJSON(req)
	if err != nil {
		return nil, err
	}

	if len(tmp.Data.Products) == 0 {
		return nil, errors.New("no products")
	}

	return tmp.Data.Products, nil
}

func (rlReq *ProductListRequest) generateCacheKeyAndURL() (string, *url.URL) {
	var cacheKey string
	var baseURL *url.URL

	baseKey := fmt.Sprintf("wb-product-list-%s-%s", rlReq.ID, rlReq.Mode)
	if rlReq.Xsubject != "" {
		baseKey = fmt.Sprintf("%s-%s", baseKey, rlReq.Xsubject)
	}

	switch rlReq.Mode {
	case ModeSeller:
		baseURL = &url.URL{
			Scheme: "https",
			Host:   "catalog.wb.ru",
			Path:   "/sellers/catalog",
		}
		cacheKey = fmt.Sprintf("%s-%d-%d", baseKey, rlReq.Page, rlReq.Limit)
	case ModeCategory:
		shard := "beauty3"
		if rlReq.ID == "9000" {
			shard = "beauty4"
		}

		baseURL = &url.URL{
			Scheme: "https",
			Host:   "catalog.wb.ru",
			Path:   fmt.Sprintf("/catalog/%s/catalog", shard),
		}
		cacheKey = fmt.Sprintf("%s-%s-%d-%d", baseKey, shard, rlReq.Page, rlReq.Limit)
	}

	return cacheKey, baseURL
}

func (rlReq *ProductListRequest) generateAPIURL(baseURL *url.URL) string {
	params := baseURL.Query()
	params.Set("appType", "1")
	params.Set("curr", "rub")
	params.Set("dest", "-1257786")
	params.Set("regions", "80,64,38,4,83,33,68,70,69,30,86,75,40,1,22,66,31,48,110,71")
	params.Set("sort", "popular")
	params.Set("spp", "0")
	if rlReq.Xsubject != "" {
		params.Set("xsubject", rlReq.Xsubject)
	}
	if rlReq.Limit > 1 {
		params.Set("limit", strconv.Itoa(rlReq.Limit))
	}
	if rlReq.Page > 1 {
		params.Set("page", strconv.Itoa(rlReq.Page))
	}

	switch rlReq.Mode {
	case ModeSeller:
		params.Set("supplier", rlReq.ID)
	case ModeCategory:
		params.Set("cat", rlReq.ID)
	}

	baseURL.RawQuery = params.Encode()
	return baseURL.String()
}
