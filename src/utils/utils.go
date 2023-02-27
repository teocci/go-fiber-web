// Package utils
// Created by RTT.
// Author: teocci@yandex.com on 2022-12월-25
package utils

import (
	"math"
	"strconv"
	"strings"
	"time"
)

func FloatToUnix(t float64) int64 {
	sec, dec := math.Modf(t)
	return time.Unix(int64(sec), int64(dec*(1e9))).Unix()
}

func StringToUInt64(s string) (n uint64) {
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}

	return
}

func EmptyUserPass(username, password string) bool {
	return strings.Trim(username, " ") == "" || strings.Trim(password, " ") == ""
}

func StringToInt(s string) (n int) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return
}
