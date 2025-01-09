package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/infrasonar/go-libagent"
)

const checkNmapDefaultInterval = 14400

var errTokenMissing = libagent.ErrMedium("environment variable 'NETWORK' is missing")
var errAlreadyRunning = libagent.ErrLow("check is already running, please increase the nmap check interval")

var lock sync.Mutex
var scanProcess *exec.Cmd

func endProcess() {
	scanProcess = nil
}

func run(cmd *exec.Cmd) error {
	// Get a pipe to read from standard out
	r, _ := cmd.StdoutPipe()

	// Use the same pipe for standard error
	cmd.Stderr = cmd.Stdout

	// Make a new channel which will be used to ensure we get all output
	done := make(chan bool)

	// Create a scanner which scans r in a line-by-line fashion
	scanner := bufio.NewScanner(r)

	go func() {
		for scanner.Scan() {
			line := scanner.Text()
			fmt.Fprintln(os.Stderr, line)
		}
		done <- true
	}()

	// Start the command and check for errors
	err := cmd.Start()
	if err != nil {
		return err
	}

	// Wait for all output to be processed
	<-done

	// Wait for the command to finish
	err = cmd.Wait()

	return err
}

func CheckNmap(_ *libagent.Check) (map[string][]map[string]any, error) {
	if !lock.TryLock() {
		return nil, errAlreadyRunning
	}
	defer lock.Unlock()

	network := os.Getenv("NETWORK")
	if network == "" {
		return nil, errTokenMissing
	}

	if _, _, err := net.ParseCIDR(network); err != nil {
		return nil, err
	}

	workFile := os.Getenv("TMP_XML_FILE")

	scanProcess = exec.Command("nmap", "-sT", "-A", "-T4", "-oX", workFile, network)
	defer endProcess()

	log.Printf("Run: %s", strings.Join(scanProcess.Args, " "))
	if err := run(scanProcess); err != nil {
		return nil, err
	}

	scan, err := parseXml(workFile)
	if err != nil {
		return nil, err
	}

	state := getState(scan)
	return state, nil
}
