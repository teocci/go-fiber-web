// Package logger
// Created by Teocci.
// Author: teocci@yandex.com on 2021-Aug-23
//go:build windows
// +build windows

package logger

import (
	"fmt"
	"io"
	"strings"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/svc/eventlog"
)

type writer struct {
	pri Level
	src string
	el  *eventlog.Log
}

func (w *writer) Write(b []byte) (n int, err error) {
	switch w.pri {
	case sInfo:
		return len(b), w.el.Info(1, string(b))
	case sWarn:
		return len(b), w.el.Warning(3, string(b))
	case sDebug:
		return len(b), w.el.Error(2, string(b))
	}
	return 0, fmt.Errorf("unrecognized severity: %v", w.pri)
}

func (w *writer) Close() error {
	return w.el.Close()
}

func NewSyslog(prefix string) (io.WriteCloser, error) {
	wl, err := newW(sDebug, prefix)
	if err != nil {
		return nil, err
	}

	return wl, nil
}

func newW(pri Level, src string) (io.WriteCloser, error) {
	// Continue if we receive "registry key already exists" or if we get
	// ERROR_ACCESS_DENIED so that we can log without administrative permissions
	// for pre-existing eventlog sources.
	if err := eventlog.InstallAsEventCreate(src, eventlog.Info|eventlog.Warning|eventlog.Error); err != nil {
		if !strings.Contains(err.Error(), "registry key already exists") &&
			err != windows.ERROR_ACCESS_DENIED {
			return nil, err
		}
	}
	el, err := eventlog.Open(src)
	if err != nil {
		return nil, err
	}
	return &writer{
		pri: pri,
		src: src,
		el:  el,
	}, nil
}
