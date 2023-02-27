// Package apictrlr
// Created by Teocci.
// Author: teocci@yandex.com on 2023-Feb-27
package apictrlr

import (
	"github.com/gofiber/fiber/v2"

	"github.com/teocci/go-fiber-web/src/model"
)

func HandleProducts(c *fiber.Ctx) error {
	sellerID := c.Params("id")
	if sellerID == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid seller id: null",
		})
	}

	products := model.ProductResponse{}
	err := products.GetJSON(sellerID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": products.Data,
	})
}
