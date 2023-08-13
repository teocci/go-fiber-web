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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid seller id: null",
		})
	}

	seller := model.SellerResponse{}
	err := seller.GetJSON(sellerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": seller,
	})
}
