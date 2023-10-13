package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "abots",
	Short: "Adam's Bag of Tricks",
}

func Execute() error {
	return rootCmd.Execute()
}
