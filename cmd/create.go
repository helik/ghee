package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/helik/ghee/controller"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create [path to Gheefile]",
	Short: "Create a resource",
	Long: `Creates resources according to a specified Gheefile. It takes one positional argument: a path to a Gheefile. Example usage:

ghee create Gheefile`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("Expects one path to Gheefile")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		gheefile, err := controller.ReadGheefile(path)
		if err != nil {
			fmt.Println(err)
			return
		}
		controller.Create(gheefile)
	},
}

func init() {
	RootCmd.AddCommand(createCmd)
}
