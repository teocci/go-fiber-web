// Package config
// Created by RTT.
// Author: teocci@yandex.com on 2021-Aug-23
package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/teocci/go-fiber-web/src/logger"
)

const (
	llName  = "log-level"
	llShort = "l"
	llDesc  = "Log level to output [fatal|error|info|debug|trace]"

	lfName  = "log-file"
	lfShort = "L"
	lfDesc  = "File that will receive the logger output"

	sysName  = "syslog"
	sysShort = "S"
	sysDesc  = "Logger will output into the syslog"

	vName  = "verbose"
	vShort = "v"
	vDesc  = "Run in verbose mode"

	cfName  = "config-file"
	cfShort = "c"
	cfDesc  = "Configuration file to load"

	wName  = "web-dir"
	wShort = "w"
	wDesc  = "Web directory where the app will work"

	sdName  = "static-dir"
	sdShort = "s"
	sdDesc  = "Static directory where the app will work"

	tsdName  = "views-dir"
	tsdShort = "T"
	tsdDesc  = "Templates directory where the app will work"

	tdName  = "tmp-dir"
	tdShort = "t"
	tdDesc  = "Temporal directory where the app will work"

	verName  = "version"
	verShort = "V"
	verDesc  = "Print version info and exit"
)

const (
	defaultConfigFile   = "config.json"
	defaultTMPDirPath   = "./tmp"
	defaultWebDirPath   = "./web"
	defaultStaticDir    = "static"
	defaultTemplatesDir = "views"

	defaultWebPort = 8080
)

type WebServer struct {
	Port int `json:"port"`
}

type ServerSetup struct {
	Web WebServer `json:"web"`
}

var (
	Syslog  = false // Run logger will output into the syslog
	Verbose = true  // Run in verbose mode
	Version = false // Print version info and exit

	File          = defaultConfigFile                           // Configuration file to load
	TMPPath       = defaultTMPDirPath                           // Temporal directory
	WebPath       = defaultWebDirPath                           // Web directory
	StaticPath    = filepath.Join(WebPath, defaultStaticDir)    // Static directory
	TemplatesPath = filepath.Join(WebPath, defaultTemplatesDir) // Templates directory

	Data = &ServerSetup{
		Web: WebServer{defaultWebPort},
	}

	LogLevel  = "info"         // Log level to output [fatal|error|info|debug|trace]
	LogFile   = ""             // File where the logger will output
	LogConfig logger.LogConfig // configuration for the logger
	Log       *logger.Logger   // Central logger for the app
)

// AddFlags adds the available cli flags
func AddFlags(cmd *cobra.Command) {
	// core
	cmd.Flags().StringVarP(&LogLevel, llName, llShort, LogLevel, llDesc)
	cmd.Flags().StringVarP(&LogFile, lfName, lfShort, LogFile, lfDesc)
	cmd.Flags().BoolVarP(&Verbose, vName, vShort, Verbose, vDesc)
	cmd.Flags().BoolVarP(&Verbose, sysName, sysShort, Verbose, sysDesc)

	cmd.PersistentFlags().StringVarP(&File, cfName, cfShort, File, cfDesc)
	cmd.PersistentFlags().StringVarP(&TMPPath, tdName, tdShort, TMPPath, tdDesc)
	cmd.PersistentFlags().StringVarP(&WebPath, wName, wShort, WebPath, wDesc)
	cmd.PersistentFlags().StringVarP(&StaticPath, sdName, sdShort, StaticPath, sdDesc)
	cmd.PersistentFlags().StringVarP(&TemplatesPath, tsdName, tsdShort, TemplatesPath, tsdDesc)

	cmd.Flags().BoolVarP(&Version, verName, verShort, Version, verDesc)
}

func fileExtension(f string) string {
	if pos := strings.LastIndexByte(f, '.'); pos != -1 {
		return f[pos+1:]
	}

	return f
}

// LoadConfigFile reads the specified config file
func LoadConfigFile() error {
	if File == "" {
		return nil
	}

	// Set defaults to whatever might be there already
	viper.SetDefault(llName, LogLevel)
	viper.SetDefault(vName, Verbose)
	viper.SetDefault(tdName, TMPPath)

	dirPath := filepath.Dir(File)
	baseName := filepath.Base(File)
	filename := strings.TrimSuffix(baseName, filepath.Ext(baseName))
	ext := fileExtension(baseName)

	viper.SetConfigName(filename)
	viper.SetConfigType(ext)
	viper.AddConfigPath(dirPath)

	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("failed to read config file - %v", err)
	}

	err = viper.Unmarshal(&Data)
	if err != nil {
		return fmt.Errorf("unable to decode into struct, %v", err)
	}

	// Set values. Config file will override commandline
	LogLevel = viper.GetString(llName)
	Verbose = viper.GetBool(vName)
	TMPPath = viper.GetString(tdName)

	return nil
}

func LoadLogConfig() {
	LogConfig = logger.LogConfig{
		Level:   LogLevel,
		Verbose: Verbose,
		Syslog:  Syslog,
		LogFile: LogFile,
	}
}
