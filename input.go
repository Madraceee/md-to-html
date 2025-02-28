//go:build !wasm

package main

import (
	"errors"
	"io"
	"log"
	"os"
)

func getContent() string {
	args := os.Args
	if len(args) != 2 {
		log.Fatal("Filename not present")
	}

	source, err := getFileData(args[1])
	if err != nil {
		log.Fatal(err)
	}

	return source
}

func getFileData(filename string) (string, error) {
	file, err := os.Open(filename)
	defer file.Close()
	if err != nil {
		return "", errors.New("File err: Unable to open file - " + err.Error())
	}

	content, err := io.ReadAll(file)
	if err != nil {
		return "", errors.New("File err: Unable to read file - " + err.Error())
	}

	return string(content), nil
}
