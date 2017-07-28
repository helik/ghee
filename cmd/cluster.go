// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"bytes"
	"fmt"
	"text/tabwriter"

	"github.com/spf13/cobra"

	"github.com/helik/ghee/controller"
)

// clusterCmd represents the cluster command
var clusterCmd = &cobra.Command{
	Use:   "cluster",
	Short: "Add and list clusters",
}

var clusterListCmd = &cobra.Command{
	Use:   "list",
	Short: "List clusters",
	Run: func(cmd *cobra.Command, args []string) {
		clusters, err := controller.GetClusters()
		if err != nil {
			fmt.Println(err)
			return
		}
		if len(clusters) == 0 {
			fmt.Println("No clusters.")
			return
		}

		// write clusters in nice table
		w := tabwriter.Writer{}
		buf := bytes.Buffer{}
		w.Init(&buf, 0, 8, 4, ' ', 0)
		fmt.Fprintf(&w, "Name\tAddress\n") // table header
		for _, cluster := range clusters {
			fmt.Fprintf(&w, "%v\t%v\n", cluster.Name, cluster.Address)
		}
		w.Flush()
		fmt.Print(buf.String())
	},
}

var clusterAddCmd = &cobra.Command{
	Use:   "add [path to Ghee cluster file]",
	Short: "Add cluster",
	Run: func(cmd *cobra.Command, args []string) {
		path := args[0]
		cluster, err := controller.ReadGheeClusterFile(path)
		if err != nil {
			fmt.Println(err)
			return
		}

		if err := controller.AddCluster(cluster); err != nil {
			fmt.Println(err)
		}
	},
	Args: cobra.ExactArgs(1),
}

var clusterRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove cluster",
}

func init() {
	clusterCmd.AddCommand(clusterListCmd)
	clusterCmd.AddCommand(clusterAddCmd)
	clusterCmd.AddCommand(clusterRemoveCmd)
	RootCmd.AddCommand(clusterCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// clusterCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// clusterCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
