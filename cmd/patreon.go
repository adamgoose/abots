package cmd

import (
	"fmt"

	"github.com/adamgoose/abots/lib"
	"github.com/adamgoose/abots/lib/patreon"
	"github.com/adamgoose/abots/lib/structure"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

var patreonCmd = &cobra.Command{
	Use:   "patreon",
	Short: "Interact with the Patreon API",
}

var patreonScrapeCmd = &cobra.Command{
	Use:   "scrape {campaign-id}",
	Short: "Scrapes the given campaign for posts and media",
	Args:  cobra.ExactArgs(1),
	RunE: lib.RunE(func(args []string, s *patreon.Scraper) error {
		return s.ScrapeCampaign(args[0])
	}),
}

var patreonCampaignCmd = &cobra.Command{
	Use:   "campaign {campaign-id}",
	Short: "Gets the given campaign",
	Args:  cobra.ExactArgs(1),
	RunE: lib.RunE(func(args []string, db *structure.DB) error {
		return db.View(func(tx *structure.Tx) error {
			bucket := "patreon"
			cid := fmt.Sprintf("campaign:%s", args[0])

			var ss patreon.ScrapeState
			if err := tx.GetStruct(bucket, fmt.Sprintf("%s:scrape_state", cid), &ss); err != nil {
				return err
			}

			spew.Dump(ss)

			mm, err := tx.SMembers(bucket, []byte(fmt.Sprintf("%s:posts", cid)))
			if err != nil {
				return err
			}

			for _, m := range mm {
				var post patreon.Entity[patreon.Post]
				if err := tx.GetStruct(bucket, string(m), &post); err != nil {
					return err
				}

				pmk, err := tx.SMembers(bucket, []byte(fmt.Sprintf("post:%s:media", post.ID)))
				if err != nil {
					return err
				}

				fmt.Printf("===\nPublished At: %s\nTitle: %s\nTeaser Text: %s\nMedia: %d\n",
					post.Attributes.PublishedAt,
					post.Attributes.Title,
					post.Attributes.TeaserText,
					len(pmk),
				)
			}

			return nil
		})
	}),
}

func init() {
	patreonCmd.AddCommand(patreonScrapeCmd)
	patreonCmd.AddCommand(patreonCampaignCmd)
	rootCmd.AddCommand(patreonCmd)
}
