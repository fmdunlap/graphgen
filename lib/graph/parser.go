package graph

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

var (
	NodeDefinitionRegex = regexp.MustCompile(`^(\w)(?: \[([\w\s:,]+)])?$`)
	EdgeDefinitionRegex = regexp.MustCompile(`^(\w)\s*(<-|->|<>)\s*(\w)(?: \[([\w\s:,.]+)])?$`)
)

type Parser struct {
	graph *Graph
}

type ParserOption func(p *Parser)

func NewParser(options ...ParserOption) *Parser {
	p := &Parser{
		graph: NewGraph(),
	}
	for _, option := range options {
		option(p)
	}
	return p
}

func (p *Parser) Parse(input string) (*Graph, error) {
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if len(line) == 0 {
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		if NodeDefinitionRegex.MatchString(line) {
			err := p.parseNodeDefinition(line)
			if err != nil {
				return nil, err
			}
			continue
		}
		if EdgeDefinitionRegex.MatchString(line) {
			err := p.parseEdgeDefinition(line)
			if err != nil {
				return nil, err
			}
			continue
		}
		return nil, &InvalidLineError{Line: line}
	}

	return p.graph, nil
}

func (p *Parser) parseNodeDefinition(line string) error {
	matches := NodeDefinitionRegex.FindStringSubmatch(line)
	id := matches[1]
	attributeString := matches[2]
	attributes, err := p.parseNodeAttributes(attributeString)
	if err != nil {
		return err
	}

	node, err := p.graph.GetNode(NodeId(id))
	if err != nil {
		var NodeNotFoundError *NodeNotFoundError
		if errors.As(err, &NodeNotFoundError) {
			node = &Node{
				ID:    NodeId(id),
				Attrs: attributes,
			}
			p.graph.AddNode(*node)
			return nil
		}
		return err
	}

	node.Attrs = attributes
	return nil
}

func (p *Parser) parseNodeAttributes(attributeString string) (NodeAttributes, error) {
	var nodeAttributes NodeAttributes

	if attributeString == "" {
		return nodeAttributes, nil
	}

	attributeStrings := strings.Split(attributeString, ",")
	for _, rawAttribute := range attributeStrings {
		attribute := strings.TrimSpace(rawAttribute)
		if attribute == "" {
			return nodeAttributes, &InvalidLineError{Line: attribute}
		}

		parts := strings.Split(attribute, ":")
		if len(parts) != 2 {
			return nodeAttributes, &InvalidLineError{Line: attribute}
		}

		key := strings.TrimSpace(parts[0])
		key = strings.ToLower(key)
		value := strings.TrimSpace(parts[1])

		switch key {
		case "label":
			nodeAttributes.Label = value
		default:
			return nodeAttributes, &InvalidAttributeError{Attribute: key, Value: value}
		}
	}

	return nodeAttributes, nil
}

func (p *Parser) parseEdgeDefinition(line string) error {
	matches := EdgeDefinitionRegex.FindStringSubmatch(line)
	leftNodeId := matches[1]
	edgeType := matches[2]
	rightNodeId := matches[3]
	attributeString := matches[4]
	attributes, err := p.parseEdgeAttributes(attributeString)
	if err != nil {
		return err
	}

	edges := make([]Edge, 0)
	switch edgeType {
	case "<-":
		edges = append(edges, Edge{
			Source: NodeId(rightNodeId),
			Target: NodeId(leftNodeId),
			Attrs:  attributes,
		})
	case "->":
		edges = append(edges, Edge{
			Source: NodeId(leftNodeId),
			Target: NodeId(rightNodeId),
			Attrs:  attributes,
		})
	case "<>":
		edges = append(edges, Edge{
			Source: NodeId(leftNodeId),
			Target: NodeId(rightNodeId),
			Attrs:  attributes,
		})
		edges = append(edges, Edge{
			Source: NodeId(rightNodeId),
			Target: NodeId(leftNodeId),
			Attrs:  attributes,
		})
	default:
		return &InvalidEdgeRelationError{Source: NodeId(leftNodeId), Target: NodeId(rightNodeId), Separator: edgeType}
	}

	for _, edge := range edges {

		if !p.graph.NodeExists(edge.Source) {
			p.graph.AddNode(Node{ID: edge.Source})
		}

		if !p.graph.NodeExists(edge.Target) {
			p.graph.AddNode(Node{ID: edge.Target})
		}

		e, err := p.graph.GetEdge(edge.Source, edge.Target)
		if err != nil {
			var EdgeNotFoundError *EdgeNotFoundError
			if errors.As(err, &EdgeNotFoundError) {
				p.graph.AddEdge(edge)
				continue
			}
			return err
		}
		e.Attrs = edge.Attrs
	}

	return nil
}

func (p *Parser) parseEdgeAttributes(attributeString string) (EdgeAttributes, error) {
	var edgeAttributes EdgeAttributes

	if attributeString == "" {
		return edgeAttributes, nil
	}

	attributeStrings := strings.Split(attributeString, ",")
	for _, rawAttribute := range attributeStrings {
		attribute := strings.TrimSpace(rawAttribute)
		if attribute == "" {
			return edgeAttributes, &InvalidLineError{Line: attribute}
		}

		parts := strings.Split(attribute, ":")
		if len(parts) != 2 {
			return edgeAttributes, &InvalidLineError{Line: attribute}
		}

		key := strings.TrimSpace(parts[0])
		key = strings.ToLower(key)
		value := strings.TrimSpace(parts[1])
		value = strings.ToLower(value)

		switch key {
		case "label":
			edgeAttributes.Label = value
		case "weight":
			weight, err := strconv.ParseFloat(value, 64)
			if err != nil {
				return edgeAttributes, &InvalidAttributeError{Attribute: key, Value: value}
			}
			edgeAttributes.Weight = weight
		default:
			return edgeAttributes, &InvalidAttributeError{Attribute: key, Value: value}
		}
	}

	return edgeAttributes, nil
}
