//go:build wasm

package main

import (
	"fmt"

	"honnef.co/go/js/dom/v2"
)

func getContent() string {
	// Get input from text field
	el := dom.GetWindow().Document().QuerySelector("#input-field")
	inputEle := el.(*dom.HTMLTextAreaElement)

	input := inputEle.Value()
	input = input + string('\n')

	return input
}

// Update after code gen
func printContent(tokens []Token) {
	outputString := ""
	for _, token := range tokens {
		outputString += fmt.Sprintf("%v\n", GetTokenString(&token))
	}

	el := dom.GetWindow().Document().GetElementByID("output-field")
	outputEle := el.(*dom.HTMLTextAreaElement)

	outputEle.SetTextContent(outputString)
}
