// Package apicontroller
// Created by Teocci.
// Author: teocci@yandex.com on 2023-Aug-13
package apicontroller

import "github.com/gofiber/fiber/v2"

func renderErrorAsJSON(c *fiber.Ctx, code int, msg string) error {
	return c.Status(code).JSON(fiber.Map{
		"error": msg,
	})
}

func renderInternalError(c *fiber.Ctx, msg string) error {
	if msg == "" {
		msg = "Internal Server Error"
	}
	return renderErrorAsJSON(c, fiber.StatusInternalServerError, msg)
}

func renderBadRequest(c *fiber.Ctx, msg string) error {
	if msg == "" {
		msg = "Bad Request"
	}

	return renderErrorAsJSON(c, fiber.StatusBadRequest, msg)
}

func renderNotFound(c *fiber.Ctx, msg string) error {
	if msg == "" {
		msg = "Not Found"
	}

	return renderErrorAsJSON(c, fiber.StatusNotFound, msg)
}

func renderUnauthorized(c *fiber.Ctx, msg string) error {
	if msg == "" {
		msg = "Unauthorized"
	}

	return renderErrorAsJSON(c, fiber.StatusUnauthorized, msg)
}

func renderForbidden(c *fiber.Ctx, msg string) error {
	if msg == "" {
		msg = "Forbidden"
	}

	return renderErrorAsJSON(c, fiber.StatusForbidden, msg)
}
