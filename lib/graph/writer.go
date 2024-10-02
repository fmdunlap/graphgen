package graph

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

type WriterOptionFunc func(*Generator)

func WithOutputFormat(format OutputFormat) func(*Writer) {
	return func(w *Writer) {
		w.OutputFormat = format
	}
}

func NewWriter(g *Graph, o ...WriterOptionFunc) *Writer {
	return &Writer{
		graph: g,
	}
}
