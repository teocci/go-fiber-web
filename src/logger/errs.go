// Package logger
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-27
package logger

import (
	"errors"
	"fmt"
)

const (
	errCanNotOpenLogFile = "%s log file cannot be opened -> %s"
	errCanNotInitSyslog  = "cannot initiate syslog -> %s"
)

func ErrCanNotOpenLogFile(p string, e error) error {
	return errors.New(fmt.Sprintf(errCanNotOpenLogFile, p, e))
}

func ErrCanNotInitSyslog(e error) error {
	return errors.New(fmt.Sprintf(errCanNotInitSyslog, e))
}
