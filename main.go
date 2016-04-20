package main

import (
	"flag"
	"fmt"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/cleung2010/go-git2consul/config"
	"github.com/cleung2010/go-git2consul/exit"
	"github.com/cleung2010/go-git2consul/git2consul"
)

var quit = make(chan bool)

func init() {
	exit.Listen(func(os.Signal) {
		log.Info("Caught interrupt signal, terminating git2consul")
		close(quit)
	})
}

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

	c := new(config.Config)
	if len(filename) > 0 {
		c = config.Load(filename)
		log.Debugf("Using configuration: %+v", c)
	} else {
		log.Info("No configuration provided")
	}

	git2consul.CloneRepos(c)

	//Wait for shutdown signal
	<-quit
}
