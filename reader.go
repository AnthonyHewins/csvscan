package csvscan

import (
	"encoding/csv"
	"io"
	"os"
	"reflect"
	"strconv"
)

// Reader will read from an input stream/file and parse out the CSV
// into the type you instantiate it with
type Reader[T any] struct {
	IgnoreHeader bool
}

// ReadFile opens the specified file, reads & parses it in its entirety,
// then closes it
func (r Reader[T]) ReadFile(filename string) ([]T, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return r.Read(f)
}

// Read takes in an io.Reader and reads it until io.EOF, parsing
// each row into a struct
func (r Reader[T]) Read(rawReader io.Reader) ([]T, error) {
	reader := csv.NewReader(rawReader)

	if r.IgnoreHeader {
		if _, err := reader.Read(); err != nil {
			return nil, err
		}
	}

	var zeroValue T
	t := reflect.TypeOf(zeroValue)

	rows := []T{}
	i, n := 1, t.NumField()
	for {
		rawRowAsStrings, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, err
		} else if m := len(rawRowAsStrings); m != n {
			return nil, newParseErr(i, 0, rawRowAsStrings, "column number mismatch: expected %v but got %v", n, m)
		}

		var row T
		v := reflect.ValueOf(&row).Elem()
		for j := 0; j < n; j++ {
			f := v.Field(j)
			s := rawRowAsStrings[j]

			switch k := f.Kind(); k {
			case reflect.Bool:
				err = setBool(f, s)
			case reflect.Uint8:
				err = setUInt(f, s, 8)
			case reflect.Uint16:
				err = setUInt(f, s, 16)
			case reflect.Uint32:
				err = setUInt(f, s, 32)
			case reflect.Uint64:
				err = setUInt(f, s, 64)
			case reflect.Uint:
				err = setUInt(f, s, strconv.IntSize)
			case reflect.Int8:
				err = setInt(f, s, 8)
			case reflect.Int16:
				err = setInt(f, s, 16)
			case reflect.Int32:
				err = setInt(f, s, 32)
			case reflect.Int64:
				err = setInt(f, s, 64)
			case reflect.Int:
				err = setInt(f, s, strconv.IntSize)
			case reflect.Float32:
				err = setFloat(f, s, 32)
			case reflect.Float64:
				err = setFloat(f, s, 64)
			case reflect.String:
				f.SetString(s)
			default:
				return nil, newParseErr(i, j, rawRowAsStrings, "unsupported type in struct: %v", k)
			}

			if err != nil {
				return nil, wrapParseErr(err, i, j, rawRowAsStrings)
			}
		}

		rows = append(rows, row)
		i++
	}

	return rows, nil
}

func setBool(f reflect.Value, s string) error {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}
	f.SetBool(b)
	return nil
}

func setFloat(f reflect.Value, s string, bits int) error {
	b, err := strconv.ParseFloat(s, bits)
	if err != nil {
		return err
	}
	f.SetFloat(b)
	return nil
}

func setUInt(f reflect.Value, s string, bits int) error {
	b, err := strconv.ParseUint(s, 10, bits)
	if err != nil {
		return err
	}
	f.SetUint(b)
	return nil
}

func setInt(f reflect.Value, s string, bits int) error {
	b, err := strconv.ParseInt(s, 10, bits)
	if err != nil {
		return err
	}
	f.SetInt(b)
	return nil
}
