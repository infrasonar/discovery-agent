package main

import (
	"fmt"
	"log"
	"os"

	"github.com/akamensky/argparse"
	"github.com/infrasonar/go-libagent"
)

func parseArgs() error {
	dname, err := os.MkdirTemp("", "sampledir")
	if err != nil {
		log.Fatal("Failed to create temp folder")
	}
	fmt.Println("Temp dir name:", dname)

	parser := argparse.NewParser("print", "Prints provided string to stdout")
	daemon := parser.Flag(
		"d",
		"daemon",
		&argparse.Options{
			Required: false,
			Help:     "Run the discovery agent as a daemon",
		},
	)
	configPath := parser.String(
		"c",
		"config-path",
		&argparse.Options{
			Required: false,
			Help:     "Path to store the asset Id (not required when an asset Id is provided)",
			Default:  "",
		},
	)
	token := parser.String(
		"t",
		"token",
		&argparse.Options{
			Required: false,
			Help:     "Agent token",
			Default:  "",
		},
	)
	assetName := parser.String(
		"n",
		"asset-name",
		&argparse.Options{
			Required: false,
			Help:     "Asset name (only required for initial run if no asset Id is given)",
			Default:  "",
		},
	)
	assetId := parser.Int(
		"i",
		"asset-id",
		&argparse.Options{
			Required: false,
			Help:     "Asset Id to use. When no asset Id is stored and no Id is provided, a new asset will be created",
			Default:  0,
		},
	)

	// Parse input
	err = parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		return err
	}
	fmt.Println(*daemon)
	fmt.Println(*configPath)
	fmt.Println(*token)
	fmt.Println(*assetName)
	fmt.Println(*assetId)
	return nil
}

func main() {
	// Start collector
	log.Printf("Starting InfraSonar Discovery Agent v%s\n", version)

	// Initialize random
	libagent.RandInit()

	// Read arguments; as this discovery might start only once, it differs from
	// other agents which are scheduled
	if parseArgs() != nil {
		return
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

	// Wait for quit
	<-quit
}
