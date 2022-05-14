# csvscan

Converts CSV values into Go structs

**Go get**

``` sh
go get -u github.com/AnthonyHewins/csvscan 
```

**Usage**


``` shell
usage: csvscan ARGUMENTS [OPTIONS]

Usage of ARGUMENTS:
  FILENAME              The filename to read from to parse. Supply multiple
                        to have all of them put in the same file
  help, -h, --help, h   Display this help text.

Usage of OPTIONS:
  -no-header
        Treat the first row as data (can't generate field names if this is the case)
  -package string
        Append a package name. Will omit if package is the empty string
```
