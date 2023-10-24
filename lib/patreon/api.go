package patreon

import (
	"fmt"

	"github.com/dghubble/sling"
)

type API struct {
	*sling.Sling
}

func NewAPI(a Authenticator) (*API, error) {
	sid, err := a.GetSessionID()
	if err != nil {
		return nil, err
	}

	return &API{
		Sling: sling.New().Base("https://www.patreon.com/api/").
			Add("Cookie", fmt.Sprintf("session_id=%s;", sid)),
	}, nil
}

func (d *API) GetUser() (*PatreonResponseOne[User], error) {
	body := PatreonResponseOne[User]{}

	if _, err := d.New().
		Get("current_user").
		QueryStruct(&CurrentUserQuery{
			UserFields: "full_name,email",
		}).
		ReceiveSuccess(&body); err != nil {
		return nil, err
	}

	return &body, nil
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
