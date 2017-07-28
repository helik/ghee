package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/helik/ghee/controller"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete [path to Ghee manifest]",
	Short: "Delete resources",
	Long: `Deletes resources according to a specified Ghee manifest. It takes one positional argument: a path to a Ghee manifest. Example usage:

ghee create Gheefile`,
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		gheefile, err := controller.ReadGheeManifest(path)
		if err != nil {
			fmt.Println(err)
			return
		}
		controller.Delete(gheefile)
	},
}

func init() {
	RootCmd.AddCommand(deleteCmd)
}
