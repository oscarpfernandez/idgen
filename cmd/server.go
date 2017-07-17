package cmd

import (
	"log"
	"net/http"

	"github.com/oscarpfernandez/idgen/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "starts REST server",
	Long:  `Starts the idgen in server mode and reachable through a REST API`,
	Run: func(cmd *cobra.Command, args []string) {
		logrus.Info("Starting server on port 8080")

		// create a server
		router := server.CreateServer()

		// starts the server on port 8080
		log.Fatal(http.ListenAndServe(":8080", router))
	},
}

func init() {
	RootCmd.AddCommand(serverCmd)
	serverCmd.Flags().Int("port", 8080, "port number")
}
