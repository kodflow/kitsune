package kitsune

import (
	"github.com/spf13/cobra"
)

var Helper *cobra.Command = &cobra.Command{
	Use:   "kitsune",
	Short: "Kitsune is a microservice-oriented framework in Go",
	Long:  "Kitsune is a powerful and flexible framework for building microservices in Go.",

	DisableAutoGenTag:     true,
	DisableFlagsInUseLine: true,
}

func init() {

	Helper.AddCommand(initCmd)
	Helper.AddCommand(buildCmd)
	Helper.AddCommand(serviceCmd)
	Helper.AddCommand(statusCmd)
	serviceCmd.AddCommand(startCmd)
	serviceCmd.AddCommand(stopCmd)
	serviceCmd.AddCommand(logCmd)

	Helper.CompletionOptions.DisableDefaultCmd = true
	Helper.CompletionOptions.DisableNoDescFlag = true

}
