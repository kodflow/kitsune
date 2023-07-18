package kitsune

import (
	"github.com/kodmain/kitsune/src/cmd/kitsune/install"
	"github.com/kodmain/kitsune/src/internal/env"
	"github.com/spf13/cobra"
)

var Helper *cobra.Command = &cobra.Command{
	Use:     "kitsune",
	Version: env.BUILD_VERSION,
	Short:   "Kitsune is a microservice-oriented framework in Go",
	Long:    "Kitsune is a powerful and flexible framework for building microservices in Go.",

	DisableAutoGenTag:     true,
	DisableFlagsInUseLine: true,
}

func init() {

	Helper.AddCommand(initCmd)
	Helper.AddCommand(buildCmd)
	Helper.AddCommand(serviceCmd)

	Helper.AddCommand(install.Cmd)

	serviceCmd.AddCommand(startCmd)
	serviceCmd.AddCommand(stopCmd)
	serviceCmd.AddCommand(logCmd)
	serviceCmd.AddCommand(statusCmd)

	Helper.CompletionOptions.DisableDefaultCmd = true
	Helper.CompletionOptions.DisableNoDescFlag = true

}
