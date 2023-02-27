// Package cmd
// Created by RTT.
// Author: teocci@yandex.com on 2021-Oct-18
package cmd

import (
	"fmt"
	"log"

	"github.com/spf13/cobra"

	"github.com/teocci/go-fiber-web/src/cmd/cmdapp"
	"github.com/teocci/go-fiber-web/src/config"
	"github.com/teocci/go-fiber-web/src/core"
	"github.com/teocci/go-fiber-web/src/logger"
)

var (
	app = &cobra.Command{
		Use:           cmdapp.Name,
		Short:         cmdapp.Short,
		Long:          cmdapp.Long,
		PreRunE:       validate,
		RunE:          runE,
		SilenceErrors: false,
		SilenceUsage:  false,
	}

	errs chan error
)

// Add supported cli commands/flags
func init() {
	cobra.OnInitialize(initConfig)

	config.AddFlags(app)
}

// Load config
func initConfig() {
	if err := config.LoadConfigFile(); err != nil {
		log.Fatal(err)
	}

	config.LoadLogConfig()
}

func validate(ccmd *cobra.Command, args []string) error {
	if config.Version {
		fmt.Printf(cmdapp.VersionTemplate, cmdapp.Name, cmdapp.Version, cmdapp.Commit)

		return nil
	}
	log.Printf("#%v", config.File)

	if !config.Verbose {
		ccmd.HelpFunc()(ccmd, args)

		return fmt.Errorf("")
	}

	return nil
}

func runE(ccmd *cobra.Command, args []string) error {
	var err error
	config.Log, err = logger.New(config.LogConfig)
	if err != nil {
		return ErrCanNotLoadLogger(err)
	}

	// make channel for errors
	errs = make(chan error)

	go runApp()

	// break if any of them return an error (blocks exit)
	if err = <-errs; err != nil {
		config.Log.Fatal(err)
	}

	return err
}

func runApp() {
	errs <- core.Start()
}

func Execute() {
	err := app.Execute()
	hasError(err)
}
