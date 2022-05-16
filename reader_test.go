package csvscan

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockCSV struct {
	A bool
	B uint8
	C uint16
	D uint32
	E uint64
	F uint
	G int8
	H int16
	I int32
	J int64
	K int
	L float32
	M float64
	N string
}

func TestRead(t *testing.T) {
	type test struct {
		name        string
		r           Reader[mockCSV]
		arg         string
		expected    []mockCSV
		expectedErr error
	}

	tests := []test{
		{
			"base case",
			Reader[mockCSV]{},
			"",
			[]mockCSV{},
			nil,
		},
		{
			"throws error when column lengths don't match",
			Reader[mockCSV]{},
			"true,2",
			nil,
			newParseErr(1, 0, []string{"true", "2"}, "column number mismatch: expected 14 but got 2"),
		},
		{
			"throws errors on invalid types",
			Reader[mockCSV]{},
			"po,2,3,4,5,6,7,8,9,10,11,0.9,1.2,string",
			nil,
			wrapParseErr(
				&strconv.NumError{Func: "ParseBool", Num: "po", Err: strconv.ErrSyntax},
				1,
				0,
				strings.Split("po,2,3,4,5,6,7,8,9,10,11,0.9,1.2,string", ","),
			),
		},
		{
			"ignores header when told to",
			Reader[mockCSV]{IgnoreHeader: true},
			"1,2,3,4\n",
			[]mockCSV{},
			nil,
		},
		{
			"reads values correctly",
			Reader[mockCSV]{},
			"true,2,3,4,5,6,7,8,9,10,11,0.9,1.2,string",
			[]mockCSV{
				{true, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 0.9, 1.2, "string"},
			},
			nil,
		},
		{
			"ignores specified columns when passed",
			Reader[mockCSV]{IgnoreCols: []int{1,5}},
			"true,2,3,4,5,6,7,8,9,10,11,0.9,1.2,string",
			[]mockCSV{
				{true, 0, 3, 4, 5, 0, 7, 8, 9, 10, 11, 0.9, 1.2, "string"},
			},
			nil,
		},
	}

	for _, v := range []int{-1, 14, 15} {
	tests = append(tests,
		test{
			fmt.Sprintf("when an ignore column is outside the range (%v not in [0,13]) it fails", v),
			Reader[mockCSV]{IgnoreCols: []int{v}},
			"true,2,3,4,5,6,7,8,9,10,11,0.9,1.2,string",
			nil,
		fmt.Errorf("invalid column value to ignore: %v. Only %v fields are available to assign to", v, 14),
		})
}

	for _, tc := range tests {
		reader := strings.NewReader(tc.arg)
		actual, actualErr := tc.r.Read(reader)
		assert.Equal(t, tc.expected, actual, tc.name)
		if !assert.Equal(t, tc.expectedErr, actualErr, tc.name) {
			fmt.Println(actualErr)
		}
	}
}
