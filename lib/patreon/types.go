package patreon

import (
	"fmt"
	"time"

	"github.com/nutsdb/nutsdb"
)

type PatreonResponse[T any] struct {
	Data     []Entity[T]     `json:"data"`
	Included []Entity[Media] `json:"included"`
}

type Entity[T any] struct {
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
	PublishedAt        time.Time    `json:"published_at"`

	Media []Entity[Media] `json:"-"`
}

func (p Post) Persist(tx *nutsdb.Tx, bucket, postID, campaignID string) error {
	if err := tx.SAdd(bucket,
		[]byte(fmt.Sprintf("campaign:%s:posts", campaignID)),
		[]byte(fmt.Sprintf("post:%s", postID)),
	); err != nil {
		return err
	}

	kvp := map[string]string{
		fmt.Sprintf("post:%s:title", postID):       p.Title,
		fmt.Sprintf("post:%s:teaser_text", postID): p.TeaserText,
	}
	for k, v := range kvp {
		if err := tx.Put(bucket, []byte(k), []byte(v), 0); err != nil {
			return err
		}
	}

	return nil
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

func (m Media) Persist(tx *nutsdb.Tx, bucket, mediaID, postID string) error {
	tid := fmt.Sprintf("media:%s", mediaID)

	if err := tx.SAdd(bucket,
		[]byte(fmt.Sprintf("post:%s:media", postID)),
		[]byte(tid),
	); err != nil {
		return err
	}

	kvp := map[string]string{
		fmt.Sprintf("%s:file_name", tid):    m.FileName,
		fmt.Sprintf("%s:mimetype", tid):     m.MimeType,
		fmt.Sprintf("%s:download_url", tid): m.DownloadURL,
	}
	for k, v := range kvp {
		if err := tx.Put(bucket, []byte(k), []byte(v), 0); err != nil {
			return err
		}
	}

	return nil
}
