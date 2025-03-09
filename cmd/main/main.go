package main

import (
	"github.com/spf13/cobra"

	"github.com/keithpaterson/postal/cmd"
	"github.com/keithpaterson/postal/logging"
)

func main() {
	defer logging.Teardown()

	rootCmd := setupCli()

	if err := rootCmd.Execute(); err != nil {
		panic(err)
	}
}

func setupCli() *cobra.Command {
	rootCmd := cmd.NewRootCommand()
	rootCmd.AddCommand(cmd.NewVersionCmd())
	rootCmd.AddCommand(cmd.NewSendCommand())
	return rootCmd
}
