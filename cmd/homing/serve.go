package main

import (
	"github.com/spf13/cobra"
	"github.com/zuxt268/homing/internal/infrastructure/server"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the homing server",
	Long:  `Start the homing server to handle API requests and web traffic.`,
	Run: func(cmd *cobra.Command, args []string) {
		server.Run()
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}