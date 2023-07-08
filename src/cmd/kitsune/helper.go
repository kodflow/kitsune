package kitsune

import (
	"github.com/kodmain/kitsune/src/cmd/kitsune/update"
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
	Helper.AddCommand(update.Cmd)

	serviceCmd.AddCommand(startCmd)
	serviceCmd.AddCommand(stopCmd)
	serviceCmd.AddCommand(logCmd)
	serviceCmd.AddCommand(statusCmd)

	Helper.CompletionOptions.DisableDefaultCmd = true
	Helper.CompletionOptions.DisableNoDescFlag = true

}
