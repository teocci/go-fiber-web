// Package model
// Created by Teocci.
// Author: teocci@yandex.com on 2023-Sep-03
package model

type AdsInfo struct {
	Pos        int    `json:"pos"`
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Type       int    `json:"type"`
	Status     int    `json:"status"`
	CampaignId int    `json:"campaign_id"`
}

type AdsListResponse struct {
	List []AdsInfo `json:"list"`
}

type AdsListRequest struct {
	Page  int `json:"page"`
	Limit int `json:"limit"`
}

func (r *AdsListResponse) GetAdsList(req AdsListRequest) (err error) {
	listReq := AdvertisementListRequest{}
	listReq.Page = req.Page
	listReq.Limit = req.Limit

	list := AdvertisementListResponse{}
	err = list.GetAllJSON(listReq)
	if err != nil {
		return
	}

	for i, v := range list.Content {
		r.List = append(r.List, AdsInfo{
			Pos:        i + 1,
			Id:         v.Id,
			Name:       v.CampaignName,
			CampaignId: v.CampaignId,
			Type:       v.Type,
			Status:     v.StatusId,
		})
	}

	return
}
