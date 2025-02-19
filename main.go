package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) == 2 {
		runFile(args[1])
	}
}

func runFile(filename string) {
	source, err := getFileData(filename)
	if err != nil {
		log.Fatal(err)
	}

	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()

	for _, token := range tokens {
		fmt.Printf("%+v\n", token)
	}
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
