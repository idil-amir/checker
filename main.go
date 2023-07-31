package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var (
	sourceFile    = "files/input/source.csv"
	targetFile    = "files/input/target.csv"
	failedFile    = "files/output/failed.csv"
	unclaimedFile = "files/output/unclaimed.csv"
)

func main() {
	// Read the source CSV file
	source, err := readCSVFile(sourceFile)
	if err != nil {
		fmt.Println("Error reading source:", err)
		return
	}

	// Read the target CSV file
	target, err := readCSVFile(targetFile)
	if err != nil {
		fmt.Println("Error reading target:", err)
		return
	}

	// Get non-existing rows in target compared to source
	failedData, unclaimedData := getUnmatchedEntries(source, target)

	// Print to the output files if there any unmatched entries
	fmt.Println()
	fmt.Println("================Logs================")
	if len(failedData) > 1 {
		err = writeDataToFile(failedData, failedFile)
		if err != nil {
			fmt.Println("Error writing non-existing rows to file:", err)
		}
	} else {
		fmt.Println("No Failed Data")
	}

	if len(unclaimedData) > 1 {
		err = writeDataToFile(unclaimedData, unclaimedFile)
		if err != nil {
			fmt.Println("Error writing non-existing rows to file:", err)
		}
	} else {
		fmt.Println("No Unclaimed Data")
	}
	fmt.Println("================End of Logs================")
	fmt.Println()
	fmt.Println("==================Summary==================")
	fmt.Printf("# of Failed Injection : %d", len(failedData)-1)
	fmt.Println()
	fmt.Printf("# of Unclaimed Medali : %d", len(unclaimedData)-1)
	fmt.Println()
	fmt.Println("===============End of Summary==============")
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

// Getfailed returns rows in source that do not exist in target
func getUnmatchedEntries(source, target [][]string) (failed, unclaimed [][]string) {
	existingRows := make(map[string]string)

	// Add all rows in target to the existingRows map
	for _, row := range target {
		newRow := []string{row[0], row[1]}
		key := strings.Join(newRow, ",")
		existingRows[key] = row[2]
	}

	// Check if each row in source exists in the existingRows map
	for i, row := range source {
		// header
		if i == 0 {
			failed = append(failed, row)
			unclaimed = append(unclaimed, row)
			continue
		}

		newRow := []string{row[0], row[1]}
		key := strings.Join(newRow, ",")
		if val, ok := existingRows[key]; ok {
			status, _ := strconv.Atoi(val)
			if status < 4 {
				newRow = append(newRow, val)
				unclaimed = append(unclaimed, newRow)
			}
		} else {
			failed = append(failed, newRow)
		}
	}

	return failed, unclaimed
}

// writeDataToFile writes the non-existing rows to a CSV file
func writeDataToFile(rows [][]string, filename string) error {
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
