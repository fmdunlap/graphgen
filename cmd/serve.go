/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"graphgen/internal/config"
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
		s := server.NewServer(EnvConfig)

		printServerInitMessage(EnvConfig)

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
}

func printServerInitMessage(c *config.Config) {
	initFmtString := `🚀 Initializing Server 🚀
- Address: %v:%v
- Environment: "%v"
- Database Host: %v:%v
`

	fmt.Printf(initFmtString,
		c.Server.Address,
		c.Server.Port,
		c.Server.Environment,
		c.Database.Host,
		c.Database.Port,
	)
}
