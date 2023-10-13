package dl

import (
	"github.com/siku2/arigo"
	"github.com/spf13/viper"
)

type Aria2Downloader struct {
	rpc  *arigo.Client
	root string
}

func NewAria2Downloader() (*Aria2Downloader, error) {
	ag, err := arigo.Dial(
		viper.GetString("aria2_rpc_url"),
		viper.GetString("aria2_rpc_secret"),
	)
	if err != nil {
		return nil, err
	}

	return &Aria2Downloader{
		rpc:  &ag,
		root: viper.GetString("aria2_root"),
	}, nil
}

func (a *Aria2Downloader) Download(url, filename string) error {
	_, err := a.rpc.AddURI([]string{url}, &arigo.Options{
		Out: filename,
		Dir: a.root,
	})
	return err
}
