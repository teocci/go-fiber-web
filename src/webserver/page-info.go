// Package webserver
// Created by RTT.
// Author: teocci@yandex.com on 2022-12월-24
package webserver

type PageInfo struct {
	Name       string `json:"name"`
	Controller string `json:"controller"`
	Action     string `json:"action,omitempty"`
	Tab        string `json:"tab,omitempty"`
	SellerID   string `json:"seller_id,omitempty"`
	Limit      int    `json:"limit,omitempty"`
}
