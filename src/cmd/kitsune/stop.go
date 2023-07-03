package kitsune

import (
	"fmt"

	"github.com/spf13/cobra"
)

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop all kitsune-services",
	Run: func(cmd *cobra.Command, args []string) {
		// Logic for down command
		fmt.Println("Stopping all services...")
	},
}
