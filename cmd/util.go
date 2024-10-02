package cmd

import (
	"graphgen/lib/graph"
	"os"
)

func loadGraphFromFile(path string) (*graph.Graph, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			panic(err)
		}
	}(f)

	stat, err := f.Stat()
	if err != nil {
		return nil, err
	}
	fileBytes := make([]byte, stat.Size())
	_, err = f.Read(fileBytes)
	if err != nil {
		return nil, err
	}

	g, err := graph.NewParser().Parse(string(fileBytes))
	if err != nil {
		return nil, err
	}
	return g, nil
}
