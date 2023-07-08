package update

import (
	"io"
	"net/http"
	"os/user"
	"path/filepath"

	"github.com/kodmain/kitsune/src/internal/storages/fs"
)

type asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func (a *asset) Download(destination string) error {
	resp, err := http.Get(a.BrowserDownloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	wheel, err := user.LookupGroup("wheel")
	if err != nil {
		return err
	}

	root, err := user.Lookup("root")
	if err != nil {
		return err
	}

	out, err := fs.CreateFile(filepath.Join(destination, a.Name), &fs.CreateOption{
		User:  root,
		Group: wheel,
		Perms: 0755,
	})

	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}
