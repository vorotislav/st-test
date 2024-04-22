// Package util содержит функцию по разбору флагов, функцию вывода версии, а так же хранит переменные версии, сборки и даты.
package util

import (
	"flag"
	"os"
)

// ParseFlags разбирает переданные флаги.
func ParseFlags() string {
	configPath := flag.String("config", "", "Path to configuration file")
	version := flag.Bool("version", false, "Show version")

	flag.Parse()

	if *version {
		printVersion()
		os.Exit(0)
	}

	return *configPath
}
