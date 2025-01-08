package main

import "os"

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}

func isDeamon() bool {
	return !(os.Getenv("DEAMON") == "" || os.Getenv("DEAMON") == "0")
}
