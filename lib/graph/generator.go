package graph

import (
	"fmt"
	"math/rand"
	"sort"
	"strconv"
)

type InsufficientIdSetError struct {
	nodeCount int
	nodeIdSet []string
}

func (e *InsufficientIdSetError) Error() string {
	return fmt.Sprintf(
		"length of node id set (%v) must be greater than or equal to request nodes (%v)",
		len(e.nodeIdSet),
		e.nodeCount,
	)
}

type InsufficientLabelSetError struct {
	labelType      string
	count          int
	labelSetLength int
}

func (e *InsufficientLabelSetError) Error() string {
	return fmt.Sprintf(
		"length of %v label set (%v) must be greater than or equal to request nodes (%v)",
		e.labelType,
		e.count,
		e.labelSetLength,
	)
}

type NonUniqueNodeIds struct {
	duplicateId string
}

func (e *NonUniqueNodeIds) Error() string {
	return fmt.Sprintf("node IDs are not unique - found %v more than once in the node id set", e.duplicateId)
}

const (
	defaultNodeCount       = 50
	defaultEdgeProbability = 0.05
)

type Generator struct {
	graph     *Graph
	nodeCount int

	nodeIdSet    *[]string
	nodeLabelSet *[]string
	// if false, node labels are assigned in the order of the provided set
	randomizeNodeLabelAssignment bool

	edgeLabelSet *[]string
	// if false, edge labels are assigned in the order of the provided set
	randomizeEdgeLabelAssignment bool
	edgeWeightMin                float64
	edgeWeightMax                float64
	// The probability that any two nodes are connected with an edge
	edgeProbability float64
	allowEdgeLoop   bool

	usableNodeLabels []string
	usableEdgeLabels []string
}

type GeneratorOptionFunc func(*Generator)

func WithNodeCount(n int) func(*Generator) {
	return func(g *Generator) {
		g.nodeCount = n
	}
}

func WithEdgeProbability(p float64, allowEdgeLoops bool) func(*Generator) {
	return func(g *Generator) {
		g.edgeProbability = p
		g.allowEdgeLoop = allowEdgeLoops
	}
}

func WithNodeIdSet(s []string) func(*Generator) {
	return func(g *Generator) {
		g.nodeIdSet = &s
	}
}

func WithNodeLabelSet(s []string, randomizeOrder bool) func(*Generator) {
	return func(g *Generator) {
		g.nodeLabelSet = &s
		g.randomizeNodeLabelAssignment = randomizeOrder
	}
}

func WithEdgeLabelSet(s []string, randomizeOrder bool) func(*Generator) {
	return func(g *Generator) {
		g.edgeLabelSet = &s
		g.randomizeEdgeLabelAssignment = randomizeOrder
	}
}

func WithEdgeWeightRange(min, max float64) func(*Generator) {
	return func(g *Generator) {
		g.edgeWeightMin = min
		g.edgeWeightMax = max
	}
}

func NewGenerator(opts ...GeneratorOptionFunc) (*Generator, error) {
	generator := &Generator{
		nodeCount:       defaultNodeCount,
		edgeProbability: defaultEdgeProbability,
		edgeWeightMin:   1.0,
		edgeWeightMax:   1.0,
	}

	for _, o := range opts {
		o(generator)
	}

	err := generator.validateGenerator()
	if err != nil {
		return nil, err
	}

	return generator, nil
}

func (gen *Generator) Generate() *Graph {
	gen.usableNodeLabels = gen.populateUsableNodeLabels()
	gen.usableEdgeLabels = gen.populateUsableEdgeLabels()

	fmt.Printf("%v", gen.usableEdgeLabels)

	g := NewGraph()

	newNodes := make([]Node, gen.nodeCount)

	for i := range gen.nodeCount {
		newNodes[i] = gen.generateNode(i)
	}

	g.AddNode(newNodes...)

	for _, source := range newNodes {
		for _, target := range newNodes {
			if !gen.allowEdgeLoop && source.ID == target.ID {
				continue
			}
			if gen.edgeProbability > rand.Float64() {
				g.AddEdge(gen.generateEdge(source, target))
			}
		}
	}

	return g
}

func (gen *Generator) validateGenerator() error {
	if gen.nodeIdSet != nil {
		if len(*gen.nodeIdSet) < gen.nodeCount {
			return &InsufficientIdSetError{
				nodeCount: gen.nodeCount,
				nodeIdSet: *gen.nodeIdSet,
			}
		}
		if nonUniqueId := gen.findNonUniqueId(); nonUniqueId != nil {
			return &NonUniqueNodeIds{duplicateId: *nonUniqueId}
		}
	}

	if gen.nodeLabelSet != nil {
		if len(*gen.nodeLabelSet) < gen.nodeCount {
			return &InsufficientLabelSetError{
				labelType:      "node",
				count:          gen.nodeCount,
				labelSetLength: len(*gen.nodeLabelSet),
			}
		}

	}

	return nil
}

func (gen *Generator) findNonUniqueId() *string {
	sortedNodeIds := make([]string, len(*gen.nodeIdSet))
	copy(sortedNodeIds, *gen.nodeIdSet)
	sort.Strings(sortedNodeIds)
	for i := range len(sortedNodeIds) - 1 {
		if sortedNodeIds[i] == sortedNodeIds[i+1] {
			return &sortedNodeIds[i]
		}
	}
	return nil
}

func (gen *Generator) generateNode(nodeIndex int) Node {
	nodeId := gen.getNodeId(nodeIndex)
	nodeLabel := gen.getNodeLabel(nodeIndex)

	return Node{
		ID: NodeId(nodeId),
		Attrs: NodeAttributes{
			Label: nodeLabel,
		},
	}
}

func (gen *Generator) getNodeId(i int) string {
	if gen.nodeIdSet == nil {
		return strconv.Itoa(i)
	}
	return (*gen.nodeIdSet)[i]
}

func (gen *Generator) getNodeLabel(i int) string {
	if gen.nodeLabelSet == nil {
		return ""
	}

	return gen.usableNodeLabels[i]
}

func (gen *Generator) generateEdge(source Node, target Node) Edge {
	edgeLabel := gen.generateEdgeLabel(source, target)
	edgeWeight := gen.generateEdgeWeight()
	return Edge{
		Source: source.ID,
		Target: target.ID,
		Attrs: EdgeAttributes{
			Label:  edgeLabel,
			Weight: edgeWeight,
		},
	}
}

func (gen *Generator) generateEdgeLabel(source Node, target Node) string {
	if gen.edgeLabelSet == nil || len(gen.usableEdgeLabels) <= 0 {
		return fmt.Sprintf("%s:%s", source.ID, target.ID)
	}

	label := gen.usableEdgeLabels[0]
	gen.usableEdgeLabels = gen.usableEdgeLabels[1:]

	return label
}

func (gen *Generator) generateEdgeWeight() float64 {
	return gen.edgeWeightMin + rand.Float64()*(gen.edgeWeightMax-gen.edgeWeightMin)
}

func (gen *Generator) populateUsableNodeLabels() []string {
	usableNodeLabels := make([]string, 0)
	if gen.nodeLabelSet == nil {
		for i := range gen.nodeCount {
			usableNodeLabels = append(usableNodeLabels, strconv.Itoa(i))
		}
		return usableNodeLabels
	}

	usableNodeLabels = make([]string, len(*gen.nodeLabelSet))
	copy(usableNodeLabels, *gen.nodeLabelSet)

	if gen.randomizeNodeLabelAssignment {
		rand.Shuffle(len(usableNodeLabels), func(i, j int) {
			usableNodeLabels[i], usableNodeLabels[j] = usableNodeLabels[j], usableNodeLabels[i]
		})
	}

	return usableNodeLabels
}

func (gen *Generator) populateUsableEdgeLabels() []string {
	if gen.edgeLabelSet == nil {
		// empty list forces fallback
		return make([]string, 0)
	}

	usableEdgeLabels := make([]string, len(*gen.edgeLabelSet))
	_ = copy(usableEdgeLabels, *gen.edgeLabelSet)

	if gen.randomizeEdgeLabelAssignment {
		rand.Shuffle(len(usableEdgeLabels), func(i, j int) {
			usableEdgeLabels[i], usableEdgeLabels[j] = usableEdgeLabels[j], usableEdgeLabels[i]
		})
	}

	fmt.Println(usableEdgeLabels)

	return usableEdgeLabels
}
