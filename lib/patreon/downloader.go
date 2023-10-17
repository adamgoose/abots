package patreon

import (
	"fmt"

	"github.com/adamgoose/abots/lib/dl"
	"github.com/adamgoose/abots/lib/structure"
	"github.com/charmbracelet/log"
	"github.com/defval/di"
)

type Downloader struct {
	di.Inject
	Log   *log.Logger
	DB    *structure.DB
	Repo  *Repo
	Aria2 *dl.Aria2Downloader
}

func (d *Downloader) DownloadCampaign(id string) error {
	return d.DB.Update(func(tx *structure.Tx) error {
		return d.Repo.ForEachCampaignPostMedia(tx, id, func(campaign Entity[Campaign], post Entity[Post], media Entity[Media]) error {
			// See if it's already downloaded
			if downloaded, _ := d.Repo.MediaIsDownloaded(tx, media.ID); downloaded {
				return nil
			}

			// Download the media
			fd := fmt.Sprintf("%s/%s", campaignID(id), postID(post.ID))
			fn := fmt.Sprintf("%s_%s", mediaID(media.ID), media.Attributes.FileName)
			log.Info("Downloading Post Media", "file_dir", fd, "file_name", fn)
			if err := d.Aria2.Download(media.Attributes.DownloadURL, fd, fn); err != nil {
				return err
			}

			// Mark it as downloaded
			return d.Repo.MarkMediaAsDownloaded(tx, media.ID)
		})
	})
}
