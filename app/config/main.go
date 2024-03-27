package config

import "flag"

const (
	DIR = "dir"
)

var CONFIG = map[string]string{}

func Config() {
	var dir string

	flag.StringVar(&dir, "dir", "/var/lib/redis", "Provide a directory where RDB files will be stored")

	flag.Parse()

	CONFIG[DIR] = dir
}
