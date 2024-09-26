/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"graphgen/lib/graph"
	"os"

	"github.com/spf13/cobra"
)

// dijkstraCmd represents the dijkstra command
var dijkstraCmd = &cobra.Command{
	Use:   "dijkstra",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		graphPath := os.Args[1]

		f, err := os.OpenFile(graphPath, os.O_RDONLY, 0644)
		if err != nil {
			panic(err)
		}
		defer func(f *os.File) {
			err := f.Close()
			if err != nil {
				panic(err)
			}
		}(f)

		stat, err := f.Stat()
		if err != nil {
			panic(err)
		}
		fileBytes := make([]byte, stat.Size())
		_, err = f.Read(fileBytes)
		if err != nil {
			panic(err)
		}

		parser := graph.NewParser()

		g, err := parser.Parse(string(fileBytes))
		if err != nil {
			panic(err)
		}
		dijkstra(g, graph.NodeId(os.Args[2]), graph.NodeId(os.Args[3]))
	},
}

func init() {
	runCmd.AddCommand(dijkstraCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// dijkstraCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// dijkstraCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func dijkstra(g *graph.Graph, source, target graph.NodeId) {
	sourceNode, err := g.GetNode(source)
	if err != nil {
		fmt.Printf("Source doesn't exist")
		return
	}
	targetNode, err := g.GetNode(target)
	if err != nil {
		fmt.Printf("Source doesn't exist")
		return
	}

	explored := make(map[string]string)
	queue := []*graph.Node{sourceNode}

	for {
		if len(queue) == 0 {
			return
		}
		node := queue[0]
		queue = queue[1:]

		if node.ID == targetNode.ID {
			backtrace(string(node.ID), string(source), explored)
			return
		}

		for _, reachable := range g.ReachableNodes(node.ID) {
			if _, ok := explored[string(reachable.ID)]; ok {
				continue
			}
			queue = append(queue, reachable)
			explored[string(reachable.ID)] = string(node.ID)
		}
	}
}

func backtrace(from, target string, explored map[string]string) {
	path := ""
	current := from
	for {
		path = current + " " + path
		if current == target {
			break
		}
		next, ok := explored[current]
		if !ok {
			break
		}
		current = next
	}
	fmt.Println(path)
}
