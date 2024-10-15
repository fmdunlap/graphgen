package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"graphgen/lib/graph"
	"log"
	"os"
)

func init() {
	runCmd.AddCommand(writeCmd)
}

var writeCmd = &cobra.Command{
	Use:   "write",
	Short: "Write a graph file dynamically",
	Run: func(cmd *cobra.Command, args []string) {
		generator, err := graph.NewGenerator(
			graph.WithNodeCount(5),
			graph.WithEdgeProbability(0.5, false),
			//graph.WithEdgeLabelSet([]string{
			//	"lol",
			//	"xyz",
			//	"abc",
			//	"giraffe",
			//	"goblin",
			//	"cheetah",
			//	"elephant",
			//	"doorknob",
			//	"tiny",
			//}, true),
			graph.WithNodeIdSet([]string{
				"alpha",
				"beta",
				"charlie",
				"delta",
				"echo",
			}),
			graph.WithNodeLabelSet([]string{
				"ABEL",
				"BRAVO",
				"CHICKEN",
				"DRIVER",
				"EAGLE",
			}, false),
		)
		if err != nil {
			log.Fatal(err)
		}
		g := generator.Generate()

		fmt.Println(g.String())

		writer := graph.NewWriter(g, graph.WithOutputFormat(graph.OutputFormatText))
		outputPath := "./tmp/output"

		outputFile, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		err = os.Truncate(outputPath, 0)
		if err != nil {
			panic(err)
		}
		defer outputFile.Close()

		err = writer.Write(outputFile)
		if err != nil {
			panic(err)
		}
	},
}
