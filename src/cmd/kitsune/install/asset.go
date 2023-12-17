package install

import (
	"fmt"
	"io"
	"net/http"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/kodflow/kitsune/src/internal/kernel/storages/fs"
)

type asset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func (a *asset) Download(destination string) error {

	if err := fs.CreateDirectory(destination); err != nil {
		return err
	}

	aNameSplit := strings.SplitN(a.Name, "-", 2)
	binaryPath := filepath.Join(destination, aNameSplit[0])

	root, err := user.Lookup("root")
	if err != nil {
		return err
	}

	resp, err := http.Get(a.BrowserDownloadURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := fs.CreateFile(binaryPath, &fs.Options{
		User:  root,
		Perms: 0755,
	})

	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)

	if err == nil {
		fmt.Printf("Download %s.\n", aNameSplit[0])
	}

	return err
}
