package patreon

import (
	"fmt"

	"github.com/adamgoose/abots/lib/structure"
	"github.com/charmbracelet/log"
	"github.com/defval/di"
)

type Repo struct {
	di.Inject
	Log *log.Logger
}

// ForEachCampaignPostMedia iterates over all media for a campaign
func (r *Repo) ForEachCampaignPostMedia(tx *structure.Tx, id string, fn func(Entity[Campaign], Entity[Post], Entity[Media]) error) error {
	// Get the Campaign
	var campaign Entity[Campaign]
	if err := tx.GetStruct(PatreonScraperBucket, campaignID(id), &campaign); err != nil {
		return err
	}

	// Get Campaign Posts
	postIDs, err := tx.SMembers(PatreonScraperBucket, []byte(fmt.Sprintf("%s:posts", campaignID(id))))
	if err != nil {
		r.Log.Warn("No posts for campaign", "key", fmt.Sprintf("%s:posts", campaignID(id)))
		return err
	}

	for _, postID := range postIDs {
		// Get the Post
		var post Entity[Post]
		if err := tx.GetStruct(PatreonScraperBucket, string(postID), &post); err != nil {
			return err
		}

		// Get Post Media
		mediaIDs, err := tx.SMembers(PatreonScraperBucket, []byte(fmt.Sprintf("%s:media", postID)))
		if err != nil {
			r.Log.Debug("No media for post", "key", fmt.Sprintf("%s:media", postID))
			continue
		}

		for _, mediaID := range mediaIDs {
			// Get the media object
			var media Entity[Media]
			if err := tx.GetStruct(PatreonScraperBucket, string(mediaID), &media); err != nil {
				return err
			}

			fn(campaign, post, media)
		}
	}
	return nil
}

// MediaIsDownloaded returns true if the media has been downloaded
func (r *Repo) MediaIsDownloaded(tx *structure.Tx, id string) (bool, error) {
	return tx.SIsMember(PatreonScraperBucket, []byte("downloaded-media"), []byte(mediaID(id)))
}

// MediaIsOrganized returns true if the media has been organized
func (r *Repo) MediaIsOrganized(tx *structure.Tx, id string) (bool, error) {
	return tx.SIsMember(PatreonScraperBucket, []byte("organized-media"), []byte(mediaID(id)))
}

// MarkMediaAsDownloaded marks the media as downloaded
func (r *Repo) MarkMediaAsDownloaded(tx *structure.Tx, id string) error {
	return tx.SAdd(PatreonScraperBucket, []byte("downloaded-media"), []byte(mediaID(id)))
}

// MarkMediaAsOrganized marks the media as organized
func (r *Repo) MarkMediaAsOrganized(tx *structure.Tx, id string) error {
	return tx.SAdd(PatreonScraperBucket, []byte("organized-media"), []byte(mediaID(id)))
}
