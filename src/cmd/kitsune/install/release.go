package install

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var data *meta = nil

type meta struct {
	TagName string  `json:"tag_name"`
	Assets  []asset `json:"assets"`
}

func latest() *meta {
	if data != nil {
		return data
	}

	resp, err := http.Get("https://api.github.com/repos/kodflow/kitsune/releases/latest")
	if err != nil || resp.StatusCode != http.StatusOK {
		fmt.Println(err.Error())
		return nil
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil
	}

	return data
}
