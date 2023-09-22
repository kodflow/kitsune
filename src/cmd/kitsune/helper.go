package kitsune

import (
	"github.com/kodmain/kitsune/src/cmd/kitsune/install"
	"github.com/kodmain/kitsune/src/cmd/kitsune/services"
	"github.com/kodmain/kitsune/src/config"
	"github.com/spf13/cobra"
)

var Helper *cobra.Command = &cobra.Command{
	Use:                   "kitsune",
	Version:               config.BUILD_VERSION,
	Short:                 "Kitsune (" + config.BUILD_VERSION + ") is a microservice-oriented framework in Go",
	Long:                  "Kitsune (" + config.BUILD_VERSION + ") is a powerful and flexible framework for building microservices in Go.",
	SilenceUsage:          true,
	DisableAutoGenTag:     true,
	DisableFlagsInUseLine: true,
}

func init() {
	Helper.SetUsageTemplate(Template)

	Helper.AddCommand(initCmd)
	Helper.AddCommand(buildCmd)
	Helper.AddCommand(services.Cmd)
	Helper.AddCommand(install.Cmd)
	Helper.AddGroup(&cobra.Group{ID: "framework", Title: "Framework Commands:"})
	Helper.AddGroup(&cobra.Group{ID: "project", Title: "Project Commands:"})

	Helper.SetHelpCommand(&cobra.Command{GroupID: "framework", Hidden: true})
	Helper.PersistentFlags().BoolP("help", "h", false, "Print usage")
	Helper.PersistentFlags().Lookup("help").Hidden = true

	Helper.CompletionOptions.DisableDefaultCmd = true
	Helper.CompletionOptions.DisableNoDescFlag = true
}
