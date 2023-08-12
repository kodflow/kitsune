package services

import (
	"fmt"

	"github.com/spf13/cobra"
)

var stop = &cobra.Command{
	Use:   "stop",
	Short: "Stop all kitsune-services",
	Run: func(cmd *cobra.Command, args []string) {
		// Logic for down command
		fmt.Println("Stopping all services...")
	},
}
