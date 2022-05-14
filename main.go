package main

import (
	"flag"
	"fmt"
	"os"
)

var helpText = `usage: csvscan ARGUMENTS [OPTIONS]

Arguments
  FILENAME					   The filename to read from to parse. Supply multiple
							   to have all of them put in the same file
  help, -h, --help, h          Display this help text.
`

type cliArgs struct {
	filename string
	noHeader bool
	packageName string
}

func main() {
	// replace global help message
	flag.Usage = usage

	if len(os.Args) <= 1 {
		help(1, "Err: not enough args, please supply a filename")
	}

	switch os.Args[1] {
	case "-h", "--help", "help", "h":
		help(0)
	default:
	}

	args := flags(os.Args)
	rows := fetchRows(args)

	var fieldNames, typesList []string
	if len(rows) >= 2 {
		fieldNames = rows[0]
		typesList = getTypes(rows[1])
	} else {
		// gotta use placeholders
		fieldNames = genGenericFieldNames(len(rows[0]))
		typesList = getTypes(rows[0])
	}

	structString := genStruct("Struct", fieldNames, typesList)

	if args.packageName != "" {
		fmt.Printf("package %v\n\n", args.packageName)
	}

	fmt.Println(structString)
}

func flags(flagArgs []string) *cliArgs {
	a := cliArgs{filename: os.Args[1]}
	fs := flag.NewFlagSet("global", flag.ExitOnError)

	fs.BoolVar(&a.noHeader, "no-header", false, "Treat the first row as data (can't generate field names if this is the case)")
	fs.StringVar(&a.packageName, "package", "", "Append a package name. Will omit if package is the empty string")
	fs.Parse(flagArgs[2:])

	return &a
}

func help(exitCode int, extraMessages ...interface{}) {
	for _, v := range extraMessages {
		fmt.Println(v)
	}

	usage()
	os.Exit(exitCode)
}

func usage() {
	fmt.Printf(helpText)
	fmt.Println("\nOptions")
	flags([]string{"--help"})
	flag.PrintDefaults()
}
