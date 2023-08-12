package services

import (
	"fmt"

	"github.com/spf13/cobra"
)

var status = &cobra.Command{
	Use:   "status",
	Short: "Display project status",
	Run: func(cmd *cobra.Command, args []string) {
		// Logic for status command
		fmt.Println("Checking the status of the microservice...")
	},
}
