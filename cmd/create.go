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
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("Path to Gheefile not provided.")
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
		fmt.Printf("%+v", gheefile)
	},
}

func init() {
	RootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// createCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
