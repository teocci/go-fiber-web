// Package webserver
// Created by RTT.
// Author: teocci@yandex.com on 2022-Apr-26
package webserver

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
)

const (
	imgURLFormat = "https://images.wbstatic.net/shops/%s_logo.jpg"
	imgEmptyPath = "./web/static/img/seller-empty-logo.jpg"
)

var (
	imgExtensions = map[string]string{
		"jpg":  "image/jpeg",
		"jpeg": "image/jpeg",
		"png":  "image/png",
	}
)

func handlePositionView(c *fiber.Ctx) error {
	sellerID := c.Params("id")
	page := PageInfo{
		Name:       "position",
		Controller: "position",
		SupplierID: sellerID,
	}

	return renderPage(c, page)
}

func handleSellerView(c *fiber.Ctx) error {
	sellerID := c.Params("id")
	limit := c.QueryInt("limit")
	page := PageInfo{
		Name:       "seller",
		Controller: "seller",
		SupplierID: sellerID,
		Limit:      limit,
	}

	return renderPage(c, page)
}

func handleLogoImage(c *fiber.Ctx) error {
	sellerID := c.Params("id")
	url := fmt.Sprintf(imgURLFormat, sellerID)

	image, err := fetchLogoImage(sellerID)
	if err != nil {
		url = imgEmptyPath
		image, err = os.ReadFile(url)
		if err != nil {
			return c.SendStatus(fiber.StatusInternalServerError)
		}
	}

	extension := filepath.Ext(url)
	contentType := imgExtensions[extension]

	c.Type(contentType)
	return c.SendStream(bytes.NewReader(image))
}

func fetchLogoImage(profileID string) (bytes []byte, err error) {
	url := fmt.Sprintf(imgURLFormat, profileID)
	var resp *http.Response
	resp, err = http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("not found")
	}

	bytes, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return
}
