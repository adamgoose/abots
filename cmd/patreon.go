package cmd

import (
	"fmt"
	"time"

	"github.com/adamgoose/abots/lib"
	"github.com/adamgoose/abots/lib/patreon"
	"github.com/nutsdb/nutsdb"
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
	RunE: lib.RunE(func(args []string, p *patreon.API, n *nutsdb.DB) error {
		r, err := p.GetVisiblePostsWithMedia(args[0])
		if err != nil {
			return err
		}

		return n.Update(func(tx *nutsdb.Tx) error {
			bucket := "patreon"

			if err := tx.SAdd(bucket, []byte("campaigns"), []byte(fmt.Sprintf("campaign:%s", args[0]))); err != nil {
				return err
			}

			kvp := map[string]string{
				fmt.Sprintf("campaign:%s:scraped_at", args[0]):          time.Now().Format("2006-01-02 15:04:05"),
				fmt.Sprintf("campaign:%s:oldest_published_at", args[0]): r[len(r)-1].Attributes.PublishedAt.Format("2006-01-02 15:04:05"),
			}
			for k, v := range kvp {
				if err := tx.Put(bucket, []byte(k), []byte(v), 0); err != nil {
					return err
				}
			}

			for _, post := range r {
				if err := post.Attributes.Persist(tx, bucket, post.ID, args[0]); err != nil {
					return err
				}

				for _, media := range post.Attributes.Media {
					if err := media.Attributes.Persist(tx, bucket, media.ID, post.ID); err != nil {
						return err
					}
				}
			}

			return nil
		})
	}),
}

func init() {
	patreonCmd.AddCommand(patreonCampaignCmd)
	rootCmd.AddCommand(patreonCmd)
}
