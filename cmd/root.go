package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "ghee",
	Short: "Simple Kubernetes multi-cluster controller",
}

func Execute() {
	RootCmd.Execute()
}
