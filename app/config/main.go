package config

import "flag"

const (
	DIR = "dir"
)

var CONFIG = map[string]string{}

func Config() {
	var dir string

	flag.StringVar(&dir, "directory", "/var/lib/redis", "Provide a directory where the endpoint files will read files")

	flag.Parse()

	CONFIG[DIR] = dir
}
