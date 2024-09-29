/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"graphgen/internal/config"
	"graphgen/internal/server"
	"log"
	"time"
)

const (
	gracefulShutdownTimeout = time.Second * 15
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
		s := server.NewServer(EnvConfig)

		printServerInitMessage(EnvConfig)
		s.StartAndBlock()
		printServerShutdownMessage()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().StringVarP(&serverAddress, "address", "a", "0.0.0.0", "Bind address for the server")
	serveCmd.Flags().IntVarP(&serverPort, "port", "p", 8080, "Bind port for the server")
}

func printServerInitMessage(c *config.Config) {
	log.Println("🚀 Initializing Server 🚀")
	log.Printf("\t- Address: %v:%v\n", c.Server.Address, c.Server.Port)
	log.Printf("\t- Environment: %v\n", c.Server.Environment)
	log.Printf("\t- Address: %v:%v\n", c.Database.Host, c.Database.Port)
}

func printServerShutdownMessage() {
	fmt.Println("Server Shutdown")
}
