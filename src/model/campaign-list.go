// Package model
// Created by Teocci.
// Author: teocci@yandex.com on 2023-Sep-03
package model

type CampaignInfo struct {
	Pos        int    `json:"pos"`
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Type       int    `json:"type"`
	Status     int    `json:"status"`
	CampaignId int    `json:"campaign_id"`
}

type CampaignListResponse struct {
	List []CampaignInfo `json:"list"`
}

type CampaignListRequest struct {
	SellerID string `json:"seller_id"`
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
}

func (r *CampaignListResponse) GetAdsList(req CampaignListRequest) (err error) {
	listReq := WBCampaignListRequest{}
	listReq.Page = req.Page
	listReq.Limit = req.Limit

	list := WBCampaignListResponse{}
	err = list.GetAllJSON(listReq)
	if err != nil {
		return
	}

	for i, v := range list.Content {
		r.List = append(r.List, CampaignInfo{
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
