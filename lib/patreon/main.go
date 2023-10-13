package patreon

import (
	"fmt"

	"github.com/dghubble/sling"
	"github.com/samber/lo"
	"github.com/spf13/viper"
)

type API struct {
	*sling.Sling
}

func NewAPI() *API {
	return &API{
		Sling: sling.New().Base("https://www.patreon.com/api/").
			Add("Cookie", fmt.Sprintf(
				"session_id=%s;",
				viper.GetString("patreon_session_id"),
			)),
	}
}

func (d *API) GetVisiblePostsWithMedia(campaignID string) ([]PatreonEntity[Post], error) {
	pq := &PostsQuery{
		CampaignID: campaignID,
		Sort:       "-published_at",
		Include:    "images",
	}

	body := PatreonResponse[Post]{}
	_, err := d.New().Get("posts").QueryStruct(pq).ReceiveSuccess(&body)
	if err != nil {
		return nil, err
	}

	images := lo.KeyBy(
		lo.Filter(body.Included, func(e PatreonEntity[Media], i int) bool {
			return e.Type == "media"
		}),
		func(e PatreonEntity[Media]) string {
			return e.ID
		},
	)

	posts := []PatreonEntity[Post]{}
	for _, post := range body.Data {
		if post.Attributes.CurrentUserCanView {
			m := []PatreonEntity[Media]{}
			for _, imageID := range post.Attributes.PostMetadata.ImageOrder {
				image := images[imageID]
				if image.Attributes.DownloadURL == "" {
					continue
				}
				m = append(m, image)
			}
			post.Attributes.Media = m
			posts = append(posts, post)
		}
	}

	return posts, nil
}
