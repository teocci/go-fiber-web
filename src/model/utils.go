// Package model
// Created by RTT.
// Author: teocci@yandex.com on 2022-12월-16
package model

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

var (
	httpClient = &http.Client{Timeout: 10 * time.Second}
)

func parseWBToken(sellerId string, mode string) string {
	if mode == "" {
		mode = "STD"
	}

	return os.Getenv(fmt.Sprintf("WB_AUTH_TOKEN_%s_%s", mode, sellerId))
}
