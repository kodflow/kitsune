package install

import (
	"strconv"
	"strings"
)

/*
var Version *release = nil

func getVersion() {

	if Version == nil {

	}

}
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

*/

func compareVersions(version1, version2 string) bool {
	if version1 == "" || version2 == "" {
		return true
	}

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

	return len(v1Nums) < len(v2Nums)
}
