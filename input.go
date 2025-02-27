//go:build !wasm

package main

import (
	"errors"
	"io"
	"log"
	"os"
)

func runFile() {
	args := os.Args
	if len(args) != 2 {
		log.Fatal("Filename not present")
	}

	source, err := getFileData(args[1])
	if err != nil {
		log.Fatal(err)
	}

	scanner := NewScanner(source)
	tokens := scanner.ScanTokens()

	for _, t := range tokens {
		DPrintf("%s %s\n", getTokenTypeString(t.TokenType), t.Lexeme)
	}

	parser := NewParser(tokens)
	parser.Parse()
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
