package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func main() {
	var (
		path        = flag.String("p", "./", "Directory path where files are located")
		csvFileName = flag.String("d", "result", "Resulting csv file name")
	)
	flag.Parse()

	directory := *path

	if !strings.HasSuffix(directory, string(filepath.Separator)) {
		directory += string(filepath.Separator)
	}

	txtFiles := getFiles(directory)

	if len(txtFiles) > 0 {
		writeContentToCsv(directory, *csvFileName, txtFiles)
	}
}

func getFiles(p string) []string {
	files, err := ioutil.ReadDir(p)
	if err != nil {
		panic(err)
	}
	return filterTxtFiles(files)
}

func filterTxtFiles(f []os.FileInfo) []string {
	var txtFiles []string
	for _, v := range f {
		// Ignore all directories, used to skip *.txt named directories
		if !v.IsDir() {
			// Append only files with .txt extension
			if filepath.Ext(v.Name()) == ".txt" {
				txtFiles = append(txtFiles, v.Name())
			}

		}
	}
	return txtFiles
}

func writeContentToCsv(path, resFileName string, files []string) {

	// Regex to escape all special characters expect dot
	reg, _ := regexp.Compile("[^a-zA-Z0-9.]+")

	// Create new result.csv file
	csvFile, err := os.Create(resFileName + ".csv")
	if err != nil {
		log.Fatalf("Cannot create file: %v", err)
	}
	defer csvFile.Close()

	csvWriter := csv.NewWriter(csvFile)
	defer csvWriter.Flush()

	for _, v := range files {
		file, err := os.Open(path + v)
		if err != nil {
			log.Fatalf("Cannot open file %s, due to error %v", v, err)
		}
		// Get fileName, minus the dot and extension
		fileName := []string{v[:len(v)-4]}

		scanner := bufio.NewScanner(file)
		// Skip first three title rows in every file
		for i := 0; i < 3; i++ {
			scanner.Scan()
		}
		for scanner.Scan() {
			// Regex used to replace all special characters. For some reason, trimming space didn't work.
			trimStr := reg.ReplaceAllString(scanner.Text(), "")
			// If line is not empty
			if len(trimStr) > 0 {
				var row []string
				// Split into two parts by dot between them
				content := strings.Split(trimStr, ".")
				// bottom part is used for handling special cases when there is no extension
				if len(content) == 2 {
					row = append(row, "."+content[1], content[0])
				} else {
					row = append(row, "", content[0])
				}
				data := append(fileName, row...)
				err := csvWriter.Write(data)
				if err != nil {
					log.Fatalf("Cannot write to csv file %v", err)
				}
			}

		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

	}
}
