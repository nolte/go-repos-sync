package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var Version = "development"
var Commit = "development"

var Time string
var User string

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Version:\t", Version)
		fmt.Println("Commit:\t\t", Commit)
		fmt.Println("Time:\t\t", Time)
		fmt.Println("User:\t\t", User)
	},
}
