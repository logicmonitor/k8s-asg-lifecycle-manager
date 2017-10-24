package main

import (
	"fmt"
	"os"

	"github.com/logicmonitor/k8s-asg-lifecycle-manager/pkg"
	"github.com/logicmonitor/k8s-asg-lifecycle-manager/pkg/config"
	"github.com/logicmonitor/k8s-asg-lifecycle-manager/pkg/stats"
	log "github.com/sirupsen/logrus"
)

func main() {
	// Application configuration
	config, err := config.GetConfig()
	if err != nil {
		fmt.Printf("Failed to retrieve configuration: %s", err)
		os.Exit(1)
	}

	// Set the logging level.
	if config.Debug {
		log.SetLevel(log.DebugLevel)
	}

	// Instantiate the base struct.
	base, err := nodeman.NewBase(config)
	if err != nil {
		log.Fatal(err.Error())
	}

	// Instantiate the application
	nodeman, err := nodeman.NewNodeMan(base)
	if err != nil {
		log.Fatal(err.Error())
	}

	go func() {
		// start the stats and health check server
		s := stats.Server{}
		s.Start()
	}()

	// Master of the Nodeman
	nodeman.Watch()
}
