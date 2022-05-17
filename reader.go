package csvscan

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"
)

// Reader will read from an input stream/file and parse out the CSV
// into the type you instantiate it with
type Reader[T any] struct {
	IgnoreHeader bool

	// Force a specific column length
	ForceColumnLength int
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
	if t.Kind() != reflect.Struct {
		return nil, fmt.Errorf("csvscan.Reader only is able to work with structs")
	}

	fieldMap, err := generateFieldMap(t)
	if err != nil {
		return nil, err
	}

	rows := []T{}
	i := 1
	for {
		rawRowAsStrings, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}

			return nil, err
		}

		var row T
		v := reflect.ValueOf(&row).Elem()

		n := len(rawRowAsStrings)
		if r.ForceColumnLength > 0 && n != r.ForceColumnLength {
			return nil, newParseErr(i, 0, rawRowAsStrings, "enforced length for columns is %v but got %v", r.ForceColumnLength, n)
		}

		for j := 0; j < len(rawRowAsStrings); j++ {
			fieldToAssign := fieldMap[j]
			if fieldToAssign == nil {
				continue
			}

			s := rawRowAsStrings[j]
			f := v.Field(*fieldToAssign)

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

// generateFieldMap examines the type's "csv" field tag which contains an
// integer. That integer tells the Reader which column in the CSV should be
// mapped to in the struct (zero indexed):
//
// 	 type Example struct {
// 	    Column1 int `csv:"1"`
//   }
//
// generateFieldMap will then generate the map
//
//   {1: 0}
//
// Which means when reading the CSV row "1,2,3,4" it will take 2 and put it
// in Column1
func generateFieldMap(t reflect.Type) (map[int]*int, error){
	n := t.NumField()
	fieldMap := make(map[int]*int, n)

	for i := 0; i < n; i++ {
		f := t.Field(i)
		csvTag := f.Tag.Get("csv")
		if csvTag == "" {
			continue
		}

		columnIndex, err := strconv.ParseInt(csvTag, 10, 32)
		if err != nil {
			return nil, fmt.Errorf("invalid CSV tag for field %v: %v", f.Name, csvTag)
		}

		if columnIndex < 0 {
			return nil, fmt.Errorf("invalid CSV tag for field %v: can't have a negative value for the index, got %v", f.Name, csvTag)
		}
		copyI := i
		fieldMap[int(columnIndex)] = &copyI
	}

	if len(fieldMap) == 0 {
		return nil, fmt.Errorf("no CSV tags specified, can't tell what struct fields to map which column to")
	}

	return fieldMap, nil
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
