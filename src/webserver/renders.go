// Package webserver
// Created by Teocci.
// Author: teocci@yandex.com on 2023-3월-21
package webserver

import (
	"github.com/gofiber/fiber/v2"
)

func renderPage(c *fiber.Ctx, page PageInfo) error {
	return c.Render("marketing", fiber.Map{
		"page": page,
	})
}

func renderErrorPage(c *fiber.Ctx, code int, msg string) error {
	return c.Status(code).Render("error", fiber.Map{
		"code": code,
		"msg":  msg,
	})
}

func renderInternalError(c *fiber.Ctx) error {
	return renderErrorPage(c, fiber.StatusInternalServerError, "Internal Server Error")
}

func renderBadRequest(c *fiber.Ctx) error {
	return renderErrorPage(c, fiber.StatusBadRequest, "Bad Request")
}

func renderNotFound(c *fiber.Ctx) error {
	return renderErrorPage(c, fiber.StatusNotFound, "Not Found")
}

func renderUnauthorized(c *fiber.Ctx) error {
	return renderErrorPage(c, fiber.StatusUnauthorized, "Unauthorized")
}

func renderForbidden(c *fiber.Ctx) error {
	return renderErrorPage(c, fiber.StatusForbidden, "Forbidden")
}
