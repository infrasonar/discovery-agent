package main

import (
	"log"
	"os"

	"github.com/infrasonar/go-libagent"
)

func main() {
	// Read arguments; as this discovery might start only once, it differs from
	// other agents which are scheduled; (sets environment variable on success)
	if err := parseArgs(); err != nil {
		os.Exit(1)
	}

	// Start collector
	log.Printf("Starting InfraSonar Discovery Agent v%s\n", version)

	// Initialize random
	libagent.RandInit()

	// Create work path (sets DISCOVERY_WORK_PATH)
	if err := createTmpXmlFile(); err != nil {
		log.Fatal(err)
	}

	// Initialize Helper (make sure to read arguments first)
	libagent.GetHelper()

	// Set-up signal handler
	quit := make(chan bool)
	go libagent.SigHandler(quit)

	// Create Collector
	collector := libagent.NewCollector("discovery", version)

	// Create Asset
	asset := libagent.NewAsset(collector)

	// asset.Kind = "Linux"
	asset.Announce()

	checkNmap := libagent.Check{
		Key:             "nmap",
		Collector:       collector,
		Asset:           asset,
		IntervalEnv:     "CHECK_NMAP_INTERVAL",
		DefaultInterval: checkNmapDefaultInterval,
		NoCount:         false,
		SetTimestamp:    false,
		Fn:              CheckNmap,
	}

	if isDeamon() {
		// Plan the check
		go checkNmap.Plan(quit)

		// Wait for quit
		<-quit
	} else {
		checkNmap.Run()
	}
}
