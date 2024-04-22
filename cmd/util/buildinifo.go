package util

import (
	"flag"
	"fmt"
)

var (
	buildVersion = "N/A" //nolint:gochecknoglobals
	buildDate    = "N/A" //nolint:gochecknoglobals
	buildCommit  = "N/A" //nolint:gochecknoglobals
)

func printVersion() {
	_, _ = fmt.Fprintf(
		flag.CommandLine.Output(),
		Version()+"\n",
	)
}

// Version выводит версию приложения в стандартный вывод.
func Version() string {
	return fmt.Sprintf(
		"Test %s %s (%s)",
		buildVersion,
		buildDate,
		buildCommit,
	)
}
