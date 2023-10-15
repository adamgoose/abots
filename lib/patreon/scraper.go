package patreon

import (
	"fmt"
	"time"

	"github.com/adamgoose/abots/lib/structure"
	"github.com/charmbracelet/log"
	"github.com/defval/di"
)

var PatreonScraperBucket = "patreon"

type Scraper struct {
	di.Inject
	Log *log.Logger
	API *API
	DB  *structure.DB

	r  *PatreonResponse[Post]
	tx *structure.Tx
}

type ScrapeState struct {
	ScrapedAt         time.Time
	LatestPublishedAt *time.Time
}

func (s *Scraper) ScrapeCampaign(id string) error {
	var scraped int
	var cursor *string
	var latestPublishedAt *time.Time

	ss := ScrapeState{}
	if err := s.DB.View(func(tx *structure.Tx) error {
		return tx.GetStruct(PatreonScraperBucket, fmt.Sprintf("%s:scrape_state", campaignID(id)), &ss)
	}); err == nil {
		s.Log.Info("Found scrape state",
			"scraped_at", ss.ScrapedAt,
			"latest_published_at", ss.LatestPublishedAt,
		)
	}

	// Scrape the API
	for {
		r, err := s.API.GetPosts(id, cursor)
		if err != nil {
			return err
		}

		if latestPublishedAt == nil {
			latestPublishedAt = &r.Data[0].Attributes.PublishedAt

			if err := s.DB.Update(func(tx *structure.Tx) error {
				return s.handleCampaign(*r, tx)
			}); err != nil {
				return err
			}
		}

		// consider breaking because of the previous sync state
		// if the first post was published before the latest scrape state
		if ss.LatestPublishedAt != nil && r.Data[0].Attributes.PublishedAt.Before(*ss.LatestPublishedAt) {
			s.Log.Info("Skipping remaining posts, as they're already scraped.",
				"latest_published_at", ss.LatestPublishedAt,
				"post_published_at", r.Data[0].Attributes.PublishedAt,
			)
			break
		}

		if err := s.handleScrapedPosts(*r); err != nil {
			return err
		}
		scraped += len(r.Data)
		s.Log.Info("Scraped a page of posts", "scraped", scraped, "total", r.Meta.Pagination.Total, "pageSize", len(r.Data))

		if scraped >= r.Meta.Pagination.Total {
			break
		}

		cursor = &r.Meta.Pagination.Cursors.Next
	}

	// handle scrape state
	ss.ScrapedAt = time.Now()
	ss.LatestPublishedAt = latestPublishedAt
	return s.DB.Update(func(tx *structure.Tx) error {
		return tx.PutStruct(PatreonScraperBucket, fmt.Sprintf("%s:scrape_state", campaignID(id)), ss)
	})
}

func (s *Scraper) handleScrapedPosts(r PatreonResponse[Post]) error {
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
		return s.handlePosts(r, tx, *c)
	})
}

func (s *Scraper) handleCampaign(r PatreonResponse[Post], tx *structure.Tx) error {
	// Fetch Campaign
	ctid, err := r.Data[0].Relationships.One("campaign")
	if err != nil {
		return nil
	}
	c, err := FindRelationship[Campaign](r.Included, *ctid)
	if err != nil {
		return err
	}

	// Add to Campaigns Set
	if err := tx.SAdd(PatreonScraperBucket, []byte("campaigns"), []byte(campaignID(c.ID))); err != nil {
		return err
	}

	// Persist the Campaign
	if err := tx.PutStruct(PatreonScraperBucket, campaignID(c.ID), *c); err != nil {
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
			[]byte(fmt.Sprintf("%s:posts", campaignID(c.ID))),
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
			[]byte(fmt.Sprintf("%s:media", postID(p.ID))),
			[]byte(mediaID(m.ID)),
		); err != nil {
			return err
		}

		// Persist the Media
		if err := tx.PutStruct(PatreonScraperBucket, mediaID(m.ID), m); err != nil {
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
