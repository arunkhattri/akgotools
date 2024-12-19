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
func RenameFiles(root, filepattern string) error {
	// ensure root ends with a separator
	root = filepath.Clean(root) + string(os.PathSeparator)

	// list all files to rename
	bhavcopies, err := filepath.Glob(root + filepattern)
	if err != nil {
		log.Fatalln("Bad Pattern,", err)
		return err
	}

	totalFiles := len(bhavcopies)
	if totalFiles == 0 {
		log.Println("No files matched the pattern")
		return nil
	}
	n := totalFiles / 4

	// Regex compilation
	var reDt, reBC *regexp.Regexp
	reDt = regexp.MustCompile(`\d{8}`)
	reBC = regexp.MustCompile(`BhavCopy`)

	// Check if filename contains "BhavCopy"
	for idx, bc := range bhavcopies {
		dir, file := filepath.Split(bc)
		// if file name contains 'BhavCopy'
		if reBC.MatchString(file) {
			dt := reDt.FindString(file)
			if dt != "" {
				newfname := filepath.Join(dir, dt+"_bhav.csv")
				if err := os.Rename(bc, newfname); err != nil {
					log.Fatalln("Failed to rename file:", err)
					return err
				}
			}
		}
		// print progress, at every 25% work done
		if (idx+1)%n == 0 {
			fmt.Printf("Renamed %d files...\n", idx+1)
		}
	}
	return nil
}

// MergeFiles merges multiple csv files into one csv file
// It checks the root directory by given file pattern to list files to be merged
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
	mergedfile, err := os.Create(root + newfn)
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

// Local Variables:
// jinx-local-words: "BhavCopy filepath fmt newfn os"
// End:
