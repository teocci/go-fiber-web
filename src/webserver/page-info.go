// Package webserver
// Created by RTT.
// Author: teocci@yandex.com on 2022-12월-24
package webserver

type PageInfo struct {
	Name       string `json:"name"`
	Action     string `json:"action,omitempty"`
	Tab        string `json:"tab,omitempty"`
	SupplierID string `json:"supplier_id,omitempty"`
	Limit      int    `json:"limit,omitempty"`
}
