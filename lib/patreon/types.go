package patreon

type PatreonResponse[T any] struct {
	Data     []PatreonEntity[T]     `json:"data"`
	Included []PatreonEntity[Media] `json:"included"`
}

type PatreonEntity[T any] struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes T
	// Relationships map[string]interface{}
}

type Post struct {
	CurrentUserCanView bool         `json:"current_user_can_view"`
	Title              string       `json:"title"`
	TeaserText         string       `json:"teaser_text"`
	PostMetadata       PostMetadata `json:"post_metadata"`

	Media []PatreonEntity[Media] `json:"-"`
}

type PostMetadata struct {
	ImageOrder []string `json:"image_order"`
}

type PostsQuery struct {
	CampaignID string `url:"filter[campaign_id],omitempty"`
	Sort       string `url:"sort,omitempty"`
	Include    string `url:"include,omitempty"`
}

type Media struct {
	FileName    string `json:"file_name"`
	MimeType    string `json:"mimetype"`
	DownloadURL string `json:"download_url"`
}
