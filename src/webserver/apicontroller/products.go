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
		return renderBadRequest(c, "Invalid supplier id: null")
	}

	req := model.ProductListRequest{
		SellerID: supplierID,
		Mode:     model.ModeSeller,
		Limit:    limit,
	}

	products := model.ProductListResponse{}
	err := products.GetIdenticalForAll(req)
	if err != nil {
		return renderInternalError(c, err.Error())
	}

	return c.JSON(fiber.Map{
		"data": products.Data,
	})
}
