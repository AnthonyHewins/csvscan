package main

import (
	"fmt"
	"strings"

	"github.com/AnthonyHewins/csvscan"
)

type Example struct {
	X            int    `csv:"0"`
	Y            int    `csv:"1"`
	IgnoredField string // has no csv tag
}

const csvWithNoHeader = "1,2,3\n4,5,6"

var csvWithHeader = fmt.Sprintf("header1,header2,header3\n%v", csvWithNoHeader)

func main() {
	csvWithOnlyValues := strings.NewReader(csvWithNoHeader)
	reader := csvscan.Reader[Example]{}
	sliceOfExample, err := reader.Read(csvWithOnlyValues)
	if err != nil {
		panic(err)
	}

	fmt.Println(sliceOfExample) // [{1 2 } {4 5 }]

	csvWithValuesAndHeader := strings.NewReader(csvWithHeader)
	reader.IgnoreHeader = true
	sliceOfExample, err = reader.Read(csvWithValuesAndHeader)
	if err != nil {
		panic(err)
	}

	fmt.Println(sliceOfExample) // Same output: [{1 2 } {4 5 }]
}
