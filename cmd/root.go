package cmd

import (
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "supreme-robot",
		Short: "The Supreme Robot",
		//TODO Long:  ``,
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
