package kitsune

import (
	"fmt"

	"github.com/spf13/cobra"
)

var buildCmd = &cobra.Command{
	Use:   "build",
	Short: "Build project",
	Run: func(cmd *cobra.Command, args []string) {
		// Logic for build command
		fmt.Println("Building the microservice project...")
	},
}
