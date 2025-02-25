package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

type rootCommand struct {
	cobra.Command
}

var rootCmd = &rootCommand{
	cobra.Command{
		Use:   "send",
		Short: "send a message",
		Run:   sendMessage,
	},
}

func (r *rootCommand) setup() {
	rootCmd.Flags().StringP("file", "f", "", "config file")

	rootCmd.MarkFlagRequired("file")
}

func sendMessage(cmd *cobra.Command, args []string) {
	fmt.Println("not implemented")
}
