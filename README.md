# csvscan

Converts CSV values into Go structs using generics and also allows reading from a CSV
using those structs via field tags

**Go get**

``` sh
go get -u github.com/AnthonyHewins/csvscan 
```

**Library usage**

See the [examples directory](github.com/AnthonyHewins/csvscan/examples)

**CLI Usage**

```
usage: csvscan ARGUMENTS [OPTIONS]

Usage of ARGUMENTS:
  FILENAME              The filename to read from to parse
  help, -h, --help, h   Display this help text.

Usage of OPTIONS:
  -no-header
        Treat the first row as data (can't generate field names if this is the case)
  -package string
        Append a package name. Will omit if package is the empty string
```
