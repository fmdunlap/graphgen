package graph

import (
	"fmt"
	"os"
)

type OutputFormat string

const (
	OutputFormatJson OutputFormat = "JSON"
	OutputFormatYaml OutputFormat = "YAML"
	OutputFormatText OutputFormat = "TEXT"
)

type Writer struct {
	graph        *Graph
	OutputFormat OutputFormat
}

type WriterOptionFunc func(*Writer)

func WithOutputFormat(format OutputFormat) func(*Writer) {
	return func(w *Writer) {
		w.OutputFormat = format
	}
}

func NewWriter(g *Graph, opts ...WriterOptionFunc) *Writer {
	writer := &Writer{
		graph:        g,
		OutputFormat: OutputFormatText,
	}

	for _, opt := range opts {
		opt(writer)
	}

	return writer
}

func (w *Writer) Write(file *os.File) error {
	var graphBytes []byte
	var err error

	switch w.OutputFormat {
	case OutputFormatText:
		graphBytes, err = w.graphTextBytes()
		if err != nil {
			return err
		}
	default:
		panic("not implemented")
	}

	_, err = file.Write(graphBytes)
	if err != nil {
		return err
	}

	return nil
}

func (w *Writer) graphTextBytes() ([]byte, error) {
	graphString := ""
	for _, node := range w.graph.Nodes {
		nodeString := node.definitionString() + "\n"
		graphString += nodeString
	}

	// Newline for style
	graphString += "\n"

	for _, edgeMap := range w.graph.Edges {
		for _, edge := range edgeMap {
			edgeString := edge.definitionString() + "\n"
			graphString += edgeString
		}
	}

	fmt.Println(graphString)

	return []byte(graphString), nil
}
