// Package webserver
// Created by RTT.
// Author: teocci@yandex.com on 2022-Apr-26
package webserver

import (
	"embed"
	"fmt"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/teocci/go-fiber-web/src/webserver/apictrlr"
	"log"
	"mime"
	"net"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"

	"github.com/teocci/go-fiber-web/src/config"
)

const (
	defaultProtocol        = "http"
	defaultPage            = "page.html"
	defaultFaviconRoute    = "/favicon.ico"
	defaultFaviconFilePath = "web/static/favicon.ico"

	formatAddress        = "%s:%d"
	formatURL            = "%s://%s/%s"
	formatRelativePath   = "/%s"
	formatStaticFilePath = "web/static/%s"
)

var (
	f       embed.FS
	address string
)

func Start() {
	address = fmt.Sprintf(formatAddress, "", config.Data.Web.Port)
	_ = mime.AddExtensionType(".js", "application/javascript")

	engine := html.New("./views", ".tmpl")
	router := fiber.New(fiber.Config{
		Views:   engine,
		Network: "tcp4",
	})

	indexRoute := fmt.Sprintf(formatRelativePath, defaultPage)
	indexFilePath := fmt.Sprintf(formatStaticFilePath, defaultPage)

	router.Static("/", "./web/static/")

	router.Static(indexRoute, indexFilePath)
	router.Static(defaultFaviconRoute, defaultFaviconFilePath)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowHeaders:     "Origin, X-Requested-With, Content-Type, Accept, Authorization, x-access-token",
		AllowMethods:     "GET, POST, PUT, DELETE",
		ExposeHeaders:    "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type",
	}))

	api := router.Group("/api/v1")
	api.Get("/seller/:id", apictrlr.HandleSeller)
	api.Get("/products/seller/:id", apictrlr.HandleProductList)
	api.Get("/identical/:id", apictrlr.HandleIdentical)

	fmt.Println("[url]", urlFormat(address))

	err := router.Listen(address)
	if err != nil {
		log.Fatalln("Start HTTP Server error", err)
	}
}

func GetLocalIp() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "localhost"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String()
}

func addressFormat(a string) string {
	s := strings.Split(a, ":")
	if s[0] == "" {
		s[0] = GetLocalIp()
	}
	return strings.Join(s[:], ":")
}

func urlFormat(a string) string {
	host := addressFormat(a)
	s := fmt.Sprintf(formatURL, defaultProtocol, host, defaultPage)

	return s
}
