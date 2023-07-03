package kitsune

import (
	"fmt"

	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "display all kitsune-services logs",
	Run: func(cmd *cobra.Command, args []string) {
		// Logic for down command
		fmt.Println("display log for all services...")
	},
}
