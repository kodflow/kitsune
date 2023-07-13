package update

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/kodmain/kitsune/src/internal/env"
	"github.com/spf13/cobra"
)

var Version = getLatestRelease()

func getLatestVersion() string {
	if env.BUILD_VERSION == "" {
		return color.YellowString("You are on a local build")
	}

	if Version.TagName == "" {
		return color.RedString("Unable to compare versions.")
	}

	if env.BUILD_VERSION == Version.TagName {
		return color.GreenString("You are on the latest version.")
	}

	return fmt.Sprintf("From %s to %s", color.RedString(env.BUILD_VERSION), color.GreenString(Version.TagName))
}

func compareVersions(version1, version2 string) bool {
	v1Parts := strings.Split(strings.TrimPrefix(version1, "v"), ".")
	v2Parts := strings.Split(strings.TrimPrefix(version2, "v"), ".")

	v1Nums := make([]int, len(v1Parts))
	v2Nums := make([]int, len(v2Parts))
	for i, part := range v1Parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			return false
		}
		v1Nums[i] = num
	}
	for i, part := range v2Parts {
		num, err := strconv.Atoi(part)
		if err != nil {
			return false
		}
		v2Nums[i] = num
	}

	for i := 0; i < len(v1Nums) && i < len(v2Nums); i++ {
		if v1Nums[i] > v2Nums[i] {
			return true
		} else if v1Nums[i] < v2Nums[i] {
			return false
		}
	}

	return len(v1Nums) > len(v2Nums)
}

func ShooldUpdate() bool {
	return compareVersions(env.BUILD_VERSION, Version.TagName)
}

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
