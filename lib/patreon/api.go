package patreon

import (
	"fmt"

	"github.com/dghubble/sling"
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

func (d *API) GetPosts(campaignID string, cursor *string) (*PatreonResponse[Post], error) {
	body := PatreonResponse[Post]{}

	if _, err := d.New().
		Get("posts").
		QueryStruct(&PostsQuery{
			CampaignID: campaignID,
			Sort:       "-published_at",
			Include:    "campaign,images",
			Cursor:     cursor,
		}).
		ReceiveSuccess(&body); err != nil {
		return nil, err
	}

	return &body, nil
}
