// Package webserver
// Created by RTT.
// Author: teocci@yandex.com on 2022-Apr-27
package webserver

import (
	"net/http"
	"os"
	"path"
	"strings"
)

type ServeFileSystem interface {
	http.FileSystem
	Exists(prefix string, path string) bool
}

type LocalFileSystem struct {
	http.FileSystem
	root    string
	indexes bool
}

func (l *LocalFileSystem) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		name := path.Join(l.root, p)
		stats, err := os.Stat(name)
		if err != nil {
			return false
		}
		if !l.indexes && stats.IsDir() {
			return false
		}

		return true
	}

	return false
}
