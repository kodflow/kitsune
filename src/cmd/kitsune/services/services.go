package services

import (
	"github.com/spf13/cobra"
)

var Cmd = &cobra.Command{
	Use:     "services",
	Short:   "Manage services",
	GroupID: "framework",
}

func init() {
	Cmd.AddCommand(start)
	Cmd.AddCommand(stop)
	Cmd.AddCommand(logs)
	Cmd.AddCommand(status)
}
