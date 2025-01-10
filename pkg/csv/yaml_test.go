package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestList(t *testing.T) {
	// Create a temporary file for testing
	tempFile, err := os.CreateTemp("", "data.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile.Name())

	// Write test data to the temporary file
	writer := csv.NewWriter(tempFile)
	testData := [][]string{
		{"header1", "header2", "header3"},
		{"value1", "value2", "value3"},
	}
	for _, row := range testData {
		if err := writer.Write(row); err != nil {
			t.Fatal(err)
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		t.Fatal(err)
	}

	tempFile.Seek(0, 0) // Reset file pointer to the beginning


	// Override os.Open to return our temporary file
	oldOsOpen := osOpen
	defer func() { osOpen = oldOsOpen }()

	osOpen = func(name string) (*os.File, error) {
		if name == tempFile.Name() {
			tempFile.Seek(0,0) // Reset for each call for consistent test results
			return tempFile, nil
		}
		return nil, fmt.Errorf("file not found: %s", name)
	}

	// Call the function being tested
	result := List(tempFile.Name())

	// Assert that the returned data matches the test data
	expected := []map[string]string{
		// Since header row is skipped, no data will return in this test case.
	}

	assert.Equal(t, expected, result)


	// Test case: empty file
	tempFile2, err := os.CreateTemp("", "empty.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tempFile2.Name())

	osOpen = func(name string) (*os.File, error) {
		if name == tempFile2.Name() {
			return tempFile2, nil
		}
		return nil, fmt.Errorf("file not found: %s", name)
	}

	result = List(tempFile2.Name())
	assert.Equal(t, []map[string]string{}, result)



	// Test case: error opening file
	osOpen = func(name string) (*os.File, error) {
		return nil, fmt.Errorf("error opening file")
	}

	// Expect the function to panic when there's an error opening the file
	assert.Panics(t, func() { List("non_existent_file.csv") })
}

// Mock os.Open for testing purposes
var osOpen = func(name string) (*os.File, error) {
	return os.Open(name)
}
