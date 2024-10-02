package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"graphgen/lib/graph"
)

func init() {
	runCmd.AddCommand(dijkstraCmd)
}

// dijkstraCmd represents the dijkstra command
var dijkstraCmd = &cobra.Command{
	Use:   "dijkstra",
	Short: "Runs dijkstra on a supplied graph file to provide shortest unweighted path from source to target",
	Args:  cobra.MatchAll(cobra.ExactArgs(3), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		g, err := loadGraphFromFile(args[0])
		if err != nil {
			panic(err)
		}
		dijkstra(g, graph.NodeId(args[1]), graph.NodeId(args[2]))
	},
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
