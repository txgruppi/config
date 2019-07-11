package config

import "fmt"

// ErrNilParser is returned when a register is called with parsers.Parser(nil).
type ErrNilParser struct {
}

func (*ErrNilParser) Error() string {
	return "trying to register a nil parser"
}

// ErrParserConflict is returned when a parser is trying to register an
// extension that is already registered.
type ErrParserConflict struct {
	Ext string
}

func (t *ErrParserConflict) Error() string {
	return fmt.Sprintf("parser conflict for %q", t.Ext)
}

// ErrMissingExtDot is returned when the extension is missing the dot prefix.
type ErrMissingExtDot struct {
	Ext string
}

func (t *ErrMissingExtDot) Error() string {
	return fmt.Sprintf("extension %q must start with a dot", t.Ext)
}

// ErrFailedToLoad is returned at any point that load process fails.
type ErrFailedToLoad struct {
	Reason error
}

func (t *ErrFailedToLoad) Error() string {
	return fmt.Sprintf("failed to load with reason: %s", t.Reason.Error())
}
