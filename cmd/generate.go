package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"graphgen/lib/graph"
	"log"
)

func init() {
	rootCmd.AddCommand(generateCmd)
}

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a graph file dynamically",
	Run: func(cmd *cobra.Command, args []string) {
		generator, err := graph.NewGenerator(
			graph.WithNodeCount(5),
			graph.WithEdgeProbability(0.5, false),
			graph.WithEdgeLabelSet([]string{
				"lol",
				"xyz",
				"abc",
				"giraffe",
				"goblin",
				"cheetah",
				"elephant",
				"doorknob",
				"tiny",
			}, true),
			graph.WithNodeIdSet([]string{
				"alpha",
				"beta",
				"charlie",
				"delta",
				"echo",
			}),
			graph.WithEdgeWeightRange(-1.0, 10.0),
		)
		if err != nil {
			log.Fatal(err)
		}
		g := generator.Generate()
		fmt.Print(g.String())
	},
}
