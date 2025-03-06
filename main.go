package main

import (
	"fmt"
	"os"
)

func main() {
	content := getContent()

	// Scanner
	scanner := NewScanner(content)
	tokens := scanner.ScanTokens()

	// Print tokens identified by scanner
	for _, t := range tokens {
		DPrintf("%s %s\n", getTokenTypeString(t.TokenType), t.Lexeme)
	}

	// Parser
	parser := NewParser(tokens)
	asts := parser.Parse()

	// Print ast on debug mode
	if DEBUG {
		astPrinter := AstPrinter{}
		for _, c := range asts {
			s, _ := c.Visit(&astPrinter)
			fmt.Println(s)
		}
	}

	htmlGen := NewHTMLGenerator()
	filename := os.Args[1]
	htmlGen.GenerateHTML(filename, asts)
}
