package kitsune

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kodmain/KitsuneFramework/internal/env"
	"github.com/spf13/cobra"
)

var forceRun bool

func init() {
	startCmd.Flags().BoolVarP(&forceRun, "forground", "f", false, "run service in forground")
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start all micro-services",
	Run: func(cmd *cobra.Command, args []string) {
		// Logic for up command
		fmt.Println("Starting the kitsune services")
		if forceRun {
			exec := exec.Command(filepath.Join(env.PATH_LIB, env.BUILD_APP_NAME, "supervisor"))
			exec.Stdout = os.Stdout
			exec.Stderr = os.Stderr
			exec.Run()
		} else {
			exec.Command(filepath.Join(env.PATH_LIB, env.BUILD_APP_NAME, "supervisor")).Start()
		}
	},
}
