// Package apicontroller
// Created by Teocci.
// Author: teocci@yandex.com on 2023-Aug-15
package apicontroller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/teocci/go-fiber-web/src/model"
)

const (
	PositionsActionSeller   = "seller"
	PositionsActionCategory = "category"
)

func HandlePositions(c *fiber.Ctx) error {
	req := model.ProductListRequest{}

	req.Mode = c.Params("action")
	if req.Mode == "" {
		return renderBadRequest(c, "Invalid action: null")
	}

	req.ID = c.Params("id")
	if req.ID == "" {
		return renderBadRequest(c, "Invalid id: null")
	}

	req.Xsubject = c.Params("xsubject")

	req.Limit = c.QueryInt("limit", 0)

	list := model.ProductPositionListResponse{}
	err := list.GetJSON(req)
	if err != nil {
		return renderBadRequest(c, err.Error())
	}

	response := response{Data: list}

	return c.JSON(response)
}
