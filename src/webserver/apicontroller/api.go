// Package apicontroller
// Created by RTT.
// Author: teocci@yandex.com on 2023-2월-27
package apicontroller

type apiResponse struct {
	Data    interface{} `json:"data,omitempty"`
	Success bool        `json:"success,omitempty"`
	Err     string      `json:"error,omitempty"`
}
