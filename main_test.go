package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFileData(t *testing.T) {
	filename := "temp.txt"
	data := "Test 123"

	file, err := os.Create(filename)
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		// Remove temporary file
		os.Remove(filename)
	}()

	_, err = file.Write([]byte(data))
	if err != nil {
		t.Error(err)
		return
	}

	file.Close()

	// Testing function
	content, err := getFileData(filename)
	if err != nil {
		t.Error(err)
		return
	}

	assert.Equal(t, content, data)
}

func TestMain(m *testing.M) {
	m.Run()
	if testing.CoverMode() != "" {
		fmt.Println("Coverage is ", testing.Coverage())
	}
}
