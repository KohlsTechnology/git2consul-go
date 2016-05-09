package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/cleung2010/go-git2consul/config"
	"github.com/cleung2010/go-git2consul/runner"
)

// var quit = make(chan bool)
//
// func init() {
// 	exit.Listen(func(os.Signal) {
// 		log.Info("Caught interrupt signal, terminating git2consul")
// 		close(quit)
// 	})
// }

func main() {
	var filename string
	var v bool
	var d bool

	flag.StringVar(&filename, "config", "", "path to config file")
	flag.BoolVar(&v, "v", false, "show version")
	flag.BoolVar(&d, "d", false, "enable debugging mode")
	flag.Parse()

	if d {
		log.SetLevel(log.DebugLevel)
	}

	if v {
		fmt.Println(Version)
		return
	}

	log.Infof("Starting git2consul version: %s", Version)

	if len(filename) == 0 {
		log.Fatal("No configuration file provided")
	}

	// Load configuration from file
	cfg, err := config.Load(filename)
	if err != nil {
		log.Fatal(err)
	}

	//////////// NOTE: This is new
	runner, err := runner.NewRunner(cfg)
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
			log.Fatal(err)
		case <-signalCh:
			log.Info("Received interrupt. Terminating git2consul")
			os.Exit(0)
		}
	}

	////////////

	// Create repos from configuration
	// repos, err := repository.LoadRepos(cfg)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// // Watch for local changes to push to KV
	// client, err := consul.NewClient(cfg)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// err = client.WatchChanges(repos)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// // Watch for remote changes to pull locally
	// err = repos.WatchRepos()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// //Wait for shutdown signal
	// <-quit
}
