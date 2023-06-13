package kitsune

import (
	"github.com/spf13/cobra"
)

var Helper *cobra.Command = &cobra.Command{
	Use:   "kitsune",
	Short: "Kitsune is a microservice-oriented framework in Go",
	Long:  "Kitsune is a powerful and flexible framework for building microservices in Go. It provides a comprehensive set of commands to manage and streamline your microservice development process.",

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
	//Helper.PersistentFlags().StringVarP(&LOG_LEVEL, "log-level", "l", LOG_LEVEL, "define log level")
	//Helper.PersistentFlags().IntVarP(&LOG_TYPE, "log-type", "t", LOG_TYPE, "define log type")
	//Helper.PersistentFlags().StringVarP(&LOG_PATH, "log", "L", LOG_PATH, "define log path")
	//Helper.PersistentFlags().BoolVarP(&ENABLE_GRPC, "enable-grpc", "G", ENABLE_GRPC, "enable GRPC")
	//Helper.PersistentFlags().BoolVarP(&ENABLE_HTTP, "enable-http", "H", ENABLE_HTTP, "enable HTTP")
	//Helper.PersistentFlags().BoolVarP(&ENABLE_TLS, "enable-tls", "s", ENABLE_TLS, "require certificats TLS for HTTP/GRPC")
	//Helper.PersistentFlags().StringVarP(&TLS_PATH, "tls", "S", TLS_PATH, "define certificats TLS path")
	//Helper.PersistentFlags().StringVarP(&UPLOAD_PATH, "upload", "u", UPLOAD_PATH, "define upload folder path")
	//Helper.PersistentFlags().StringVarP(&STATIC_PATH, "static", "f", STATIC_PATH, "define static folder path")
	//Helper.PersistentFlags().StringVarP(&TMP_PATH, "tmp", "F", TMP_PATH, "define tmp folder path")
	//Helper.PersistentFlags().StringVarP(&ENVIRONMENT, "env", "e", ENVIRONMENT, "define log level")

	Helper.CompletionOptions.DisableDefaultCmd = true
	Helper.CompletionOptions.DisableNoDescFlag = true

}
