package main

import "fmt"

var (
	DEBUG = false
)

func DPrintf(pattern string, content ...any) {
	if DEBUG {
		fmt.Printf(pattern, content...)
	}
}
