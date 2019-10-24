/*
Copyright © 2019 Thorsten Klein <iwilltry42@gmail.com>

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package get

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/rancher/k3d/pkg/cluster"
	"github.com/rancher/k3d/pkg/runtimes"
	k3d "github.com/rancher/k3d/pkg/types"
	"github.com/spf13/cobra"

	log "github.com/sirupsen/logrus"

	"github.com/liggitt/tabwriter"
)

// NewCmdGetCluster returns a new cobra command
func NewCmdGetCluster() *cobra.Command {

	// create new command
	cmd := &cobra.Command{
		Use:     "cluster",
		Aliases: []string{"clusters"},
		Short:   "Get cluster",
		Long:    `Get cluster.`,
		Run: func(cmd *cobra.Command, args []string) {
			log.Debugln("get cluster called")
			c, rt := parseGetClusterCmd(cmd, args)
			var existingClusters []*k3d.Cluster
			if c == nil { // Option a)  no cluster name specified -> get all clusters
				found, err := cluster.GetClusters(rt)
				if err != nil {
					log.Fatalln(err)
				}
				existingClusters = append(existingClusters, found...)
			} else { // Option b) cluster name specified -> get specific cluster
				found, err := cluster.GetCluster(c, rt)
				if err != nil {
					log.Fatalln(err)
				}
				existingClusters = append(existingClusters, found)
			}
			// print existing clusters
			printClusters(existingClusters)
		},
	}

	// add subcommands

	// done
	return cmd
}

func parseGetClusterCmd(cmd *cobra.Command, args []string) (*k3d.Cluster, runtimes.Runtime) {
	// --runtime
	rt, err := cmd.Flags().GetString("runtime")
	if err != nil {
		log.Fatalln("No runtime specified")
	}
	runtime, err := runtimes.GetRuntime(rt)
	if err != nil {
		log.Fatalln(err)
	}

	// Args = cluster name
	if len(args) == 0 {
		return nil, runtime
	}

	cluster := &k3d.Cluster{Name: args[0]} // TODO: validate name first?

	return cluster, runtime
}

// TODO: improve (tabular output or output similar to kubectl)
func printClusters(clusters []*k3d.Cluster) {

	tabwriter := tabwriter.NewWriter(os.Stdout, 6, 4, 3, ' ', tabwriter.RememberWidths)
	defer tabwriter.Flush()

	headers := []string{"NAME", "MASTERS", "WORKERS"} // TODO: add status
	_, err := fmt.Fprintf(tabwriter, "%s\n", strings.Join(headers, "\t"))
	if err != nil {
		log.Fatalln("Failed to print headers")
	}

	sort.Slice(clusters, func(i, j int) bool {
		return clusters[i].Name < clusters[j].Name
	})

	for _, cluster := range clusters {
		masterCount := 0
		workerCount := 0
		for _, node := range cluster.Nodes {
			if node.Role == k3d.MasterRole {
				masterCount++
			} else if node.Role == k3d.WorkerRole {
				workerCount++
			}
		}
		fmt.Fprintf(tabwriter, "%s\t%d\t%d\n", cluster.Name, masterCount, workerCount)
	}
}