package dl

import (
	"github.com/charmbracelet/log"
	"github.com/defval/di"
)

type NopDownloader struct {
	di.Inject
	Log *log.Logger
}

func (d *NopDownloader) Download(url, dir, out string) error {
	d.Log.Info("NopDownloader: Downloading", "url", url, "dir", dir, "out", out)
	return nil
}
