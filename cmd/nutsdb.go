package cmd

import (
	"fmt"

	"github.com/adamgoose/abots/lib"
	"github.com/adamgoose/abots/lib/structure"
	"github.com/nutsdb/nutsdb"
	"github.com/spf13/cobra"
)

var nutsdbCmd = &cobra.Command{
	Use:     "nutsdb",
	Short:   "Interact with the NutsDB database",
	Aliases: []string{"nuts", "db"},
}

var nutsdbKeysCmd = &cobra.Command{
	Use:   "keys {bucket} [prefix]",
	Short: "List all keys in the given bucket",
	Args:  cobra.RangeArgs(1, 2),
	RunE: lib.RunE(func(args []string, db *nutsdb.DB) error {
		if len(args) == 1 {
			args = append(args, "")
		}

		return db.View(func(tx *nutsdb.Tx) error {
			entries, err := tx.PrefixScan(args[0], []byte(args[1]), 0, 999)
			if err != nil {
				return err
			}

			for _, entry := range entries {
				fmt.Printf("%s\n", entry.Key)
			}

			return nil
		})
	}),
}

var nutsdbSetsCmd = &cobra.Command{
	Use:   "sets {bucket}",
	Short: "List all sets in the given bucket",
	Args:  cobra.ExactArgs(1),
	RunE: lib.RunE(func(args []string, db *nutsdb.DB) error {
		return db.View(func(tx *nutsdb.Tx) error {
			return tx.SKeys(args[0], "*", func(key string) bool {
				fmt.Println(key)
				return true
			})
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
	RunE: lib.RunE(func(args []string, db *structure.DB) error {
		return db.View(func(tx *structure.Tx) error {
			val, err := tx.Get(args[0], []byte(args[1]))
			if err != nil {
				return err
			}

			fmt.Println(val)
			return nil
		})
	}),
}

func init() {
	nutsdbCmd.AddCommand(nutsdbKeysCmd)
	nutsdbCmd.AddCommand(nutsdbSetsCmd)
	nutsdbCmd.AddCommand(nutsdbSetCmd)
	nutsdbCmd.AddCommand(nutsdbGetCmd)
	rootCmd.AddCommand(nutsdbCmd)
}
