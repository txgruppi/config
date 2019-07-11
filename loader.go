package config

import (
	"os"
	"path"
	"sort"
	"strings"

	"github.com/txgruppi/config/parsers"
)

// Loader registers parsers and loads files into a given data type.
type Loader interface {
	// RegisterParser add a new parser to N extensions.
	// It will fail if: 1. the parser is `nil`; 2. any of its supported
	// extensions is already registered.
	RegisterParser(parsers.Parser) error
	// SupportedExtensions returns a list of registered extensions.
	SupportedExtensions() []string
	// Load find and loads files into the given data type.
	Load(v interface{}) (*Info, error)
}

// NewLoader returns a new Loader instance
func NewLoader() Loader {
	return &loader{parsers: map[string]parsers.Parser{}}
}

type loader struct {
	parsers map[string]parsers.Parser
}

func (t *loader) RegisterParser(parser parsers.Parser) error {
	if parser == nil {
		return &ErrNilParser{}
	}
	for _, ext := range parser.Exts() {
		if !strings.HasPrefix(ext, ".") {
			return &ErrMissingExtDot{Ext: ext}
		}
		if t.parsers[ext] != nil {
			return &ErrParserConflict{Ext: ext}
		}
		t.parsers[ext] = parser
	}
	return nil
}

func (t *loader) SupportedExtensions() []string {
	exts := make([]string, len(t.parsers))
	i := 0
	for k := range t.parsers {
		exts[i] = k
		i++
	}
	sort.Strings(exts)
	return exts
}

func (t *loader) Load(v interface{}) (*Info, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return nil, &ErrFailedToLoad{Reason: err}
	}
	hostname, err := os.Hostname()
	if err != nil {
		return nil, &ErrFailedToLoad{Reason: err}
	}
	deployment, err := getDeployment()
	if err != nil {
		return nil, &ErrFailedToLoad{Reason: err}
	}
	files := makePossibleFilepaths(configDir, hostname, deployment, t.SupportedExtensions())
	info := &Info{ConfigFolder: configDir, LoadedFiles: []string{}}
	for _, file := range files {
		_, err := os.Stat(file)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return info, err
		}
		ext := path.Ext(file)
		if err := t.parsers[ext].ParseFile(file, v); err != nil {
			return info, err
		}
		info.LoadedFiles = append(info.LoadedFiles, file)
	}

	return info, nil
}
