// Package model
// Created by Teocci.
// Author: teocci@yandex.com on 2023-Sep-03
package model

import (
	"encoding/json"
	"fmt"
	"github.com/teocci/go-fiber-web/src/utils"
	"net/http"
	"net/url"
	"os"
	"time"
)

type WBCampaignContent struct {
	Id                 int           `json:"id"`
	Type               int           `json:"type"`
	StatusId           int           `json:"statusId"`
	ImagePath          string        `json:"imagePath"`
	CategoryUid        string        `json:"categoryUid"`
	CategoryName       string        `json:"categoryName"`
	CampaignId         int           `json:"campaignId"`
	CampaignName       string        `json:"campaignName"`
	BrandName          string        `json:"brandName"`
	StartDate          time.Time     `json:"startDate"`
	FinishDate         *time.Time    `json:"finishDate"`
	Position           interface{}   `json:"position"`
	PredictedStartDate interface{}   `json:"predictedStartDate"`
	PredictedEndDate   interface{}   `json:"predictedEndDate"`
	CarouselCard       []interface{} `json:"carouselCard"`
	Products           []struct {
		Nm      int    `json:"nm"`
		Name    string `json:"name"`
		Subject struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"subject"`
		Brand struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"brand"`
		Kind struct {
			Id   int    `json:"id"`
			Name string `json:"name"`
		} `json:"kind"`
		Categories []interface{} `json:"categories"`
	} `json:"products"`
	TotalCount  int         `json:"TotalCount"`
	CreateDate  time.Time   `json:"createDate"`
	ChangeDate  time.Time   `json:"ChangeDate"`
	UrlType     int         `json:"UrlType"`
	Url         string      `json:"Url"`
	DisableDate *time.Time  `json:"disableDate"`
	PromoSetId  interface{} `json:"PromoSetId"`
}

type WBCampaignListResponse struct {
	HttpStatus int    `json:"httpStatus"`
	Error      string `json:"error"`
	Code       int    `json:"code"`
	Counts     struct {
		PageCount    int `json:"pageCount"`
		TotalCount   int `json:"totalCount"`
		ActiveCount  int `json:"activeCount"`
		PauseCount   int `json:"pauseCount"`
		DraftCount   int `json:"draftCount"`
		ArchiveCount int `json:"archiveCount"`
	} `json:"counts"`
	Content []WBCampaignContent `json:"content"`
}

type WBCampaignListRequest struct {
	SellerID string `json:"seller_id"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

func (alr *WBCampaignListResponse) GetAllJSON(req WBCampaignListRequest) (err error) {
	if req.Page == 0 {
		req.Page = 1
	}

	if req.Limit == 0 {
		req.Limit = 100
	}

	page := utils.AnyToString(req.Page)
	size := utils.AnyToString(req.Limit)

	baseURL := &url.URL{
		Scheme: "https",
		Host:   "cmp.wildberries.ru",
		Path:   "/backend/api/v3/atrevds",
	}

	params := baseURL.Query()
	params.Set("pageNumber", page)
	params.Set("pageSize", size)
	params.Set("search", "")
	params.Set("status", "%5B11,4,9%5D")
	params.Set("order", "createDate")
	params.Set("direction", "desc")
	params.Set("type", "%5B2,3,4,5,6,7,8,9%5D")
	baseURL.RawQuery = params.Encode()

	apiURL := baseURL.String()
	fmt.Printf("WB_API_URL: %#v\n", apiURL)

	cookies := []*http.Cookie{
		{
			Name:   "x-supplier-id-external",
			Value:  os.Getenv("X_SUPPLIER_ID_EXTERNAL"),
			Domain: ".wildberries.ru",
			Secure: true,
		},
		{
			Name:     "WBToken",
			Value:    os.Getenv("WB_TOKEN"),
			Domain:   "cmp.wildberries.ru",
			HttpOnly: true,
			Secure:   true,
		},
	}

	headers := map[string]string{
		"User-Agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/116.0.0.0 Safari/537.36",
		"X-User-Id":  os.Getenv("X_USER_ID"),
	}

	r, err := utils.GetWithCookiesAndHeaders(apiURL, cookies, headers)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	err = json.NewDecoder(r.Body).Decode(&alr)

	return err
}
