// Package model
// Created by Teocci.
// Author: teocci@yandex.com on 2023-Sep-03
package model

import (
	"errors"
	"fmt"
	"sync"
)

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

type CampaignRAPListRequest struct {
	SellerID  string              `json:"seller_id"`
	Type      WBCampaignType      `json:"type"`
	Limit     int                 `json:"limit"`
	Offset    int                 `json:"offset"`
	Order     WBCampaignOrder     `json:"order"`
	Direction WBCampaignDirection `json:"direction"`
}

func (r *CampaignListResponse) GetRAPAdsList(req CampaignRAPListRequest) (err error) {
	if req == (CampaignRAPListRequest{}) {
		return errors.New("invalid request")
	}

	if req.SellerID == "" {
		return errors.New("seller id is required")
	}

	fmt.Printf("GetRAPAdsList Seller: %s\n", req.SellerID)

	statuses := []WBCampaignStatus{
		WBCampaignStatusReady,
		WBCampaignStatusActive,
		WBCampaignStatusPaused,
	}

	// Create a channel to collect results from goroutines
	resultChan := make(chan []CampaignInfo)

	// Create a WaitGroup to wait for all goroutines to finish
	var wg sync.WaitGroup

	for _, status := range statuses {
		// Increment the WaitGroup counter
		wg.Add(1)

		go func(status WBCampaignStatus) {
			defer wg.Done()

			fmt.Printf("Go Routine status: %d\n", status)

			listReq := WBPCampaignListRequest{
				SellerID:  req.SellerID,
				Status:    status,
				Type:      req.Type,
				Limit:     req.Limit,
				Offset:    req.Offset,
				Order:     req.Order,
				Direction: req.Direction,
			}

			list := WBPCampaignListResponse{}
			if err := list.GetAllJSON(listReq); err != nil {
				return
			}

			// Collect the campaign info from this goroutine
			var tmp []CampaignInfo
			for i, v := range list.Content {
				tmp = append(tmp, CampaignInfo{
					Pos:    i + 1,
					Id:     v.AdvertId,
					Name:   v.Name,
					Type:   v.Type,
					Status: v.Status,
				})
			}

			fmt.Printf("GetRAPAdsList tmp[%d]: %#v\n", len(tmp), tmp)

			// Send the collected campaign info to the result channel
			resultChan <- tmp
		}(status)
	}

	// Start a goroutine to close the result channel when all goroutines are done
	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Collect results from the channel
	for campaignInfo := range resultChan {
		r.List = append(r.List, campaignInfo...)
	}

	return
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
