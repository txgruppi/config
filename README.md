[![GoDoc](https://img.shields.io/badge/godoc-reference-546FA7.svg?style=flat-square)](https://godoc.org/github.com/txgruppi/config)
![Codeship](https://img.shields.io/codeship/55287ea0-85e7-0137-8460-5ef343a80245.svg?style=flat-square)
[![Codecov](https://img.shields.io/codecov/c/github/txgruppi/config.svg?style=flat-square)](https://codecov.io/github/txgruppi/run)
[![Go Report Card](https://img.shields.io/badge/go_report-A+-brightgreen.svg?style=flat-square)](https://goreportcard.com/report/github.com/txgruppi/config)

# config

Quick and easy way to load config files based on a simple set of rules.

Project inspired by https://github.com/lorenwest/node-config

## Important stuff

### Supported files

Before you can load any file you must register parsers using `Loader.RegisterParser`.

Each parser has a list of supported extensions that will be used to find files to load.

### Config folder

By default the load will try to find the files based on the environment variable name given to it (defaults to `CONFIG_DIR`). If the variable name is empty or the variable value is empty, it will look for files in `./config`.

### File load order

```
default.{ext}
{deployment}.{ext}
{hostname}.{ext}
{hostname}-{deployment}.{ext}
local.{ext}
local-{deployment}.{ext}
```

Where

- `{ext}` is one of the registered extensions.
- `{deployment}` is the deployment name, from the `$ENV` environment variable. (No default value, ignored if empty)
- `{hostname}` is the value returned from `os.Hostname()` with no changes. (No default value, ignored if empty)

## Installation

```
go get -u github.com/txgruppi/config
```

## Example

```
package main

import (
	"fmt"
	"log"

	"github.com/txgruppi/config"
	"github.com/txgruppi/config/parsers/json"
)

type Config struct {
	Server struct {
		Bind string `json:"bind"`
		Port int    `json:"port"`
	} `json:"server"`
}

func main() {
	loader := NewLoader()
	if err := loader.RegisterParser(json.NewParser()); err != nil {
		log.Fatal(err)
	}
	var config Config
	info, err := loader.Load(&config)
	if err != nil {
		log.Fatal(err)
	}
  fmt.Printf("Looked for files in: %s\n", info.ConfigFolder)
	fmt.Printf("Loaded files: %v\n", info.LoadedFiles)
	fmt.Printf("Loaded config: %v\n", config)
}
```
