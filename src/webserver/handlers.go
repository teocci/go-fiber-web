// Package webserver
// Created by RTT.
// Author: teocci@yandex.com on 2022-Apr-26
package webserver

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"io"
	"net/http"
	"path/filepath"
)

const (
	imgURLFormat = "https://images.wbstatic.net/shops/%s_logo.jpg"
	imgDirPath   = "./web/static/img"
)

var (
	imgExtensions = map[string]string{
		"jpg":  "image/jpeg",
		"jpeg": "image/jpeg",
		"png":  "image/png",
	}
)

func HandleLogoImage(c *fiber.Ctx) error {
	profileID := c.Params("id")
	url := fmt.Sprintf(imgURLFormat, profileID)

	image, err := fetchLogoImage(profileID)
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
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
