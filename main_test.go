package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetFiles(t *testing.T) {
	cf, err := os.Getwd()
	if err != nil {
		t.Fatal()
	}
	s := getFiles(cf + string(filepath.Separator) + "testdata")
	if len(s) != 9 {
		t.Errorf("Got incorrect number of files. Expected 9, got %s.", s)
	}
}

func TestWriteContent(t *testing.T) {
	cf, err := os.Getwd()
	if err != nil {
		t.Fatal()
	}
	dir := cf + string(filepath.Separator) + "testdata" + string(filepath.Separator)
	s := getFiles(dir)
	writeContentToCsv(dir, "result", s)

	pathToFile := cf + string(filepath.Separator) + "result.csv"

	if _, err := os.Stat(pathToFile); os.IsNotExist(err) {
		t.Fatal("Resulting csv was not created.")
	}

	os.Remove(pathToFile)

}
