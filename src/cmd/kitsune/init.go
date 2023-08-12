package kitsune

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:     "init",
	Short:   "Initialize project",
	GroupID: "project",
	Run: func(cmd *cobra.Command, args []string) {
		// Logic for init command
		fmt.Println("Initializing project...")
	},
}
