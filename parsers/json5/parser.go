package json5

import (
	"io/ioutil"

	"github.com/txgruppi/config/parsers"
	"github.com/yosuke-furukawa/json5/encoding/json5"
)

// NewParser returns a new JSON parser
func NewParser() parsers.Parser {
	return &parser{}
}

type parser struct {
}

func (t *parser) Exts() []string {
	return []string{".json5"}
}

func (t *parser) ParseFile(filepath string, v interface{}) error {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	if err := json5.Unmarshal(data, v); err != nil {
		return err
	}

	return nil
}
