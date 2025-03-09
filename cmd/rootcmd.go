package cmd

import (
	"github.com/keithpaterson/postal/logging"

	"github.com/spf13/cobra"
)

const (
	dryRunFlag  = "dry-run"
	debugFlag   = "debug"
	logfileFlag = "logfile"
	errfileFlag = "logerr"
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
	cmd.PersistentFlags().Bool(debugFlag, false, "turn on debugging output")
	cmd.PersistentFlags().String(logfileFlag, "", "log non-errors to this file")
	cmd.PersistentFlags().String(errfileFlag, "", "log errors to this file")

	return cmd
}

func setupLogging(cmd *cobra.Command) {
	var err error
	var debug bool
	if debug, err = cmd.Flags().GetBool(debugFlag); err != nil {
		panic("ERROR: failed to read debug flag!")
	}
	var logfile, errfile string
	if logfile, err = cmd.Flags().GetString(logfileFlag); err != nil {
		panic("ERROR: failed to read logfile flag")
	}
	if errfile, err = cmd.Flags().GetString(errfileFlag); err != nil {
		panic("ERROR: failed to read errfile flag")
	}

	logging.Setup(debug, logfile, errfile)
}
