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
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid action: null",
		})
	}
	productID := c.Params("id")
	if productID == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid seller id: null",
		})
	}

	products := model.IdenticalProductsResponse{}
	err := products.GetJSON(productID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": products,
	})
}
