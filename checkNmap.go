package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/infrasonar/go-libagent"
)

const checkNmapDefaultInterval = 14400

var lock sync.Mutex

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
			fmt.Println(line)
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
		return nil, &libagent.CheckError{Sev: libagent.Low, Err: errors.New("check is already running, please increase the nmap check interval")}
	}
	defer lock.Unlock()

	network := os.Getenv("NETWORK")
	if network == "" {
		return nil, errors.New("environment variable 'NETWORK' is missing")
	}

	workFile := os.Getenv("TMP_XML_FILE")

	cmd := exec.Command("nmap", "-sT", "-A", "-T4", "-oX", workFile, network)
	if err := run(cmd); err != nil {
		return nil, err
	}

	scan, err := parseXml(workFile)
	if err != nil {
		return nil, err
	}

	state := getState(scan)

	return state, nil
}
