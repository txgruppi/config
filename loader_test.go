package config_test

import (
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/txgruppi/config"
)

type p struct {
	mock.Mock
}

func (t *p) Exts() []string {
	args := t.Called()
	if args.Get(0) == nil {
		return nil
	}
	return args.Get(0).([]string)
}

func (t *p) ParseFile(f string, v interface{}) error {
	return t.Called(f, v).Error(0)
}

func createFile(filepath string) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	return file.Close()
}

func TestLoader(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "config-test-")
	assert.Nil(t, err)

	defaultJson := path.Join(tmpDir, "default.json")
	developmentYml := path.Join(tmpDir, "development.yml")
	localDevelopmentYaml := path.Join(tmpDir, "local-development.yaml")
	assert.Nil(t, createFile(path.Join(tmpDir, "default.json5")))
	assert.Nil(t, createFile(defaultJson))
	assert.Nil(t, createFile(developmentYml))
	assert.Nil(t, createFile(localDevelopmentYaml))

	os.Setenv("CUSTOM_CONFIG_DIR", tmpDir)
	os.Setenv("ENV", "development")

	loader := config.NewLoader()
	loader.SetEnvironmentVariableName("CUSTOM_CONFIG_DIR")
	assert.NotNil(t, loader)

	p0 := &p{}
	p1 := &p{}
	p2 := &p{}
	p3 := &p{}

	p0.
		On("Exts").Return([]string{".json"}).
		On("ParseFile", defaultJson, mock.Anything).Return(nil).Once()
	p1.
		On("Exts").Return([]string{".yaml", ".yml"}).
		On("ParseFile", developmentYml, mock.Anything).Return(nil).Once().
		On("ParseFile", localDevelopmentYaml, mock.Anything).Return(nil).Once()
	p2.
		On("Exts").Return([]string{".yml"})
	p3.
		On("Exts").Return([]string{"xml"})

	assert.NotNil(t, loader.SupportedExtensions())
	assert.Len(t, loader.SupportedExtensions(), 0)
	assert.Nil(t, loader.RegisterParser(p0))
	assert.Nil(t, loader.RegisterParser(p1))
	assert.NotNil(t, loader.SupportedExtensions())
	assert.Len(t, loader.SupportedExtensions(), 3)
	assert.Equal(t, []string{".json", ".yaml", ".yml"}, loader.SupportedExtensions())

	assert.EqualError(t, loader.RegisterParser(nil), "trying to register a nil parser")
	assert.Len(t, loader.SupportedExtensions(), 3)
	assert.Equal(t, []string{".json", ".yaml", ".yml"}, loader.SupportedExtensions())

	assert.EqualError(t, loader.RegisterParser(p2), `parser conflict for ".yml"`)
	assert.Len(t, loader.SupportedExtensions(), 3)
	assert.Equal(t, []string{".json", ".yaml", ".yml"}, loader.SupportedExtensions())

	assert.EqualError(t, loader.RegisterParser(p3), `extension "xml" must start with a dot`)
	assert.Len(t, loader.SupportedExtensions(), 3)
	assert.Equal(t, []string{".json", ".yaml", ".yml"}, loader.SupportedExtensions())

	info, err := loader.Load(nil)
	assert.Nil(t, err)
	assert.Equal(t, &config.Info{
		ConfigFolder: tmpDir,
		LoadedFiles:  []string{defaultJson, developmentYml, localDevelopmentYaml},
	}, info)

	mock.AssertExpectationsForObjects(t, p0, p1)
}
