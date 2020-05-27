/*
Copyright 2019 Kohl's Department Stores, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/KohlsTechnology/git2consul-go/config"
	"github.com/KohlsTechnology/git2consul-go/pkg/version"
	"github.com/KohlsTechnology/git2consul-go/runner"
	"github.com/apex/log"
	"github.com/apex/log/handlers/json"
	"github.com/apex/log/handlers/text"
)

// Exit code represented as int values for particular errors.
const (
	ExitCodeError = 10 + iota
	ExitCodeFlagError
	ExitCodeConfigError

	ExitCodeOk int = 0
)

func main() {

	var (
		filename     string
		printVersion bool
		debug        bool
		once         bool
		logfmt       string
	)

	flag.StringVar(&filename, "config", "", "path to config file")
	flag.BoolVar(&printVersion, "version", false, "show version")
	flag.BoolVar(&debug, "debug", false, "enable debugging mode")
	flag.BoolVar(&once, "once", false, "run git2consul once and exit")
	// allow switching logformat. Structured output helps with parsers
	flag.StringVar(&logfmt, "logfmt", "text", "specify log format [text | json] ")
	flag.Parse()

	if printVersion {
		version.Print()
		return
	}

	// Init checks
	if len(filename) == 0 {
		log.Error("No configuration file provided")
		flag.Usage()
		os.Exit(ExitCodeFlagError)
	}

	if debug {
		log.SetLevel(log.DebugLevel)
	}

	// TODO: Accept other logger inputs
	switch logfmt {
	case "text":
		log.SetHandler(text.New(os.Stderr))
	case "json":
		log.SetHandler(json.New(os.Stderr))
	}

	log.WithField("caller", "main").Infof("Starting git2consul version: %s", version.Version)

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
		case err := <-runner.ErrCh:
			log.WithError(err).Error("Runner error")
			os.Exit(ExitCodeError)
		case <-runner.SndDoneCh: // Used for cases like -once, where program is not terminated by interrupt
			log.Info("Terminating git2consul")
			os.Exit(ExitCodeOk)
		case <-signalCh:
			log.Info("Received interrupt. Cleaning up...")
			runner.Stop()
		}
	}
}
