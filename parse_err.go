package csvscan

import (
	"fmt"
)

type ParseErr struct {
	Line int
	Col int
	RawLine any
	Inner error
}

// Shorthand to create parse err from an error, also a shorthand that should inline from the compiler
func wrapParseErr(inner error, line, col int, rawLine []string) *ParseErr{
	return &ParseErr{Line: line, Col: col, RawLine: rawLine, Inner: inner}
}

// Shorthand to create parse err that the compiler will likely just inline
func newParseErr(line, col int, rawLine []string, errMsg string, fmtArgs ...any) *ParseErr{
	return &ParseErr{Line: line, Col: col, RawLine: rawLine, Inner: fmt.Errorf(errMsg, fmtArgs...)}
}

func (p *ParseErr) Error() string {
	return fmt.Sprintf("Line %v, Column %v: %v\nRaw line: %v", p.Line, p.Col, p.Inner, p.RawLine)
}

func (p *ParseErr) Unwrap() error {
	return p.Inner
}

func (p *ParseErr) Is(e error) bool {
	_, ok := e.(*ParseErr)
	return ok
}
