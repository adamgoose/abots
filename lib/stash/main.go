package stash

import (
	"github.com/shurcooL/graphql"
	"github.com/spf13/viper"
)

type Stash struct {
	*graphql.Client
}

func NewStash() *Stash {
	return &Stash{
		Client: graphql.NewClient(viper.GetString("stash_api_url"), nil),
	}
}
