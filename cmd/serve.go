/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"graphgen/internal/server"
)

var (
	serverPort    int
	serverAddress string
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		s := server.NewServer()

		err := s.ListenAndServe()
		if err != nil {
			panic(fmt.Sprintf("cannot start server: %s", err))
		}
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&serverAddress, "address", "a", "0.0.0.0", "Bind address for the server")
	serveCmd.Flags().IntVarP(&serverPort, "port", "p", 8080, "Bind port for the server")

	viper.BindPFlag("address", serveCmd.Flags().Lookup("address"))
	viper.BindPFlag("port", serveCmd.Flags().Lookup("port"))
}

func initConfig() {
	viper.AutomaticEnv()
}
