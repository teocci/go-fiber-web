// Package apictrlr
// Created by RTT.
// Author: teocci@yandex.com on 2023-2월-27
package apictrlr

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/teocci/go-fiber-web/src/model"
)

const (
	urlFormatSeller = "https://www.wildberries.ru/webapi/seller/data/short/%s"
)

func HandleSeller(c *fiber.Ctx) error {
	sellerID := c.Params("id")
	if sellerID == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid seller id: null",
		})
	}

	url := fmt.Sprintf(urlFormatSeller, sellerID)

	seller := model.SellerResponse{}
	err := seller.GetJSON(url)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": seller,
	})
}
