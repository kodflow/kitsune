package install

import (
	"fmt"
	"os/user"
	"runtime"
	"strings"

	"github.com/kodmain/kitsune/src/internal/env"
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
	PreRunE: func(cmd *cobra.Command, args []string) error {
		current, err := user.Current()
		if err != nil {
			return err
		}

		if current.Uid != "0" {
			return fmt.Errorf("require admin privilege")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		// Convert arguments to map for efficient lookup
		argsMap := make(map[string]bool)
		for _, arg := range args {
			argsMap[arg+"-"+runtime.GOOS+"-"+runtime.GOARCH] = true
		}

		// Logic for status command
		if compareVersions(env.BUILD_VERSION, latest().TagName) {
			fmt.Println("kitsune install to latest version", latest().TagName)
			for _, asset := range latest().Assets {
				var err error = nil
				isKitsune := strings.Contains(asset.Name, "kitsune-"+runtime.GOOS+"-"+runtime.GOARCH)
				isService := strings.Contains(asset.Name, runtime.GOOS+"-"+runtime.GOARCH)
				if len(args) == 0 || argsMap[asset.Name] || isKitsune {
					if isKitsune {
						err = asset.Download(env.PATH_BIN)
					} else if isService {
						err = asset.Download(env.PATH_SERVICES)
					}
				}

				if err != nil {
					fmt.Println("Failed to install", err.Error())
					break
				}
			}
		} else {
			fmt.Println("kitsune is up to date")
		}
	},
}
