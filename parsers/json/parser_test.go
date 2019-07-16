package json_test

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/txgruppi/config/parsers/json"
)

type server struct {
	Bind string
	Port int
}

type config struct {
	HTTPServer server `json:"server"`
}

func TestParser(t *testing.T) {
	t.Run("NewParser()", func(t *testing.T) {
		parser := json.NewParser()
		assert.NotNil(t, parser)
	})

	t.Run("Parser.Exts()", func(t *testing.T) {
		parser := json.NewParser()
		assert.Equal(t, []string{".json"}, parser.Exts())
	})

	t.Run("Parser.ParseFile()", func(t *testing.T) {
		cases := []struct {
			name       string
			data       []byte
			shouldFail bool
			expected   *config
			subject    *config
		}{
			{"no file", nil, true, nil, nil},
			{"empty file", []byte(""), true, nil, nil},
			{"empty json object", []byte("{}"), false, &config{}, nil},
			{
				"sample data",
				[]byte(`{"server":{"bind":"0.0.0.0","port":80}}`),
				false,
				&config{HTTPServer: server{"0.0.0.0", 80}},
				nil,
			},
			{
				"only replace fields present in file",
				[]byte(`{"server":{"bind":"127.0.0.1"}}`),
				false,
				&config{
					HTTPServer: server{
						Bind: "127.0.0.1",
						Port: 8080,
					},
				},
				&config{
					HTTPServer: server{
						Bind: "0.0.0.0",
						Port: 8080,
					},
				},
			},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				var err error
				var file *os.File
				file, err = ioutil.TempFile("", "config-test-data.*.json")
				assert.Nil(t, err)
				if c.data != nil {
					_, err = file.Write(c.data)
					assert.Nil(t, err)
					assert.Nil(t, file.Close())
				} else {
					assert.Nil(t, file.Close())
					assert.Nil(t, os.Remove(file.Name()))
				}

				parser := json.NewParser()
				err = parser.ParseFile(file.Name(), &c.subject)
				if c.shouldFail {
					assert.NotNil(t, err)
					return
				}
				assert.Nil(t, err)
				assert.Equal(t, c.expected, c.subject)
			})
		}
	})
}
