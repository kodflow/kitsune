package kitsune

import (
	"fmt"

	"github.com/spf13/cobra"
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start all micro-services",
	Run: func(cmd *cobra.Command, args []string) {
		// Logic for up command
		fmt.Println("Starting the microservice project...")
	},
}
