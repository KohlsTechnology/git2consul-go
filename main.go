package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/cleung2010/go-git2consul/config"
	"github.com/cleung2010/go-git2consul/runner"
)

const (
	ExitCodeError = 10 + iota
	ExitCodeFlagError
	ExitCodeConfigError

	ExitCodeOk int = 0
)

func main() {
	var filename string
	var version bool
	var debug bool
	var once bool

	flag.StringVar(&filename, "config", "", "path to config file")
	flag.BoolVar(&version, "version", false, "show version")
	flag.BoolVar(&debug, "debug", false, "enable debugging mode")
	flag.BoolVar(&once, "once", false, "run git2consul once and exit")
	flag.Parse()

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	if version {
		fmt.Println(Version)
		return
	}

	// TODO: Accept other logger inputs
	log.SetHandler(text.New(os.Stderr))

	log.Infof("Starting git2consul version: %s", Version)

	if len(filename) == 0 {
		log.Error("No configuration file provided")
		os.Exit(ExitCodeFlagError)
	}

	// Load configuration from file
	cfg, err := config.Load(filename)
	if err != nil {
		log.Errorf("(config): %s", err)
		os.Exit(ExitCodeConfigError)
	}

	runner, err := runner.NewRunner(cfg, once)
	if err != nil {
		log.Errorf("(runner): %s", err)
		os.Exit(ExitCodeConfigError)
	}
	go runner.Start()

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	for {
		select {
		case <-runner.DoneCh:
			log.Info("Finished execution of git2consul")
			os.Exit(ExitCodeOk)
		case err := <-runner.ErrCh:
			log.WithError(err).Error("Runner error")
			os.Exit(ExitCodeError)
		case <-signalCh:
			log.Info("Received interrupt. Terminating git2consul")
			os.Exit(ExitCodeOk)
		}
	}
}
