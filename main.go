package main

import (
	"os"

	"github.com/adamgoose/abots/cmd"
	"github.com/adamgoose/abots/lib"
	"github.com/adamgoose/abots/lib/dl"
	"github.com/adamgoose/abots/lib/patreon"
	"github.com/adamgoose/abots/lib/stash"
	"github.com/charmbracelet/log"
	"github.com/defval/di"
	"github.com/muesli/termenv"
	"github.com/nutsdb/nutsdb"
	"github.com/spf13/viper"
)

func main() {
	if err := lib.App.Apply(
		// Configure Logging
		di.Provide(func() *log.Logger {
			l := log.New(os.Stdout)
			l.SetLevel(log.ParseLevel(viper.GetString("log_level")))
			l.SetTimeFormat("3:04.05pm")
			if viper.GetBool("log_color") {
				l.SetColorProfile(termenv.TrueColor)
			}

			return l
		}),

		// Configure Storage
		di.Provide(func(l *log.Logger) (*nutsdb.DB, func(), error) {
			db, err := nutsdb.Open(
				nutsdb.DefaultOptions,
				nutsdb.WithDir(viper.GetString("db_path")),
			)

			cleanup := func() {
				if err := db.Close(); err != nil {
					l.Fatal(err)
				}
			}

			return db, cleanup, err
		}),

		// Configure APIs
		di.Provide(patreon.NewAuthenticator),
		di.Provide(patreon.NewAPI),
		di.Provide(dl.NewAria2Downloader),
		di.Provide(stash.NewStash),
	); err != nil {
		log.Fatal(err)
	}

	cmd.Execute()
	lib.App.Cleanup()
}

func init() {
	viper.SetEnvPrefix("abots")
	viper.AutomaticEnv()

	viper.SetDefault("db_path", "./.db")
}
