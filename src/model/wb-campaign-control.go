// Package model
// Created by Teocci.
// Author: teocci@yandex.com on 2023-Sep-05
package model

import (
	"fmt"
	"github.com/teocci/go-fiber-web/src/utils"
	"net/url"
	"os"
)

type WBCampaignControlRequest struct {
	SellerID   string `json:"seller_id"`
	Action     string `json:"action"`
	CampaignID string `json:"campaign_id"`
}

func (wbcr *WBCampaignControlRequest) SendAction() (err error) {
	baseURL := &url.URL{
		Scheme: "https",
		Host:   "advert-api.wb.ru",
		Path:   fmt.Sprintf("/adv/v0/%s", wbcr.Action),
	}

	q := baseURL.Query()
	q.Set("id", wbcr.CampaignID)
	baseURL.RawQuery = q.Encode()

	apiURL := baseURL.String()
	fmt.Printf("WB_API_URL: %#v\n", apiURL)

	headers := map[string]string{
		"Authorization": os.Getenv(fmt.Sprintf("WB_AUTH_TOKEN_%s", wbcr.SellerID)),
	}

	r, err := utils.GetWithHeaders(apiURL, headers)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	if r.StatusCode == 200 {
		return nil
	}

	return fmt.Errorf("error: %d", r.StatusCode)
}
