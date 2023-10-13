package cmd

import (
	"fmt"

	"github.com/adamgoose/abots/lib"
	"github.com/nutsdb/nutsdb"
	"github.com/spf13/cobra"
)

var nutsdbCmd = &cobra.Command{
	Use:     "nutsdb",
	Short:   "Interact with the NutsDB database",
	Aliases: []string{"nuts", "db"},
}

var nutsdbKeysCmd = &cobra.Command{
	Use:   "keys {bucket}",
	Short: "List all keys in the given bucket",
	Args:  cobra.ExactArgs(1),
	RunE: lib.RunE(func(args []string, db *nutsdb.DB) error {
		return db.View(func(tx *nutsdb.Tx) error {
			entries, err := tx.GetAll(args[0])
			if err != nil {
				return err
			}

			for _, entry := range entries {
				fmt.Println(string(entry.Key))
			}

			return nil
		})
	}),
}

var nutsdbSetCmd = &cobra.Command{
	Use:   "set {bucket} {key} {value}",
	Short: "Set a key/value pair in the given bucket",
	Args:  cobra.ExactArgs(3),
	RunE: lib.RunE(func(args []string, db *nutsdb.DB) error {
		return db.Update(func(tx *nutsdb.Tx) error {
			return tx.Put(args[0], []byte(args[1]), []byte(args[2]), 0)
		})
	}),
}

var nutsdbGetCmd = &cobra.Command{
	Use:   "get {bucket} {key}",
	Short: "Get the value for a given key in the given bucket",
	Args:  cobra.ExactArgs(2),
	RunE: lib.RunE(func(args []string, db *nutsdb.DB) error {
		return db.View(func(tx *nutsdb.Tx) error {
			entry, err := tx.Get(args[0], []byte(args[1]))
			if err != nil {
				return err
			}

			fmt.Println(string(entry.Value))

			return nil
		})
	}),
}

func init() {
	nutsdbCmd.AddCommand(nutsdbKeysCmd)
	nutsdbCmd.AddCommand(nutsdbSetCmd)
	nutsdbCmd.AddCommand(nutsdbGetCmd)
	rootCmd.AddCommand(nutsdbCmd)
}
