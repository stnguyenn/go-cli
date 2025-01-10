package csv

import (
	"encoding/csv"
	"fmt"
	"os"
)

func List(path string) []map[string]string {
	ml := []map[string]string{}

	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	//write(file)

	reader := csv.NewReader(file)
	data, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}

	for _, row := range data {
		for _, col := range row {
			fmt.Printf("%s,", col)
		}
		fmt.Println()
	}

	return ml
}

// func write(f *os.File) {
// 	writer := csv.NewWriter(f)

// 	headers := []string{"name", "age", "gender"}
// 	writer.Write(headers)

// 	data := [][]string{
// 		{"Alice", "25", "Female"},
// 		{"Bob", "30", "Male"},
// 		{"Charlie", "35", "Male"},
// 	}

// 	for _, row := range data {
// 		writer.Write(row)
// 	}

// 	defer writer.Flush()
// }
