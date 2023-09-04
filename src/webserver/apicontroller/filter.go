// Package apicontroller
// Created by Teocci.
// Author: teocci@yandex.com on 2023-Aug-11
package apicontroller

import (
	"github.com/gofiber/fiber/v2"

	"github.com/teocci/go-fiber-web/src/model"
)

func HandleFilter(c *fiber.Ctx) error {
	action := c.Params("action")
	if action == "" {
		return renderBadRequest(c, "Invalid action: null")
	}

	id := c.Params("id")
	if id == "" {
		return renderBadRequest(c, "Invalid id: null")
	}

	req := model.FilterRequest{
		ID:   id,
		Mode: action,
	}

	rf := model.FilterResponse{}
	err := rf.GetJSON(req)
	if err != nil {
		return renderBadRequest(c, err.Error())
	}

	response := apiResponse{Data: rf.Data.Filters}

	return c.JSON(response)
}
