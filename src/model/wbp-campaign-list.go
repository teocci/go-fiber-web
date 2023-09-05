// Package model
// Created by Teocci.
// Author: teocci@yandex.com on 2023-Sep-05
package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/teocci/go-fiber-web/src/utils"
	"net/http"
	"net/url"
	"time"
)

type WBCampaignStatus int

const (
	WBCampaignStatusReady   WBCampaignStatus = 4
	WBCampaignStatusOver                     = 7
	WBCampaignStatusRefused                  = 8
	WBCampaignStatusActive                   = 9
	WBCampaignStatusPaused                   = 11
)

type WBCampaignType int

const (
	WBCampaignTypeCatalog       WBCampaignType = 4
	WBCampaignTypeContent                      = 5
	WBCampaignTypeSearch                       = 6
	WBCampaignTypeRecommended                  = 7
	WBCampaignTypeAuto                         = 8
	WBCampaignTypeSearchCatalog                = 9
)

type WBCampaignOrder string

const (
	WBCampaignOrderCreate WBCampaignOrder = "create"
	WBCampaignOrderChange                 = "change"
	WBCampaignOrderId                     = "id"
)

type WBCampaignDirection string

const (
	WBCampaignDirectionAsc  WBCampaignDirection = "asc"
	WBCampaignDirectionDesc                     = "desc"
)

type WBPCampaignInfo struct {
	AdvertId        int       `json:"advertId"`
	Name            string    `json:"name"`
	Type            int       `json:"type"`
	Status          int       `json:"status"`
	DailyBudget     int       `json:"dailyBudget,omitempty"`
	CreateTime      time.Time `json:"createTime,omitempty"`
	ChangeTime      time.Time `json:"changeTime,omitempty"`
	StartTime       time.Time `json:"startTime,omitempty"`
	EndTime         time.Time `json:"endTime,omitempty"`
	SearchPlusState bool      `json:"searchPluseState,omitempty"`
}

type WBPCampaignListResponse struct {
	Content []WBPCampaignInfo
}

type WBPCampaignListRequest struct {
	SellerID  string              `json:"seller_id"`
	Status    WBCampaignStatus    `json:"status"`
	Type      WBCampaignType      `json:"type"`
	Limit     int                 `json:"limit"`
	Offset    int                 `json:"offset"`
	Order     WBCampaignOrder     `json:"order"`
	Direction WBCampaignDirection `json:"direction"`
}

func IsValidWBCampaignStatus(value WBCampaignStatus) bool {
	switch value {
	case WBCampaignStatusReady,
		WBCampaignStatusOver,
		WBCampaignStatusRefused,
		WBCampaignStatusActive,
		WBCampaignStatusPaused:
		return true
	default:
		return false
	}
}

func IsValidCampaignType(value WBCampaignType) bool {
	switch value {
	case WBCampaignTypeCatalog,
		WBCampaignTypeContent,
		WBCampaignTypeSearch,
		WBCampaignTypeRecommended,
		WBCampaignTypeAuto,
		WBCampaignTypeSearchCatalog:
		return true
	default:
		return false
	}
}

func IsValidCampaignOrder(value WBCampaignOrder) bool {
	switch value {
	case WBCampaignOrderCreate,
		WBCampaignOrderChange,
		WBCampaignOrderId:
		return true
	default:
		return false
	}
}

func IsValidCampaignDirection(value WBCampaignDirection) bool {
	switch value {
	case WBCampaignDirectionAsc,
		WBCampaignDirectionDesc:
		return true
	default:
		return false
	}
}

func (wblr *WBPCampaignListRequest) GetToken() (token string, found bool) {
	if wblr.SellerID == "" {
		return
	}

	token = parseWBToken(wblr.SellerID, "ADV")
	if token == "" {
		return
	}

	return token, true
}

func (alr *WBPCampaignListResponse) GetAllJSON(req WBPCampaignListRequest) (err error) {
	token, found := req.GetToken()
	if !found {
		return errors.New("token not found")
	}

	if req.Status == 0 {
		return errors.New("status is required")
	}

	if !IsValidWBCampaignStatus(req.Status) {
		return errors.New("invalid status")
	}

	if req.Type != 0 && !IsValidCampaignType(req.Type) {
		return errors.New("invalid type")
	}

	order := string(req.Order)
	if order == "" {
		order = WBCampaignOrderChange
	}

	direction := string(req.Direction)
	if direction == "" {
		direction = WBCampaignDirectionDesc
	}

	baseURL := &url.URL{
		Scheme: "https",
		Host:   "advert-api.wb.ru",
		Path:   "/adv/v0/adverts",
	}
	params := baseURL.Query()
	params.Set("status", utils.AnyToString(req.Status))

	if req.Type > 0 {
		params.Set("type", utils.AnyToString(req.Type))
	}

	if req.Limit > 0 {
		params.Set("limit", utils.AnyToString(req.Limit))
	}

	if req.Offset > 0 {
		params.Set("offset", utils.AnyToString(req.Offset))
	}

	params.Set("order", order)
	params.Set("direction", direction)

	baseURL.RawQuery = params.Encode()

	apiURL := baseURL.String()

	headers := map[string]string{
		"Authorization": token,
	}

	fmt.Printf("WB_API_URL: %s\n", apiURL)
	r, err := utils.GetWithHeaders(apiURL, headers)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode != http.StatusNoContent && r.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("error: %d", r.StatusCode))
	}

	if r.StatusCode == http.StatusNoContent {
		return
	}

	var data []WBPCampaignInfo
	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(&data)
	if err != nil {
		fmt.Printf("WBPCampaignListResponse: %#v\n", err)
		return err
	}

	fmt.Printf("WBPCampaignListResponse: %#v\n", data)

	alr.Content = data

	return
}
