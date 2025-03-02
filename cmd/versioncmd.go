package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show the current version",
		Run:   showVersion,
	}
}

var appVersion = "1.0.0"

func showVersion(_ *cobra.Command, _ []string) {
	fmt.Println(appVersion)
}
