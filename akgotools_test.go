package akgotools

import (
	"os"
	"path/filepath"
	"testing"
)

// Helper function to create test files
func createTestFiles(dir string, fileNames []string) error {
	for _, name := range fileNames {
		filePath := filepath.Join(dir, name)
		file, err := os.Create(filePath)
		if err != nil {
			return err
		}
		defer file.Close() // Close file after creation
	}
	return nil
}

// Helper function to get all file names in a directory
func getFileNames(dir string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var fileNames []string
	for _, entry := range entries {
		fileNames = append(fileNames, entry.Name())
	}
	return fileNames, nil
}

func TestRenameFiles(t *testing.T) {
	// Step 1: Create a temporary directory
	tempDir, err := os.MkdirTemp("", "testdir")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir) // Clean up

	// Step 2: Create test files
	testFileNames := []string{
		"BhavCopy_0_0_0_20241210_file1.csv",
		"BhavCopy_0_0_0_20241211_file2.csv",
		"BhavCopy_0_0_0_20241212_file3.txt", // Non-matching file
		"BhavCopy_0_0_0_20241212_file3.csv", // Non-matching file
		"BhavCopy_0_0_0_20241213_file4.csv",
		"BhavCopy_0_0_0_20241214_file4.csv",
	}
	err = createTestFiles(tempDir, testFileNames)
	if err != nil {
		t.Fatalf("Failed to create test files: %v", err)
	}

	// Log files before calling RenameFiles
	initialFiles, err := getFileNames(tempDir)
	if err != nil {
		t.Fatalf("Failed to list files before RenameFiles: %v", err)
	}
	t.Logf("Initial files: %v", initialFiles)

	// Step 3: Call RenameFiles
	pattern := "*.csv"
	err = RenameFiles(tempDir, pattern) // Replace with your actual function
	if err != nil {
		t.Fatalf("RenameFiles failed: %v", err)
	}

	// Log files after calling RenameFiles
	finalFiles, err := getFileNames(tempDir)
	if err != nil {
		t.Fatalf("Failed to list files after RenameFiles: %v", err)
	}
	t.Logf("Final files: %v", finalFiles)

	// Step 4: Verify renamed files
	expectedFileNames := []string{
		"20241210_bhav.csv",
		"20241211_bhav.csv",
		"20241212_bhav.csv",
		"20241213_bhav.csv",
		"20241214_bhav.csv",
		"BhavCopy_0_0_0_20241212_file3.txt", // Non-matching file should remain unchanged
	}

	expectedMap := make(map[string]struct{})
	for _, name := range expectedFileNames {
		expectedMap[name] = struct{}{}
	}

	for _, actualFileName := range finalFiles {
		if _, exists := expectedMap[actualFileName]; !exists {
			t.Errorf("Unexpected file found: %s", actualFileName)
		}
		delete(expectedMap, actualFileName)
	}

	for missingFile := range expectedMap {
		t.Errorf("Expected file not found: %s", missingFile)
	}
}

// Local Variables:
// jinx-local-words: "testdir"
// End:
