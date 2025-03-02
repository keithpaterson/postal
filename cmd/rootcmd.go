package cmd

import "github.com/spf13/cobra"

const (
	dryRunFlag = "dry-run"
)

func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "postal",
		Short: "Postal is a CLI replacement for Postman.",
		Long: `Postal allows you to compose and send HTTP requests from the command line
  by concatenating configurations, injecting command-line arguments, environment variables,
  and more.`,
	}

	cmd.PersistentFlags().BoolP(dryRunFlag, "d", false, "dry run will not perform the operation")

	return cmd
}
