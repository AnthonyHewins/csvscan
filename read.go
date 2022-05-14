package main

import (
	"encoding/csv"
	"os"
)

func fetchRows(args *cliArgs) [][]string {
	f, err := os.Open(args.filename)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	numberOfRows := 2
	if args.noHeader {
		// no header -> we can't scan in names. Gotta go straight
		// to guessing data types
		numberOfRows--
	}

	rows := make([][]string, numberOfRows, numberOfRows)
	reader:= csv.NewReader(f)
	for i := 0; i < numberOfRows; i++ {
		row, err := reader.Read()
		if err != nil {
			panic(err)
		}

		rows[i] = row
	}

	return rows
}
