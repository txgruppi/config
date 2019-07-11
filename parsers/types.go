package parsers

// Parser represents a data format data can be parsed and loaded into a type.
type Parser interface {
	// Exts returns a list of extensions for this data format.
	Exts() []string
	// ParseFile tries to parse and load the filepath into the given type.
	ParseFile(filepath string, v interface{}) error
}
