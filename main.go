package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/go-sql-driver/mysql"

	"sqlcmder/internal/app"
	"sqlcmder/components"
	"sqlcmder/internal/helpers/logger"
)

var version = "dev"

func main() {
	defaultConfigPath, err := app.DefaultConfigFile()
	if err != nil {
		log.Fatalf("Error getting default config file: %v", err)
	}
	flag.Usage = func() {
		f := flag.CommandLine.Output()
		fmt.Fprintln(f, "sqlcmder")
		fmt.Fprintln(f, "")
		fmt.Fprintf(f, "Usage:  %s [options] [connection_url]\n\n", os.Args[0])
		fmt.Fprintln(f, "  connection_url")
		fmt.Fprintln(f, "        database URL to connect to. Omit to start in picker mode")
		fmt.Fprintln(f, "")
		fmt.Fprintln(f, "Options:")
		flag.PrintDefaults()
	}
	// Get executable directory for default log file
	exePath, _ := os.Executable()
	exeDir := filepath.Dir(exePath)
	defaultLogFile := filepath.Join(exeDir, "sqlcmder.log")

	configFile := flag.String("config", defaultConfigPath, "config file to use")
	printVersion := flag.Bool("version", false, "Show version")
	logLevel := flag.String("loglevel", "debug", "Log level")
	logFile := flag.String("logfile", defaultLogFile, "Log file")
	flag.Parse()

	if *printVersion {
		println("SQLCmder version: ", version)
		os.Exit(0)
	}

	slogLevel, err := logger.ParseLogLevel(*logLevel)
	if err != nil {
		log.Fatalf("Error parsing log level: %v", err)
	}
	logger.SetLevel(slogLevel)

	// Always enable logging to file
	if err := logger.SetFile(*logFile); err != nil {
		log.Fatalf("Error setting log file: %v", err)
	}

	logger.Info("Starting SQLCmder...", nil)

	if err := mysql.SetLogger(log.New(io.Discard, "", 0)); err != nil {
		log.Fatalf("Error setting MySQL logger: %v", err)
	}

	// First load the config.
	if err = app.LoadConfig(*configFile); err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	// Now we can initialize the main pages.
	mainPages := components.MainPages()

	// Parse the command line arguments.
	args := flag.Args()

	switch len(args) {
	case 0:
		// Launch into the connection picker.
	case 1:
		// Set a connection from the command line.
		err := components.InitFromArg(args[0])
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatal("Only a single connection is allowed")
	}

	if err = app.App.Run(mainPages, *configFile); err != nil {
		log.Fatalf("Error running app: %v", err)
	}
}
