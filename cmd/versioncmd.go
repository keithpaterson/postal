package cmd

import (
	"fmt"
	"runtime/debug"

	"github.com/spf13/cobra"
)

func init() {
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" {
			appVersion = info.Main.Version
		}
		return
	}
	// panic?
	panic("failed to read build information")
}

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show the current version",
		Run:   showVersion,
	}
}

var appVersion = "(devel)"

func showVersion(_ *cobra.Command, _ []string) {
	fmt.Println(appVersion)
}
