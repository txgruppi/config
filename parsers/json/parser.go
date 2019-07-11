package json

import (
	"encoding/json"
	"io/ioutil"

	"github.com/txgruppi/config/parsers"
)

// NewParser returns a new JSON parser
func NewParser() parsers.Parser {
	return &parser{}
}

type parser struct {
}

func (t *parser) Exts() []string {
	return []string{".json"}
}

func (t *parser) ParseFile(filepath string, v interface{}) error {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(data, v); err != nil {
		return err
	}

	return nil
}
