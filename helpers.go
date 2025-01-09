package main

import "os"

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func isDaemon() bool {
	return !(os.Getenv("DAEMON") == "" || os.Getenv("DAEMON") == "0")
}
