package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var mytetsCmd = &cobra.Command{
	Use:   "mytets",
	Short: "A brief description of your command",
	Long:  `A longer description of your command.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("mytets called")
	},
}

func main() {
	if err := mytetsCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
