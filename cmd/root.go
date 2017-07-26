package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "ghee",
	Short: "Simple Kubernetes multi-cluster controller",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("root")
	},
}

func Execute() {
	RootCmd.Execute()
}
