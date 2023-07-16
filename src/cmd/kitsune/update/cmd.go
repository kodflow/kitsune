package update

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/kodmain/kitsune/src/internal/env"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:   "update",
	Short: "Update kitsune (" + getLatestVersion() + ")",
	Run: func(cmd *cobra.Command, args []string) {
		// Logic for status command
		fmt.Println("update to", Version.TagName)
		for _, asset := range Version.Assets {
			if !(strings.HasSuffix(asset.Name, ".md5") || strings.HasSuffix(asset.Name, ".sha1")) {
				var err error = nil
				if strings.Contains(asset.Name, "kitsune-"+runtime.GOOS+"-"+runtime.GOARCH) {
					err = asset.Download(env.PATH_BIN)
				} else if strings.Contains(asset.Name, runtime.GOOS+"-"+runtime.GOARCH) {
					err = asset.Download(env.PATH_SERVICES)
				}

				if err != nil {
					fmt.Println("Failed to do update", err.Error())
					break
				}
			}
		}
	},
}
