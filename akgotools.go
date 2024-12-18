// Package akgotools implements utility routines for manipulating NSE bhavcopies
// downloaded from NSE site.
// This package utilities are meant only for those files.
package akgotools

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
)

// RenameFiles rename the files
func RenameFiles(root, filepattern string) {

	// list all files to rename
	bhavcopies, err := filepath.Glob(root + filepattern)
	if err != nil {
		log.Fatalln("Bad Pattern,", err)
	}

	totalFiles := len(bhavcopies)
	n := totalFiles / 4

	// Check if filename contains "BhavCopy"
	for idx, bc := range bhavcopies {
		dir, file := filepath.Split(bc)
		matched, err := regexp.MatchString(`BhavCopy`, file)
		if err != nil {
			log.Fatalln("Incorrect regex;", err)
			return
		}
		// if it contains, then get the date part from file
		if matched {
			re, err := regexp.Compile(`\d{8}`)
			if err != nil {
				log.Fatalln("regex not correct")
			}
			dt := re.FindString(file)
			newfname := filepath.Join(dir, dt+"_bhav.csv")
			os.Rename(bc, newfname)
		}
		// print progress, at every 25% work done
		if (idx+1)%n == 0 {
			fmt.Printf("Renamed %d files...\n", idx+1)
		}
	}
}

// MergeFiles merges multiple csv files into one csv file
// It checks the root directory by given filepattern to list files to be merged
// and merges in newfn csv file
func MergeFiles(root, filepattern, newfn string) {
	// get all the relevant files
	bhavcopies, err := filepath.Glob(root + filepattern)
	if err != nil {
		log.Fatalln("Bad Pattern,", err)
	}

	totalFiles := len(bhavcopies)
	n := totalFiles / 4

	// merge files
	mergedfile, err := os.Create(newfn)
	if err != nil {
		log.Fatalln("Error creating file:", err)
	}
	defer mergedfile.Close()

	writer := csv.NewWriter(mergedfile)
	defer writer.Flush()

	for idx, bc := range bhavcopies {
		f, err := os.Open(bc)
		if err != nil {
			log.Fatalln("Error:", err)
			return
		}
		defer f.Close()
		reader := csv.NewReader(f)
		records, err := reader.ReadAll()
		if err != nil {
			log.Fatalln("Error reading file:", err)
			return
		}
		if idx > 0 {
			records = records[1:]
		}
		for _, record := range records {
			if err := writer.Write(record); err != nil {
				log.Fatalln("Error writing record to file:", err)
			}
		}
		if (idx+1)%n == 0 {
			fmt.Printf("Merged %d files...\n", idx+1)
		}
	}
}
