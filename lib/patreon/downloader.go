package patreon

import (
	"fmt"

	"github.com/adamgoose/abots/lib/dl"
	"github.com/adamgoose/abots/lib/structure"
	"github.com/charmbracelet/log"
	"github.com/defval/di"
)

var PatreonDownloaderBucket = "downloads:patreon"

type Downloader struct {
	di.Inject
	Log   *log.Logger
	DB    *structure.DB
	Aria2 *dl.Aria2Downloader
}

func (d *Downloader) DownloadCampaign(id string) error {
	return d.DB.Update(func(tx *structure.Tx) error {
		// Get Campaign Posts
		postIDs, err := tx.SMembers(PatreonScraperBucket, []byte(fmt.Sprintf("%s:posts", campaignID(id))))
		if err != nil {
			d.Log.Warn("No posts for campaign", "key", fmt.Sprintf("%s:posts", campaignID(id)))
			return err
		}

		for _, postID := range postIDs {
			// Get Post Media
			mediaIDs, err := tx.SMembers(PatreonScraperBucket, []byte(fmt.Sprintf("%s:media", postID)))
			if err != nil {
				d.Log.Debug("No media for post", "key", fmt.Sprintf("%s:media", postID))
				continue
			}

			for _, mediaID := range mediaIDs {
				// See if it's already downloaded
				if downloaded, _ := tx.SIsMember(PatreonDownloaderBucket, []byte("downloaded-media"), mediaID); downloaded {
					continue
				}

				// Get the media object
				var media Entity[Media]
				if err := tx.GetStruct(PatreonScraperBucket, string(mediaID), &media); err != nil {
					return err
				}

				fd := fmt.Sprintf("%s/%s", campaignID(id), postID)
				fn := fmt.Sprintf("%s_%s", mediaID, media.Attributes.FileName)
				log.Info("Downloading Post Media", "file_dir", fd, "file_name", fn)
				if err := d.Aria2.Download(media.Attributes.DownloadURL, fd, fn); err != nil {
					return err
				}

				if err := tx.SAdd(PatreonDownloaderBucket, []byte("downloaded-media"), mediaID); err != nil {
					return err
				}
			}
		}

		return nil
	})
}
