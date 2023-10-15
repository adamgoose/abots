package cmd

import (
	"github.com/spf13/cobra"
)

var Version = "development"

var rootCmd = &cobra.Command{
	Use:     "abots",
	Short:   "Adam's Bag of Tricks",
	Version: Version,
}

func Execute() error {
	return rootCmd.Execute()
}
