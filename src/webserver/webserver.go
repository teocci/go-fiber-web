// Package webserver
// Created by RTT.
// Author: teocci@yandex.com on 2022-Apr-26
package webserver

import (
	"fmt"
	"log"
	"mime"
	"net"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/teocci/go-fiber-web/src/config"
	"github.com/teocci/go-fiber-web/src/webserver/apicontroller"
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
	address string
)

func Start() {
	address = fmt.Sprintf(formatAddress, "", config.Data.Web.Port)
	_ = mime.AddExtensionType(".js", "application/javascript")

	engine := html.New("./views", ".tmpl")
	router := fiber.New(fiber.Config{
		Views:       engine,
		ViewsLayout: "layouts/base",
		Network:     "tcp4",
	})

	indexRoute := fmt.Sprintf(formatRelativePath, defaultPage)
	indexFilePath := fmt.Sprintf(formatStaticFilePath, defaultPage)

	router.Static("/", "./web/static/")

	router.Static(indexRoute, indexFilePath)
	router.Static(defaultFaviconRoute, defaultFaviconFilePath)

	router.Use(CORSMiddleware())

	router.Get("/seller/logo/:id", handleLogoImage)
	router.Get("/seller/:id", handleSellerView)
	router.Get("/positions/:id", handlePositionsView)
	router.Get("/ads/:id", handleAdsView)

	api := router.Group("/api/v1")
	api.Get("/seller/:id", apicontroller.HandleSeller)
	api.Get("/marketing/:id", apicontroller.HandleMarketing)
	api.Get("/marketing/:id/:action/:campaign", apicontroller.HandleMarketing)
	api.Get("/filters/:action/:id", apicontroller.HandleFilter)
	api.Get("/filters/:action/:id/:key", apicontroller.HandleFilter)
	api.Get("/positions/:action/:id", apicontroller.HandlePositions)
	api.Get("/positions/:action/:id/:xsubject", apicontroller.HandlePositions)
	api.Get("/positions/:action/:id/cat/:cat", apicontroller.HandlePositions)
	api.Get("/positions/:action/:id/cat/:cat/:xsubject", apicontroller.HandlePositions)
	api.Get("/products/seller/:id", apicontroller.HandleProductList)
	api.Get("/identical/:id", apicontroller.HandleIdentical)

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
	defer func(conn net.Conn) {
		_ = conn.Close()
	}(conn)

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
