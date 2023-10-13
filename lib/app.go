package lib

import (
	"context"

	"github.com/defval/di"
	"github.com/spf13/cobra"
)

var App *di.Container

func init() {
	App, _ = di.New()
}

func RunE(runE interface{}) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		App.ProvideValue(cmd)
		App.ProvideValue(args)
		App.ProvideValue(cmd.Context(), di.As(new(context.Context)))
		return App.Invoke(runE)
	}
}
