package main

import (
	"log"

	"github.com/adamgoose/abots/cmd"
	"github.com/adamgoose/abots/lib"
	"github.com/adamgoose/abots/lib/dl"
	"github.com/adamgoose/abots/lib/patreon"
	"github.com/defval/di"
	"github.com/spf13/viper"
)

func main() {
	if err := lib.App.Apply(
		di.Provide(patreon.NewAPI),
		di.Provide(dl.NewAria2Downloader),
	); err != nil {
		log.Fatal(err)
	}

	cmd.Execute()
}

func init() {
	viper.SetEnvPrefix("abots")
	viper.AutomaticEnv()
}
