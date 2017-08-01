package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/helik/ghee/controller"
)

// createCmd represents the create command
var updateCmd = &cobra.Command{
	Use:   "update [path to Ghee manifest]",
	Short: "Update resources",
	Long: `Updates resources according to a specified Ghee manifest. It takes one positional argument: a path to a Ghee manifest. Example usage:

ghee update Gheefile`,
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		gheefile, err := controller.ReadGheeManifest(path)
		if err != nil {
			fmt.Println(err)
			return
		}
		controller.Update(gheefile)
	},
	Args: cobra.ExactArgs(1),
}

func init() {
	RootCmd.AddCommand(updateCmd)
}
