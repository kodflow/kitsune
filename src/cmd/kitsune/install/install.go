package install

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/kodmain/kitsune/src/internal/env"
	"github.com/kodmain/kitsune/src/internal/observability/logger"
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:                   "install <service1> <service2> ... (all by default)",
	Short:                 "Install kitsune",
	Long:                  "Install all kitsune services or specifie what you want",
	DisableFlagParsing:    true,
	DisableAutoGenTag:     true,
	DisableFlagsInUseLine: true,
	DisableSuggestions:    true,
	GroupID:               "framework",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if env.USER.Uid != "0" {
			return fmt.Errorf("require admin access")
		}

		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		argsMap := make(map[string]bool)
		for _, arg := range args {
			argsMap[arg+"-"+runtime.GOOS+"-"+runtime.GOARCH] = true
		}

		if compareVersions(env.BUILD_VERSION, latest().TagName) {
			logger.Default().Message("looking for latest version " + color.CyanString(latest().TagName))
			var err error = nil
			for _, asset := range latest().Assets {
				isKitsune := strings.Contains(asset.Name, "kitsune-"+runtime.GOOS+"-"+runtime.GOARCH)
				isService := strings.Contains(asset.Name, runtime.GOOS+"-"+runtime.GOARCH)

				if len(args) == 0 || argsMap[asset.Name] || isKitsune {
					if isKitsune {
						err = asset.Download(env.PATH_BIN)
					} else if isService {
						err = asset.Download(env.PATH_SERVICES)
					}

					if err != nil {
						return err
					}
				}
			}
		}

		return nil
	},
}
