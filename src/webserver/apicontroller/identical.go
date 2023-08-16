// Package apicontroller
// Created by Teocci.
// Author: teocci@yandex.com on 2023-Feb-27
package apicontroller

import (
	"github.com/gofiber/fiber/v2"

	"github.com/teocci/go-fiber-web/src/model"
)

func HandleIdentical(c *fiber.Ctx) error {
	productID := c.Params("id")
	if productID == "" {
		return renderBadRequest(c, "Invalid product id: null")
	}

	products := model.IdenticalProductsResponse{}
	err := products.GetJSON(productID)
	if err != nil {
		return renderInternalError(c, err.Error())
	}

	return c.JSON(fiber.Map{
		"data": products,
	})
}
