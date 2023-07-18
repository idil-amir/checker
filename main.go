package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strings"
)

var (
	file1      = "files/input/file1.csv"
	file2      = "files/input/file2.csv"
	outputFile = "files/output/non_existing_rows.csv"
)

func main() {

	// Read the first CSV file
	rows1, err := readCSVFile(file1)
	if err != nil {
		fmt.Println("Error reading file1:", err)
		return
	}

	// Read the second CSV file
	rows2, err := readCSVFile(file2)
	if err != nil {
		fmt.Println("Error reading file2:", err)
		return
	}

	// Get non-existing rows in file2 compared to file1
	nonExistingRows := getNonExistingRows(rows1, rows2)

	if len(nonExistingRows) == 0 {
		fmt.Println("All rows in file1 exist in file2.")
	} else {
		fmt.Printf("The following rows in file1 do not exist in file2. Writing to %s\n", outputFile)
		err := writeNonExistingRowsToFile(nonExistingRows, outputFile)
		if err != nil {
			fmt.Println("Error writing non-existing rows to file:", err)
		}
	}
}

// ReadCSVFile reads a CSV file and returns its rows
func readCSVFile(filename string) ([][]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// GetNonExistingRows returns rows in rows1 that do not exist in rows2
func getNonExistingRows(rows1, rows2 [][]string) [][]string {
	existingRows := make(map[string]bool)

	// Add all rows in rows2 to the existingRows map
	for _, row := range rows2 {
		key := strings.Join(row, ",")
		existingRows[key] = true
	}

	nonExistingRows := [][]string{}

	// Check if each row in rows1 exists in the existingRows map
	for i, row := range rows1 {
		// header
		if i == 0 {
			nonExistingRows = append(nonExistingRows, row)
		}

		key := strings.Join(row, ",")
		if !existingRows[key] {
			nonExistingRows = append(nonExistingRows, row)
		}
	}

	return nonExistingRows
}

// WriteNonExistingRowsToFile writes the non-existing rows to a CSV file
func writeNonExistingRowsToFile(rows [][]string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, row := range rows {
		err := writer.Write(row)
		if err != nil {
			return err
		}
	}

	return nil
}
