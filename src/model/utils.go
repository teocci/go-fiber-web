// Package model
// Created by RTT.
// Author: teocci@yandex.com on 2022-12월-16
package model

import (
	"net/http"
	"time"
)

var (
	httpClient = &http.Client{Timeout: 10 * time.Second}
)
