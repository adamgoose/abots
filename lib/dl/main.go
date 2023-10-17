package dl

type Downloader interface {
	Download(url, dir, out string) error
}
