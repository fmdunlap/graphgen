package graphparse

import (
    "fmt"
    "sort"
)

type Graph struct {
    Nodes map[NodeId]*Node
    Edges map[NodeId]map[NodeId]*Edge
}

func NewGraph() *Graph {
    return &Graph{
        Nodes: map[NodeId]*Node{},
        Edges: map[NodeId]map[NodeId]*Edge{},
    }
}

func (g *Graph) AddNode(node ...Node) {
    for _, n := range node {
        g.Nodes[n.ID] = &n
    }
}

func (g *Graph) RemoveNode(id NodeId) {
    delete(g.Nodes, id)
}

func (g *Graph) GetNode(id NodeId) (*Node, error) {
    if node, ok := g.Nodes[id]; ok {
        return node, nil
    }
    return nil, &NodeNotFoundError{ID: id}
}

func (g *Graph) NodeExists(id NodeId) bool {
    _, err := g.GetNode(id)
    return err == nil
}

func (g *Graph) ReachableNodes(id NodeId) []*Node {
    reachableNodes := make([]*Node, 0)

    sourceEdges, ok := g.Edges[id]
    // If there are no edges with this node's id as the source, the there are no reachable edges.
    if !ok {
        return reachableNodes
    }

    for targetId, _ := range sourceEdges {
        targetNode, err := g.GetNode(targetId)
        if err != nil {
            panic(err)
        }
        reachableNodes = append(reachableNodes, targetNode)
    }
    return reachableNodes
}

func (g *Graph) AddEdge(edge Edge) {
    sourceEdges, ok := g.Edges[edge.Source]
    if !ok {
        g.Edges[edge.Source] = make(map[NodeId]*Edge)
        sourceEdges = g.Edges[edge.Source]
    }

    sourceEdges[edge.Target] = &edge

}

func (g *Graph) RemoveEdge(source NodeId, target NodeId) {
    if _, ok := g.Edges[source]; !ok {
        return
    }
    delete(g.Edges[source], target)
}

func (g *Graph) GetEdge(source NodeId, target NodeId) (*Edge, error) {
    sourceEdges, ok := g.Edges[source]
    if !ok {
        return nil, &EdgeNotFoundError{
            Source: source,
            Target: target,
        }
    }
    edge, ok := sourceEdges[target]
    if !ok {
        return nil, &EdgeNotFoundError{
            Source: source,
            Target: target,
        }
    }
    return edge, nil
}

func (g *Graph) EdgeExists(source NodeId, target NodeId) bool {
    _, err := g.GetEdge(source, target)
    return err == nil
}

func (g *Graph) String() string {
    var output string
    for _, node := range g.Nodes {
        label := string(node.ID)
        if node.Attrs.Label != "" {
            label = node.Attrs.Label
        }
        output += fmt.Sprintf("%v\n", label)
    }

    output += "\n"

    sortedEdges := make([]Edge, 0)
    for _, sourceEdge := range g.Edges {
        for _, edge := range sourceEdge {
            sortedEdges = append(sortedEdges, *edge)
        }
    }
    sort.Slice(sortedEdges, func(i, j int) bool {
        return sortedEdges[i].Source < sortedEdges[j].Source
    })

    for _, edge := range sortedEdges {
        source := g.Nodes[edge.Source]
        target := g.Nodes[edge.Target]

        sourceLabel := string(source.ID)
        if source.Attrs.Label != "" {
            sourceLabel = source.Attrs.Label
        }
        targetLabel := string(target.ID)
        if target.Attrs.Label != "" {
            targetLabel = target.Attrs.Label
        }

        output += fmt.Sprintf("%v -> %v\n", sourceLabel, targetLabel)
    }
    return output
}
