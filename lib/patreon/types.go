package patreon

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/samber/lo"
)

type PatreonResponse[T any] struct {
	Data     []Entity[T]               `json:"data"`
	Included []Entity[json.RawMessage] `json:"included"`
	Meta     ResponseMeta              `json:"meta"`
}

type Entity[T any] struct {
	ID            string        `json:"id"`
	Type          string        `json:"type"`
	Attributes    T             `json:"attributes"`
	Relationships Relationships `json:"relationships"`
}

type ResponseMeta struct {
	Pagination struct {
		Cursors struct {
			Next string `json:"next"`
		} `json:"cursors"`
		Total int `json:"total"`
	} `json:"pagination"`
}

type Relationships map[string]Relationship

func (r Relationships) One(name string) (*RelationshipData, error) {
	rel, ok := r[name]
	if !ok {
		return nil, errors.New("not found")
	}

	return rel.One()
}

func (r Relationships) Many(name string) ([]RelationshipData, error) {
	rel, ok := r[name]
	if !ok {
		return nil, errors.New("not found")
	}

	return rel.Many()
}

type Relationship struct {
	Data json.RawMessage `json:"data"`
}

func (r Relationship) One() (*RelationshipData, error) {
	var d RelationshipData
	if err := json.Unmarshal(r.Data, &d); err != nil {
		return nil, err
	}

	return &d, nil
}

func (r Relationship) Many() ([]RelationshipData, error) {
	var d []RelationshipData
	if err := json.Unmarshal(r.Data, &d); err != nil {
		return nil, err
	}

	return d, nil
}

type RelationshipData struct {
	ID   string `json:"id"`
	Type string `json:"type"`
}

func FindIncludeByType[T any](inc []Entity[json.RawMessage], t string) (*Entity[T], error) {
	x, ok := lo.Find(inc, func(e Entity[json.RawMessage]) bool {
		return e.Type == t
	})
	if !ok {
		return nil, errors.New("not found")
	}

	return UnmarshalEntity[T](x)
}

func FindInclude[T any](inc []Entity[json.RawMessage], t, id string) (*Entity[T], error) {
	x, ok := lo.Find(inc, func(e Entity[json.RawMessage]) bool {
		return e.Type == t && e.ID == id
	})
	if !ok {
		return nil, errors.New("not found")
	}

	return UnmarshalEntity[T](x)
}

func FindRelationship[T any](inc []Entity[json.RawMessage], rd RelationshipData) (*Entity[T], error) {
	return FindInclude[T](inc, rd.Type, rd.ID)
}

func UnmarshalEntity[T any](e Entity[json.RawMessage]) (*Entity[T], error) {
	ne := Entity[T]{
		ID:   e.ID,
		Type: e.Type,
	}

	if err := json.Unmarshal(e.Attributes, &ne.Attributes); err != nil {
		return nil, err
	}

	return &ne, nil
}

type Campaign struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Post struct {
	CurrentUserCanView bool         `json:"current_user_can_view"`
	Title              string       `json:"title"`
	TeaserText         string       `json:"teaser_text"`
	PostMetadata       PostMetadata `json:"post_metadata"`
	PublishedAt        time.Time    `json:"published_at"`
}

type PostMetadata struct {
	ImageOrder []string `json:"image_order"`
}

type PostsQuery struct {
	CampaignID string  `url:"filter[campaign_id],omitempty"`
	Sort       string  `url:"sort,omitempty"`
	Include    string  `url:"include,omitempty"`
	Cursor     *string `url:"page[cursor],omitempty"`
}

type Media struct {
	FileName    string `json:"file_name"`
	MimeType    string `json:"mimetype"`
	DownloadURL string `json:"download_url"`
}
