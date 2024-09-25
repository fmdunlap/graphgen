package graphparse

import "fmt"

type NodeNotFoundError struct {
    ID NodeId
}

func (e *NodeNotFoundError) Error() string {
    return fmt.Sprintf("Node %v not found", e.ID)
}

type EdgeNotFoundError struct {
    Source NodeId
    Target NodeId
}

func (e *EdgeNotFoundError) Error() string {
    return fmt.Sprintf("Edge %v -> %v not found", e.Source, e.Target)
}

type InvalidLineError struct {
    Line string
}

func (e *InvalidLineError) Error() string {
    return fmt.Sprintf("Invalid line: %s", e.Line)
}

type InvalidAttributeError struct {
    Attribute string
    Value     string
}

func (e *InvalidAttributeError) Error() string {
    return fmt.Sprintf("Invalid attribute: %s = %s", e.Attribute, e.Value)
}

type InvalidEdgeRelationError struct {
    Source    NodeId
    Target    NodeId
    Separator string
}

func (e *InvalidEdgeRelationError) Error() string {
    return fmt.Sprintf("Invalid edge relation: %v %s %v", e.Source, e.Separator, e.Target)
}
