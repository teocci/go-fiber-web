// Package apicontroller
// Created by Teocci.
// Author: teocci@yandex.com on 2023-Aug-15
package apicontroller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/teocci/go-fiber-web/src/model"
)

func HandleMarketing(c *fiber.Ctx) error {
	req := model.AdsListRequest{
		Page:  c.QueryInt("page", 0),
		Limit: c.QueryInt("limit", 0),
	}

	list := model.AdsListResponse{}
	err := list.GetAdsList(req)
	if err != nil {
		return renderBadRequest(c, err.Error())
	}

	response := response{Data: list.List}

	return c.JSON(response)
}
