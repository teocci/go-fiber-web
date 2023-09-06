// Package apicontroller
// Created by Teocci.
// Author: teocci@yandex.com on 2023-Aug-15
package apicontroller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/teocci/go-fiber-web/src/model"
)

func HandleMarketing(c *fiber.Ctx) error {
	sellerID := c.Params("id")
	if sellerID == "" {
		return renderBadRequest(c, "Invalid id: null")
	}

	action := c.Params("action")
	if action == "" {
		action = "list"
	}

	switch action {
	case "list":
		return handleMarketingList(c, sellerID)
	case "start", "pause", "stop":
		return handleMarketingControl(c, sellerID, action)
	default:
		return renderBadRequest(c, fmt.Sprintf("Invalid action: %s", action))
	}
}

func handleMarketingControl(c *fiber.Ctx, id, action string) error {
	req := model.WBCampaignControlRequest{
		SellerID:   id,
		Action:     action,
		CampaignID: c.Params("campaign"),
	}

	err := req.SendAction()
	if err != nil {
		fmt.Printf("error: %#v", err)
		return renderBadRequest(c, err.Error())
	}

	resp := apiResponse{
		Success: true,
	}

	return c.JSON(resp)
}

func handleMarketingList(c *fiber.Ctx, id string) error {
	req := model.CampaignRAPListRequest{
		SellerID: id,
		Limit:    c.QueryInt("limit", 0),
	}

	list := model.CampaignListResponse{}
	err := list.GetRAPAdsList(req)
	if err != nil {
		return renderBadRequest(c, err.Error())
	}

	resp := apiResponse{Data: list.List}

	return c.JSON(resp)
}
