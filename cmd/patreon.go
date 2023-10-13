package cmd

import (
	"fmt"

	"github.com/adamgoose/abots/lib"
	"github.com/adamgoose/abots/lib/patreon"
	"github.com/spf13/cobra"
)

var patreonCmd = &cobra.Command{
	Use:   "patreon",
	Short: "Interact with the Patreon API",
}

var patreonCampaignCmd = &cobra.Command{
	Use:   "campaign {id}",
	Short: "Lists posts and media file names for the given campaign id",
	Args:  cobra.ExactArgs(1),
	RunE: lib.RunE(func(args []string, p *patreon.API) error {
		r, err := p.GetVisiblePostsWithMedia(args[0])
		if err != nil {
			return err
		}

		for _, post := range r {
			fmt.Println(post.Attributes.Title)
			for _, image := range post.Attributes.Media {
				fmt.Printf("- %s\n", image.Attributes.FileName)
			}
		}

		return nil
	}),
}

func init() {
	patreonCmd.AddCommand(patreonCampaignCmd)
	rootCmd.AddCommand(patreonCmd)
}
