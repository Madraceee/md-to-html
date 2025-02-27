package main

import "fmt"

var (
	DEBUG = false
)

func DPrintf(pattern string, content ...interface{}) {
	if DEBUG {
		fmt.Printf(pattern, content...)
	}
}
