package update

import (
	"encoding/json"
	"io"
	"net/http"
)

type release struct {
	TagName string  `json:"tag_name"`
	Assets  []asset `json:"assets"`
}

func getLatestRelease() *release {

	resp, err := http.Get("https://api.github.com/repos/kodmain/kitsune/releases/latest")
	if err != nil || resp.StatusCode != http.StatusOK {
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	var version *release
	err = json.Unmarshal(body, &version)
	if err != nil {
		return nil
	}

	return version
}
