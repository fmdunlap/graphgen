package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"graphgen/internal/config"
	"graphgen/internal/server"
	"log"
)

var (
	serverPort    int
	serverAddress string
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Run the graphgen server",
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
