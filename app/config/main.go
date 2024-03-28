package config

import (
	"flag"
	"fmt"
	"os"
)

const (
	DIR = "dir"
)

var CONFIG = map[string]string{}

func Config() {
	var dir string

  homePath := os.Getenv("HOME")
  defaultDir := fmt.Sprintf("%s/http-go-tmp", homePath)
	flag.StringVar(&dir, "directory", defaultDir, "Provide a directory where the endpoint files will read files")

	flag.Parse()

	CONFIG[DIR] = dir
}
