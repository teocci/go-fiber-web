// Package apicontroller
// Created by RTT.
// Author: teocci@yandex.com on 2023-2월-27
package apicontroller

import (
	"github.com/gofiber/fiber/v2"

	"github.com/teocci/go-fiber-web/src/model"
)

func HandleSeller(c *fiber.Ctx) error {
	sellerID := c.Params("id")
	if sellerID == "" {
		return renderBadRequest(c, "Invalid seller id: null")
	}

	seller := model.SellerResponse{}
	err := seller.GetJSON(sellerID)
	if err != nil {
		return renderInternalError(c, err.Error())
	}

	response := response{Data: seller}

	return c.JSON(response)
}
