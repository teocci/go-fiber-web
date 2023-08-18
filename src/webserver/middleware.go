// Package webserver
// Created by RTT.
// Author: teocci@yandex.com on 2021-Nov-03
package webserver

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CORSMiddleware() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowHeaders:     "Origin, X-Requested-With, Content-Type, Accept, Authorization, x-access-token",
		AllowMethods:     "GET, POST, PUT, DELETE",
		ExposeHeaders:    "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type",
	})
}
