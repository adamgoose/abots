package patreon

import (
	"fmt"
	"time"

	"github.com/adamgoose/abots/lib/structure"
	"github.com/defval/di"
)

var PatreonScraperBucket = "patreon"

type Scraper struct {
	di.Inject
	API *API
	DB  *structure.DB

	r  *PatreonResponse[Post]
	tx *structure.Tx
}

type ScrapeState struct {
	ScrapedAt         time.Time
	OldestPublishedAt time.Time
}

func (s *Scraper) ScrapeCampaign(id string) error {
	// Scrape the API
	r, err := s.API.GetPosts(id)
	if err != nil {
		return err
	}

	// Get the Campaign
	ctid, err := r.Data[0].Relationships.One("campaign")
	if err != nil {
		return nil
	}
	c, err := FindRelationship[Campaign](r.Included, *ctid)
	if err != nil {
		return err
	}

	return s.DB.Update(func(tx *structure.Tx) error {
		if err := s.handleScrapeState(*r, tx, *c); err != nil {
			return err
		}

		if err := s.handleCampaign(*r, tx, *c); err != nil {
			return err
		}

		if err := s.handlePosts(*r, tx, *c); err != nil {
			return err
		}

		return nil
	})
}

func (s *Scraper) handleScrapeState(r PatreonResponse[Post], tx *structure.Tx, c Entity[Campaign]) error {
	ctid, err := r.Data[0].Relationships.One("campaign")
	if err != nil {
		return nil
	}

	return tx.PutStruct(PatreonScraperBucket, fmt.Sprintf("campaign:%s:scrape_state", ctid.ID), ScrapeState{
		ScrapedAt:         time.Now(),
		OldestPublishedAt: r.Data[len(r.Data)-1].Attributes.PublishedAt,
	})
}

func (s *Scraper) handleCampaign(r PatreonResponse[Post], tx *structure.Tx, c Entity[Campaign]) error {
	// Add to Campaigns Set
	if err := tx.SAdd(PatreonScraperBucket, []byte("campaigns"), []byte(campaignID(c.ID))); err != nil {
		return err
	}

	// Persist the Campaign
	if err := tx.PutStruct(PatreonScraperBucket, campaignID(c.ID), c); err != nil {
		return err
	}

	return nil
}

func (s *Scraper) handlePosts(r PatreonResponse[Post], tx *structure.Tx, c Entity[Campaign]) error {
	for _, post := range r.Data {
		if !post.Attributes.CurrentUserCanView {
			continue
		}

		// Add to Posts set
		if err := tx.SAdd(PatreonScraperBucket,
			[]byte("posts"),
			[]byte(postID(post.ID)),
		); err != nil {
			return err
		}

		// Add to Campaign Posts Set
		if err := tx.SAdd(PatreonScraperBucket,
			[]byte(fmt.Sprintf("%s:posts", c.ID)),
			[]byte(postID(post.ID)),
		); err != nil {
			return err
		}

		// Persist the Post
		if err := tx.PutStruct(PatreonScraperBucket, postID(post.ID), post); err != nil {
			return err
		}

		// Handle Post Media
		if err := s.handlePostMedia(r, tx, post); err != nil {
			return err
		}

	}

	return nil
}

func (s *Scraper) handlePostMedia(r PatreonResponse[Post], tx *structure.Tx, p Entity[Post]) error {
	for _, id := range p.Attributes.PostMetadata.ImageOrder {
		m, err := FindInclude[Media](r.Included, "media", id)
		if err != nil {
			return err
		}

		if m.Attributes.DownloadURL == "" {
			continue
		}

		// Add to Media Set
		if err := tx.SAdd(PatreonScraperBucket,
			[]byte("media"),
			[]byte(mediaID(m.ID)),
		); err != nil {
			return err
		}

		// Add to Post Media Set
		if err := tx.SAdd(PatreonScraperBucket,
			[]byte(fmt.Sprintf("%s:media", postID(id))),
			[]byte(mediaID(m.ID)),
		); err != nil {
			return err
		}

		// Persist the Media
		if err := tx.PutStruct(PatreonScraperBucket, fmt.Sprintf("media:%s", m.ID), m); err != nil {
			return err
		}
	}

	return nil
}

func postID(id string) string {
	return fmt.Sprintf("post:%s", id)
}

func mediaID(id string) string {
	return fmt.Sprintf("media:%s", id)
}

func campaignID(id string) string {
	return fmt.Sprintf("campaign:%s", id)
}
