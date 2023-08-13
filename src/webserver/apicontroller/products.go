// Package apicontroller
// Created by Teocci.
// Author: teocci@yandex.com on 2023-Feb-27
package apicontroller

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/teocci/go-fiber-web/src/model"
)

func HandleProductList(c *fiber.Ctx) error {
	supplierID := c.Params("id")
	fmt.Println(string(c.Request().URI().QueryString()))
	limit := c.QueryInt("limit")
	if supplierID == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Invalid seller id: null",
		})
	}

	//if limit == 0 {
	//	limit = 10
	//}

	products := model.ProductListResponse{}
	err := products.GetAll(supplierID, limit)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"data": products.Data,
	})
}
