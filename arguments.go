package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"regexp"
	"strconv"

	"github.com/akamensky/argparse"
)

var reToken = regexp.MustCompile(`^[0-9a-f]{32}$`)
var tokenValidation = func(args []string) error {
	if !reToken.MatchString(args[0]) {
		return errors.New("invalid token")
	}
	return nil
}

func parseArgs() error {
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
			Default:  "/etc/infrasonar",
		},
	)
	token := parser.String(
		"t",
		"token",
		&argparse.Options{
			Required: false,
			Help:     "Agent token",
			Validate: tokenValidation,
		},
	)
	network := parser.String(
		"n",
		"network",
		&argparse.Options{
			Required: false,
			Help:     "Network to scan. For example: 192.168.1.0/24",
		},
	)
	assetName := parser.String(
		"",
		"asset-name",
		&argparse.Options{
			Required: false,
			Help:     "Asset name (only required for initial run if no asset Id is given)",
		},
	)
	assetId := parser.Int(
		"i",
		"asset-id",
		&argparse.Options{
			Required: false,
			Help:     "Asset Id to use. When no asset Id is stored and no Id is provided, a new asset will be created",
		},
	)
	apiUri := parser.String(
		"",
		"api-uri",
		&argparse.Options{
			Required: false,
			Help:     "InfraSonar API URI",
			Default:  "https://api.infrasonar.com",
		},
	)
	skipVerify := parser.Flag(
		"",
		"skip-verify",
		&argparse.Options{
			Required: false,
			Help:     "Skip the certificate verification test",
		},
	)
	checkNmapInterval := parser.Int(
		"",
		"check-nmap-interval",
		&argparse.Options{
			Required: false,
			Help:     "Interval for the NMAP check in seconds (recommended interval >= 60 as this is the minimum interval accepted by InfraSonar)",
			Default:  checkNmapDefaultInterval,
			Validate: func(args []string) error {
				if interval, err := strconv.Atoi(args[0]); err == nil {
					if interval < 0 || interval > 86400 {
						return errors.New("expecting an interval between 0 and 86400")
					}
				}
				return nil
			},
		},
	)
	parseFile := parser.String(
		"",
		"parse-file",
		&argparse.Options{
			Required: false,
			Help:     "Parse a scan XML file, print the parsed state output in JSON format, and quit",
		},
	)
	versionArg := parser.Flag(
		"v",
		"version",
		&argparse.Options{
			Required: false,
			Help:     "Print version and quit",
		},
	)

	// Parse input
	err := parser.Parse(os.Args)
	if err != nil {
		// In case of error print error and print usage
		// This can also be done by passing -h or --help flags
		fmt.Print(parser.Usage(err))
		return err
	}

	if *versionArg {
		fmt.Printf("Version: %s\n", version)
		os.Exit(0)
	}

	if *parseFile != "" {
		scan, err := parseXml(*parseFile)
		if err != nil {
			log.Fatal(err)
		}
		state := getState(scan)
		out, err := json.MarshalIndent(state, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s\n", out)
		os.Exit(0)
	}

	if os.Getenv("DAEMON") == "" {
		os.Setenv("DAEMON", fmt.Sprintf("%d", btoi(*daemon)))
	}
	if os.Getenv("CONFIG_PATH") == "" && *configPath != "" {
		os.Setenv("CONFIG_PATH", *configPath)
	}
	if os.Getenv("TOKEN") == "" && *token != "" {
		os.Setenv("TOKEN", *token)
	}
	if os.Getenv("NETWORK") == "" && *network != "" {
		os.Setenv("NETWORK", *network)
	}
	if os.Getenv("ASSET_NAME") == "" && *assetName != "" {
		os.Setenv("ASSET_NAME", *assetName)
	}
	if os.Getenv("ASSET_ID") == "" && *assetId > 0 {
		os.Setenv("ASSET_ID", fmt.Sprintf("%d", *assetId))
	}
	if os.Getenv("API_URI") == "" && *apiUri != "" {
		os.Setenv("API_URI", *apiUri)
	}
	if os.Getenv("SKIP_VERIFY") == "" {
		os.Setenv("SKIP_VERIFY", fmt.Sprintf("%d", btoi(*skipVerify)))
	}
	if os.Getenv("CHECK_NMAP_INTERVAL") == "" {
		os.Setenv("CHECK_NMAP_INTERVAL", fmt.Sprintf("%d", *checkNmapInterval))
	}
	return nil
}
