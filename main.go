package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	file1                   = "files/input/file1.csv"
	file2                   = "files/input/file2.csv"
	outputFile              = "files/output/non_existing_rows.csv"
	nonClaimedQuestUserFile = "files/output/non_claimed_quest_user.csv"
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
	nonExistingRows, nonClaimedQuestUser := getNonExistingRows(rows1, rows2)

	if len(nonExistingRows) == 0 {
		fmt.Println("All rows in file1 exist in file2.")
	} else {
		fmt.Println("================Summary================")
		fmt.Printf("# of Failed Injection : %d", len(nonExistingRows)-1)
		fmt.Println()
		fmt.Printf("# of Not Claimed Medali : %d", len(nonClaimedQuestUser)-1)
		fmt.Println()
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
func getNonExistingRows(rows1, rows2 [][]string) (nonExistingRows, nonClaimedQuestUserData [][]string) {
	existingRows := make(map[string]bool)

	// Add all rows in rows2 to the existingRows map
	for i, row := range rows2 {
		newRow := []string{row[0], row[1]}
		key := strings.Join(newRow, ",")
		existingRows[key] = true
		if i == 0 {
			nonClaimedQuestUserData = append(nonClaimedQuestUserData, row)
		}
		status, _ := strconv.Atoi(row[2])
		if status == 2 {
			nonClaimedQuestUserData = append(nonClaimedQuestUserData, row)
		}
	}

	err := writeNonExistingRowsToFile(nonClaimedQuestUserData, nonClaimedQuestUserFile)
	if err != nil {
		fmt.Println("Error writing non-existing rows to file:", err)
	}

	// Check if each row in rows1 exists in the existingRows map
	for i, row := range rows1 {
		// header
		if i == 0 {
			nonExistingRows = append(nonExistingRows, row)
		}

		newRow := []string{row[0], row[1]}
		key := strings.Join(newRow, ",")
		if !existingRows[key] {
			nonExistingRows = append(nonExistingRows, row)
		}
	}

	return nonExistingRows, nonClaimedQuestUserData
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
