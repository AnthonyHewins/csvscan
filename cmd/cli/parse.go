package main

import (
	"fmt"
	"strconv"
)

func getTypes(row []string) []string {
	typeList := make([]string, len(row))

	for i, v := range row {
		typeList[i] = discoverType(v)
	}

	return typeList
}

func discoverType(s string) string {
	if _, err := strconv.ParseBool(s); err == nil {
		return "bool"
	}

	if _, err := strconv.ParseInt(s, 10, 64); err == nil {
		return "int"
	}

	if _, err := strconv.ParseFloat(s, 64); err == nil {
		return "float64"
	}

	// return best guess
	return "string"
}

func genGenericFieldNames(count int) []string {
	f := make([]string, count, count)
	for i := 0; i < count; i++ {
		f[i] = fmt.Sprintf("Field%v", i)
	}
	return f
}
