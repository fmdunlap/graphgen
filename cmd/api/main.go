package main

import (
	"fmt"
	"graphgen/lib/graph"
	"os"
)

func main() {
	//s := server.NewServer()
	//
	//err := s.ListenAndServe()
	//if err != nil {
	//    panic(fmt.Sprintf("cannot start server: %s", err))
	//}

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

	graph, err := parser.Parse(string(fileBytes))
	if err != nil {
		panic(err)
	}
	dijkstra(graph, graph.NodeId(os.Args[2]), graph.NodeId(os.Args[3]))
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
