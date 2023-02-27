// Package cmd
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-27
package cmd

import (
	"errors"
	"fmt"
	"log"
)

const (
	errFileDoesNotExist = "%s file does not exist"
	errCanNotLoadLogger = "cannot load logger -> %s"
)

func hasError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func ErrFileDoesNotExist(f string) error {
	return errors.New(fmt.Sprintf(errFileDoesNotExist, f))
}

func ErrCanNotLoadLogger(e error) error {
	return errors.New(fmt.Sprintf(errCanNotLoadLogger, e))
}
