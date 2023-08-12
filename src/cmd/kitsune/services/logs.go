package services

import (
	"fmt"

	"github.com/spf13/cobra"
)

var logs = &cobra.Command{
	Use:   "logs",
	Short: "display all kitsune-services logs",
	Run: func(cmd *cobra.Command, args []string) {
		// Logic for down command
		fmt.Println("display log for all services...")
	},
}
