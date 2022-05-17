package csvscan

import (
	"io"
	"strings"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
)

type mockCSV struct {
	Skip bool
	A bool `csv:"0"`
	B uint8 `csv:"1"`
	C uint16 `csv:"2"`
	D uint32 `csv:"3"`
	E uint64 `csv:"4"`
	F uint `csv:"5"`
	G int8 `csv:"6"`
	H int16 `csv:"7"`
	I int32 `csv:"8"`
	J int64 `csv:"9"`
	K int `csv:"10"`
	L float32 `csv:"11"`
	M float64 `csv:"12"`
	N string `csv:"13"`
	NeverTouched int `csv:"1000"`
}

func TestRead(mainTest *testing.T) {
	mainTest.Run("fails when the instantiated type is invalid", func(t *testing.T) {
		actual, err := Reader[int]{}.Read(io.Reader(nil))
		cupaloy.SnapshotT(t,actual, err)
	})

	mainTest.Run("fails on invalid CSV tags that arent ints", func(t *testing.T) {
		type X struct { Y int `csv:"aosdpad"` }
		actual, err := Reader[X]{}.Read(io.Reader(nil))
		cupaloy.SnapshotT(t,actual, err)
	})

	mainTest.Run("fails on invalid CSV tags that are negative indices", func(t *testing.T) {
		type X struct { Y int `csv:"-1"` }
		actual, err := Reader[X]{}.Read(io.Reader(nil))
		cupaloy.SnapshotT(t,actual, err)
	})

	type test struct {
		name        string
		r           Reader[mockCSV]
		arg         string
	}

	tests := []test{
		{
			"base case",
			Reader[mockCSV]{},
			"",
		},
		{
			"throws errors on invalid types",
			Reader[mockCSV]{},
			"po,2,3,4,5,6,7,8,9,10,11,0.9,1.2,string",
		},
		{
			"ignores header when told to",
			Reader[mockCSV]{IgnoreHeader: true},
			"1,2,3,4\n1,2,3,4\n",
		},
		{
			"enforces column length when specified, erroring if it doesn't satisfy",
			Reader[mockCSV]{ForceColumnLength: 5},
			"1,2,3,4\n",
		},
		{
			"reads values correctly",
			Reader[mockCSV]{ForceColumnLength: 14},
			"true,2,3,4,5,6,7,8,9,10,11,0.9,1.2,string",
		},
	}

	for _, tc := range tests {
		reader := strings.NewReader(tc.arg)

		mainTest.Run(tc.name, func(t *testing.T) {
			actual, actualErr := tc.r.Read(reader)
			cupaloy.SnapshotT(t, tc.name, actual, actualErr)
		})
	}
}
