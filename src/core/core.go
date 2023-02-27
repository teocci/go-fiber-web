// Package core
// Created by RTT.
// Author: teocci@yandex.com on 2022-Apr-26
package core

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/teocci/go-fiber-web/src/webserver"
)

func Start() error {
	pid := os.Getpid()
	fmt.Println("PID:", pid)

	go webserver.Start()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		sig := <-sigs
		log.Println(sig)
		done <- true
	}()
	log.Println("Server start awaiting signal")
	<-done
	log.Println("Server stop working by signal")

	return nil
}
