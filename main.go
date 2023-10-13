package main

import (
	"log"

	"github.com/adamgoose/abots/cmd"
	"github.com/adamgoose/abots/lib"
	"github.com/adamgoose/abots/lib/dl"
	"github.com/adamgoose/abots/lib/patreon"
	"github.com/defval/di"
	"github.com/nutsdb/nutsdb"
	"github.com/spf13/viper"
)

func main() {
	if err := lib.App.Apply(
		di.Provide(patreon.NewAPI),
		di.Provide(dl.NewAria2Downloader),
		di.Provide(func() (*nutsdb.DB, error) {
			return nutsdb.Open(
				nutsdb.DefaultOptions,
				nutsdb.WithDir(viper.GetString("db_path")),
			)
		}),
	); err != nil {
		log.Fatal(err)
	}

	cmd.Execute()
}

func init() {
	viper.SetEnvPrefix("abots")
	viper.AutomaticEnv()

	viper.SetDefault("db_path", "./.db")
}
